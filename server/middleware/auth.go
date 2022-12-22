package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ahmadmirdas/julo-test/config"
	"github.com/ahmadmirdas/julo-test/utils"
	"github.com/ahmadmirdas/julo-test/utils/log"
	"github.com/ahmadmirdas/julo-test/utils/response"
	"github.com/golang-jwt/jwt/v4"
)

type key int

const (
	Customer key = iota
)

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cfg := config.Config.JWTCfg
			authorizationHeader := r.Header.Get("Authorization")
			url_to_skip_auth_check := []string{
				"/api/v1/init",
			}
			skip_check := utils.Contains(r.URL.Path, url_to_skip_auth_check)
			if skip_check {
				next.ServeHTTP(w, r)
				return
			}
			if !strings.Contains(authorizationHeader, "Token") {
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response.ResponseAPI{
					Error_: &response.ApiError{
						Code:    http.StatusInternalServerError,
						Message: "invalid token",
					},
				})
				return
			}
			tokenString := strings.Replace(authorizationHeader, "Token ", "", -1)

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("signing method invalid")
				} else if method != JWT_SIGNING_METHOD {
					return nil, fmt.Errorf("signing method invalid")
				}

				return []byte(cfg.SignKey), nil
			})
			if err != nil {
				log.WithContext(context.Background()).Errorf("Error jwt parse", err)
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response.ResponseAPI{
					Error_: &response.ApiError{
						Code:    http.StatusInternalServerError,
						Message: err.Error(),
					},
				})
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				log.WithContext(context.Background()).Errorf("Error claims", err)
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response.ResponseAPI{
					Error_: &response.ApiError{
						Code:    http.StatusInternalServerError,
						Message: err.Error(),
					},
				})
				return
			}

			ctx := context.WithValue(r.Context(), Customer, claims)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

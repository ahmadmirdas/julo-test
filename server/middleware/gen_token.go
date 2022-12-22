package middleware

import (
	"time"

	"github.com/ahmadmirdas/julo-test/config"
	"github.com/golang-jwt/jwt/v4"
)

var JWT_SIGNING_METHOD = jwt.SigningMethodHS256

type MyClaims struct {
	jwt.RegisteredClaims
	CustomerXId string `json:"customer_xid"`
}

func GenerateToken(customerXId string) (string, error) {
	cfg := config.Config.JWTCfg
	claims := MyClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.Exp) * time.Hour)),
		},
		CustomerXId: customerXId,
	}

	token := jwt.NewWithClaims(
		JWT_SIGNING_METHOD,
		claims,
	)

	signedToken, err := token.SignedString([]byte(cfg.SignKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

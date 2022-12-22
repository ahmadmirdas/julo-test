package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ahmadmirdas/julo-test/config"
	"github.com/ahmadmirdas/julo-test/config/database"
	"github.com/ahmadmirdas/julo-test/handler"
	"github.com/ahmadmirdas/julo-test/repository/database/models"
	"github.com/ahmadmirdas/julo-test/server/middleware"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func RunServer() {
	cfg := config.Config
	cfgDb := cfg.PostgresCfg
	paramCfgDB := database.ParamConn{
		Username:    cfgDb.Username,
		Password:    cfgDb.Password,
		Host:        cfgDb.Host,
		Port:        cfgDb.Port,
		Database:    cfgDb.Database,
		MaxConn:     cfgDb.MaxConn,
		MinIdleConn: cfgDb.MinIdleConn,
		MaxRetries:  cfgDb.MaxRetries,
	}

	db := database.DbConn(paramCfgDB)

	err := db.Ping(context.Background())
	if err != nil {
		logrus.Fatalf("Ping DB error: %v", err)
		log.Fatalln(err)
	}

	walletRepo := models.NewDBWalletRepo(db)
	handlerAPI := handler.NewHandlerWallet(walletRepo)

	// Declare a new router
	r := mux.NewRouter()
	apiV1 := r.PathPrefix("/api/v1").Subrouter()

	apiV1.HandleFunc("/init", handlerAPI.InitAccountWallet).Methods(http.MethodPost)
	apiV1.HandleFunc("/wallet", handlerAPI.EnableWallet).Methods(http.MethodPost)
	apiV1.HandleFunc("/wallet", handlerAPI.ViewWalletBalance).Methods(http.MethodGet)
	apiV1.HandleFunc("/wallet/deposits", handlerAPI.DepositWallet).Methods(http.MethodPost)
	apiV1.HandleFunc("/wallet/withdrawals", handlerAPI.WithdrawWallet).Methods(http.MethodPost)
	apiV1.HandleFunc("/wallet", handlerAPI.DisableWallet).Methods(http.MethodPatch)
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(middleware.AuthMiddleware())

	srv := &http.Server{
		Handler:      r,
		Addr:         ":5000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Starting web on port 5000")
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

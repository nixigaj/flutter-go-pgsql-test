package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"time"
)

func serve(cfg *Config, sc *syncController, dbPool *pgxpool.Pool) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		apiHandler(w, r, dbPool)
	})

	srv := &http.Server{
		Addr:    cfg.Bind,
		Handler: mux,
	}

	go func() {
		log.Infof("Listening on %s", cfg.Bind)
		err := srv.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			sc.fatalChan <- fmt.Errorf("HTTP server startup: %v", err)
		}
	}()

	<-sc.stopChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Errorf("HTTP stutdown: %v", err)
	}

	sc.wg.Done()
}

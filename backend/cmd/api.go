package main

import (
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

func apiHandler(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	defer log.Debugf("Handled HTTP request from %s", r.Host)

	if r.URL.Path == "/" {
		_, err := fmt.Fprint(w, "Hello from app API!")
		if err != nil {
			log.Errorf("HTTP write: %v", err)
			return
		}
		return
	}

	if r.URL.Path == "/sql-hello" {
		row := dbPool.QueryRow(context.Background(), "SELECT 'Hello, world!'")

		var helloResp string
		err := row.Scan(&helloResp)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Errorf("Row scan fail: %v", err)
			return
		}

		_, err = fmt.Fprintf(w, "SQL response: %s", helloResp)
		if err != nil {
			log.Errorf("HTTP write: %v", err)
			return
		}
		return
	}

	http.NotFound(w, r)
}

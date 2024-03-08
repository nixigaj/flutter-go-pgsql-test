package main

import (
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
		countResp, err := incrementHello(dbPool)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Errorf("Database error: %v", err)
			return
		}

		_, err = fmt.Fprintf(w, "SQL response: %s", countResp)
		if err != nil {
			log.Errorf("HTTP write: %v", err)
			return
		}
		return
	}

	http.NotFound(w, r)
}

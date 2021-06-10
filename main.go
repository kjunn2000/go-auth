package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kjunn2000/go-auth/internal/go-auth/handlers"
	"github.com/kjunn2000/go-auth/internal/go-auth/middleware"
)

func main() {
	mr := mux.NewRouter()

	pr := mr.Methods("POST").Subrouter()
	gr := mr.Methods("GET").Subrouter()

	pr.HandleFunc("/api/v1/login", handlers.LoginHandler)
	pr.HandleFunc("/api/v1/refresh-token", handlers.RefreshTokenHandler)
	pr.HandleFunc("/api/v1/account/opening", handlers.AccOpeningHandler)

	gr.HandleFunc("/api/v1/home", handlers.HomeHandler)
	gr.Use(middleware.JwtTokenVerifier)

	srv := &http.Server{
		Handler:      mr,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

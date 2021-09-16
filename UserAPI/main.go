package main

import (
	"UserAPI/handlers"
	"UserAPI/middleware"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	sm := mux.NewRouter()

	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/v1/users/{id:[0-9]+}", handlers.GetID)
	getR.HandleFunc("/v1/users", handlers.GetUsers)

	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/v1/users", handlers.CreateUser)
	postR.Use(middleware.ValidateUser)

	deleteR := sm.Methods(http.MethodDelete).Subrouter()
	deleteR.HandleFunc("/v1/users/{id:[0-9]+}", handlers.DeleteUser)

	s := &http.Server{
		Addr:         "localhost:9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	fmt.Println("Received Terminate signal. Shutting down. ", sig)
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.Shutdown(tc)

}

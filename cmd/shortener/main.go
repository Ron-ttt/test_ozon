package main

import (
	"log"
	"net/http"
	"testozon/internal/app/handlers"
	"testozon/internal/app/middleware"

	"github.com/gorilla/mux"
)

type GRPCsever struct{}

func main() {

	hw := handlers.Init()

	r := mux.NewRouter()
	r.Use(middleware.Logger1, middleware.GzipMiddleware)
	r.HandleFunc("/{id}", hw.Redirect).Methods(http.MethodGet)
	r.HandleFunc("/", hw.IndexPage).Methods(http.MethodPost)

	//s := grpc.NewServer()
	log.Println("server is running")
	err := http.ListenAndServe(hw.Localhost, r)
	if err != nil {
		panic(err)
	}

}

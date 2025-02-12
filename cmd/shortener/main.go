package main

import (
	"log"
	"net"
	"net/http"
	"testozon/internal/app/grpcfunc"
	"testozon/internal/app/handlers"
	"testozon/internal/app/middleware"

	pb "testozon/internal/app/proto"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

type GRPCsever struct{}

func main() {

	hw := handlers.Init()
	grpcfunc.Init()
	r := mux.NewRouter()
	r.Use(middleware.Logger1, middleware.GzipMiddleware)
	r.HandleFunc("/{id}", hw.Redirect).Methods(http.MethodGet)
	r.HandleFunc("/", hw.IndexPage).Methods(http.MethodPost)

	s := grpc.NewServer()
	go func() {
		listen, err := net.Listen("tcp", ":3200")
		if err != nil {
			log.Fatal(err)
		}

		pb.RegisterShortenerServer(s, &grpcfunc.ShortenerServer{})

		log.Println("gRPC server  is running")
		if err := s.Serve(listen); err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("server is running")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

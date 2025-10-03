package main

import (
	"REST_API_app/internal/user"
	"REST_API_app/pkg/logging"
	"net"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := logging.GetLogger()

	logger.Info("create router")
	router := httprouter.New()

	logger.Info("register user handler")
	handler := user.NewHandler(logger)
	handler.Register(router)

	start(router)
}

func start(router *httprouter.Router) {
	logger := logging.GetLogger()
	logger.Info("start application")
	listener, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
	}

	logger.Info("server is listening 127.0.0.1:1234")
	logger.Fatal(server.Serve(listener))
}

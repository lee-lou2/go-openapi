package main

import (
	"go-openapi/cmd/api"
	"go-openapi/config"
	"log"
	"net/http"
)

func main() {
	// API 서버 실행
	server := api.Server()
	serverPort := config.GetEnv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "3000"
	}
	log.Fatal(http.ListenAndServe(":"+serverPort, server))
}

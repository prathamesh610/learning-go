package main

import (
	"log"

	"github.com/prathameshj610/go-microservices/internal/database"
	"github.com/prathameshj610/go-microservices/internal/server"
)


func main(){
	db, err := database.NewDataBaseClient()
	if err != nil {
		log.Fatalf("failed to initialize db client: %s", err)
	}

	srv := server.NewEchoServer(db)

	if err := srv.Start(); err != nil {
		log.Fatal(err.Error())
	}
}
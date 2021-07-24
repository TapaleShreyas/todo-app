package main

import (
	"go-server/router"
	"log"
	"net/http"
)

func main() {
	router := router.Router()
	log.Fatal(http.ListenAndServe(":8081", router))
}

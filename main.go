package main

import (
	"fmt"
	"github.com/pianisimo/go-postgresql-pq/router"
	"log"
	"net/http"
)

func main() {
	r := router.Router()
	fmt.Println("Starting server on port 8080....")

	log.Fatal(http.ListenAndServe(":8080", r))
}

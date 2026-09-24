package main

import (
	"log"

	"github.com/s444v/spots/internal/server"
)

func main() {
	s := server.CreateServer()
	log.Fatal(s.ListenAndServe())

}

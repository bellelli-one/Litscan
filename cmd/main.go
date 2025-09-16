package main

import (
	"log"
	"RIP/internal/api"
)

func main() {
	log.Println("Application start")
	api.StartServer()
	log.Println("Application stop")
}
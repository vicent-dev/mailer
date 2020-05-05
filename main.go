package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if len(os.Args) >= 2 {
		for i := 1; i < len(os.Args); i++ {
			sendEmail(os.Args[i], "Google no ha roto cosas", "Todo bien :D")
			log.Print("sent " + os.Args[i])
		}
		return
	}

	startApi()
}

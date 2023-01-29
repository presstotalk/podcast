package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/ninan.fm/podcast/rest"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = rest.Start()
	if err != nil {
		log.Fatalln(err)
	}
}

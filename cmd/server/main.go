package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/ninan.fm/podcast/rest"
)

func main() {
	godotenv.Load()

	err := rest.Start()
	if err != nil {
		log.Fatalln(err)
	}
}

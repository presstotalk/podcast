package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ninan.fm/podcast/feed"
)

func main() {
	if len(os.Args) <= 1 {
		help()
	}
	url := os.Args[1]

	podcast, err := feed.ParseURL(url)
	if err != nil {
		log.Fatal(err)
	}

	content, err := json.Marshal(podcast)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(content))
}

func help() {
	log.Fatalln("parse <feed-url>")
}

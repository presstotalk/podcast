package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ninan.fm/podcast/database"
	"github.com/ninan.fm/podcast/feed"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := database.New()
	podcastRepo := database.NewPodcastRepository(db)
	episodeRepo := database.NewEpisodeRepository(db)

	podcast, err := podcastRepo.Get(os.Getenv("AIRTABLE_BASE_ID"))
	if err != nil {
		log.Fatalln(err)
	}

	episodes, err := episodeRepo.ListAll(os.Getenv("AIRTABLE_BASE_ID"))
	if err != nil {
		log.Fatalln(err)
	}

	podcast.Episodes = episodes

	content, err := feed.Generate(podcast)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(content))
}

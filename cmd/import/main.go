package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ninan.fm/podcast/database"
	"github.com/ninan.fm/podcast/importer"
	"github.com/ninan.fm/podcast/storage"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if len(os.Args) <= 1 {
		help()
	}
	url := os.Args[1]

	db := database.New()

	imp := importer.New(&importer.Options{
		PodcastRepository: database.NewPodcastRepository(db),
		EpisodeRepository: database.NewEpisodeRepository(db),
		Storage:           storage.LocalStorage{Folder: os.Getenv("STORAGE_FOLDER")},
	})

	err = imp.Import(os.Getenv("AIRTABLE_BASE_ID"), url)
	if err != nil {
		log.Fatalln(err)
	}
}

func help() {
	log.Fatalln("parse <feed-url>")
}

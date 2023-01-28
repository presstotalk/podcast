package importer

import (
	"fmt"

	"log"

	pod "github.com/ninan.fm/podcast"
)

func (i Importer) uploadFiles(podcast *pod.Podacst) error {
	log.Println("Uploading...")
	return i.uploadPodcastFiles(podcast)
}

func (i Importer) uploadPodcastFiles(podcast *pod.Podacst) error {
	log.Println("> podcast cover image")
	// TODO: convert image format
	err := i.storage.UploadFromURL("cover.jpg", podcast.CoverImageURL)
	if err != nil {
		return err
	}

	for _, episode := range podcast.Episodes {
		err = i.uploadEpisodeFiles(episode)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i Importer) uploadEpisodeFiles(episode *pod.Episode) error {
	log.Printf("> %s > cover image", episode.Title)
	if episode.CoverImageURL != nil {
		err := i.storage.UploadFromURL(fmt.Sprintf("%s/cover.jpg", episode.Title), *episode.CoverImageURL)
		if err != nil {
			return err
		}
	}
	// TODO: support other auto formats
	log.Printf("> %s > audio file", episode.Title)
	return i.storage.UploadFromURL(fmt.Sprintf("%s/audio.mp3", episode.Title), *&episode.AudioURL)
}

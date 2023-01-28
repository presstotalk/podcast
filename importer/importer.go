package importer

import (
	"log"

	"github.com/ninan.fm/podcast/database"
	"github.com/ninan.fm/podcast/feed"
	"github.com/ninan.fm/podcast/storage"
)

type Importer struct {
	podcastRepository database.PodcastRepository
	episodeRepository database.EpisodeRepository
	storage           storage.Storage
}

type Options struct {
	PodcastRepository database.PodcastRepository
	EpisodeRepository database.EpisodeRepository
	Storage           storage.Storage
}

func New(options *Options) *Importer {
	return &Importer{
		options.PodcastRepository,
		options.EpisodeRepository,
		options.Storage,
	}
}

func (i Importer) Import(id string, feedURL string) error {
	log.Println("Parse rss feed")
	podcast, err := feed.ParseURL(feedURL)
	if err != nil {
		return err
	}

	log.Println("Save metadata into database")
	err = i.podcastRepository.Create(id, podcast)
	if err != nil {
		return err
	}

	for _, episode := range podcast.Episodes {
		err = i.episodeRepository.Create(id, episode)
		if err != nil {
			return err
		}
	}

	return nil
	//return i.uploadFiles(podcast)
}

package rest

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/ninan.fm/podcast/database"
	"github.com/ninan.fm/podcast/feed"
)

func Start() error {
	e := echo.New()

	e.GET("/feeds/podcast", func(c echo.Context) error {
		db := database.New()
		podcastRepo := database.NewPodcastRepository(db)
		episodeRepo := database.NewEpisodeRepository(db)

		podcast, err := podcastRepo.Get(os.Getenv("AIRTABLE_BASE_ID"))
		if err != nil {
			return err
		}

		episodes, err := episodeRepo.ListAll(os.Getenv("AIRTABLE_BASE_ID"))
		if err != nil {
			return err
		}

		podcast.Episodes = episodes

		content, err := feed.Generate(podcast)
		if err != nil {
			return err
		}

		c.Response().Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
		c.Response().Write(content)

		return nil
	})

	return e.Start(":" + os.Getenv("SERVER_PORT"))
}

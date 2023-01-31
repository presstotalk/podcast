package feed

import (
	"time"

	rss "github.com/eduncan911/podcast"
	"github.com/ninan.fm/podcast"
)

func Generate(podcast *podcast.Podacst) ([]byte, error) {
	now := time.Now()
	feed := rss.New(podcast.Title, podcast.URL, podcast.Description, nil, &now)

	for _, author := range podcast.Authors {
		feed.AddAuthor(author.Name, author.Email)
	}

	feed.AddImage(podcast.CoverImageURL)
	for _, category := range podcast.Categories {
		feed.AddCategory(category.Name, extractSubCategories(category.SubCategory))
	}

	for _, episode := range podcast.Episodes {
		item := rss.Item{
			GUID:        episode.GUID,
			Title:       episode.Title,
			Link:        episode.URL,
			Description: episode.Description,
			PubDate:     &episode.PublishedAt,
		}
		if episode.CoverImageURL != nil {
			item.AddImage(*episode.CoverImageURL)
		}
		item.AddEnclosure(episode.AudioURL, rss.MP3, episode.AudioSize)
		item.Author = &rss.Author{
			Name:  episode.Authors[0].Name,
			Email: episode.Authors[0].Email,
		}
		feed.AddItem(item)
	}

	return feed.Bytes(), nil
}

func extractSubCategories(category *podcast.Category) []string {
	categories := make([]string, 0)
	for true {
		if category == nil {
			break
		}
		categories = append(categories, category.Name)
		category = category.SubCategory
	}
	return categories
}

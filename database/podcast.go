package database

import (
	"errors"
	"os"

	"github.com/mehanizm/airtable"
	pod "github.com/ninan.fm/podcast"
	"github.com/samber/lo"
)

type PodcastRepository interface {
	Create(id string, podcast *pod.Podacst) error
	Get(id string) (*pod.Podacst, error)
}

func NewPodcastRepository(db *airtable.Client) PodcastRepository {
	return &podcastRepository{db}
}

type podcastRepository struct {
	db *airtable.Client
}

func (p *podcastRepository) Create(id string, podcast *pod.Podacst) error {
	table := p.db.GetTable(id, PODCAST_TABLE)
	_, err := table.AddRecords(&airtable.Records{
		Typecast: true,
		Records: []*airtable.Record{
			{
				Fields: map[string]interface{}{
					"Title":       podcast.Title,
					"Description": podcast.Description,
					"Authors": lo.Map(podcast.Authors, func(author *pod.Author, _ int) string {
						return author.String()
					}),
					"Categories": lo.Map(podcast.Categories, func(category *pod.Category, _ int) string {
						return category.String()
					}),
				},
			},
		},
	})
	return err
}

func (p *podcastRepository) Get(id string) (*pod.Podacst, error) {
	table := p.db.GetTable(id, PODCAST_TABLE)
	result, err := table.GetRecords().FromView("Grid view").ReturnFields("Title", "Authors", "Description", "Categories").Do()
	if err != nil {
		return nil, err
	}

	if len(result.Records) < 1 {
		return nil, errors.New("No podcast info")
	}

	record := result.Records[0]

	return &pod.Podacst{
		Title: record.Fields["Title"].(string),
		URL:   os.Getenv("PODCAST_URL"),
		Authors: lo.Map(record.Fields["Authors"].([]interface{}), func(val interface{}, _ int) *pod.Author {
			author := pod.NewAuthor(val.(string))
			return &author
		}),
		Description: record.Fields["Description"].(string),
		Categories: lo.Map(record.Fields["Categories"].([]interface{}), func(val interface{}, _ int) *pod.Category {
			category := pod.NewCategory(val.(string))
			return &category
		}),
		CoverImageURL: os.Getenv("ASSETS_BASE_URL") + "/cover.jpg",
	}, nil
}

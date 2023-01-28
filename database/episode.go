package database

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-module/carbon/v2"
	"github.com/mehanizm/airtable"
	pod "github.com/ninan.fm/podcast"
	"github.com/samber/lo"
)

type EpisodeRepository interface {
	Create(podcastId string, episode *pod.Episode) error
	ListAll(podcastId string) ([]*pod.Episode, error)
}

func NewEpisodeRepository(db *airtable.Client) EpisodeRepository {
	return &episodeRepository{db}
}

type episodeRepository struct {
	db *airtable.Client
}

func (r episodeRepository) Create(podcastId string, episode *pod.Episode) error {
	table := r.db.GetTable(podcastId, EPISODES_TABLE)
	_, err := table.AddRecords(&airtable.Records{
		Typecast: true,
		Records: []*airtable.Record{
			{
				Fields: map[string]interface{}{
					"GUID":        episode.GUID,
					"Title":       episode.Title,
					"Description": episode.Description,
					"Season":      episode.SeasonNumber,
					"Episode":     episode.EpisodeNumber,
					"Authors": lo.Map(episode.Authors, func(author *pod.Author, _ int) string {
						return author.String()
					}),
					"Publish Time": episode.PublishedAt,
					"Published":    true,
				},
			},
		},
	})
	return err
}

func (r episodeRepository) ListAll(podcastId string) ([]*pod.Episode, error) {
	table := r.db.GetTable(podcastId, EPISODES_TABLE)
	result, err := table.GetRecords().FromView("Grid view").WithFilterFormula("Published").ReturnFields("Title", "Authors", "Description", "Publish Time", "Published", "Season", "Episode", "GUID", "Slug").Do()
	if err != nil {
		return nil, err
	}

	if len(result.Records) < 1 {
		return nil, errors.New("No podcast info")
	}

	episodes := make([]*pod.Episode, len(result.Records))

	for i, record := range result.Records {
		coverImage := os.Getenv("ASSETS_BASE_URL") + "/cover.jpg"
		publishedAt := carbon.Parse(record.Fields["Publish Time"].(string))
		guid := getOptionalString(record.Fields["GUID"], record.ID)
		slug := getOptionalString(record.Fields["Slug"], guid)

		episodes[i] = &pod.Episode{
			GUID:  guid,
			Title: record.Fields["Title"].(string),
			URL:   fmt.Sprintf(os.Getenv("EPISODE_URL"), slug),
			Authors: lo.Map(record.Fields["Authors"].([]interface{}), func(val interface{}, _ int) *pod.Author {
				author := pod.NewAuthor(val.(string))
				return &author
			}),
			Description:   record.Fields["Description"].(string),
			PublishedAt:   publishedAt.ToStdTime(),
			SeasonNumber:  getOptionalString(record.Fields["Season"], ""),
			EpisodeNumber: getOptionalString(record.Fields["Episode"], ""),
			CoverImageURL: &coverImage,
			AudioURL:      os.Getenv("ASSETS_BASE_URL") + "/audio.mp3",
		}
	}

	return episodes, nil
}

func getOptionalString(val interface{}, defVal string) string {
	switch v := val.(type) {
	case string:
		return v
	default:
		return defVal
	}
}

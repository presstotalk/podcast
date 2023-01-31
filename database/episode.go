package database

import (
	"fmt"
	"net/http"
	"net/url"
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
	episodes := make([]*pod.Episode, 0)
	pageSize := 100
	offset := ""

	for true {
		result, err := table.
			GetRecords().
			FromView("Grid view").
			WithFilterFormula("Published").
			ReturnFields("Title", "Authors", "Description", "Publish Time", "Published", "Season", "Episode", "GUID", "Slug").
			WithOffset(offset).
			PageSize(pageSize).
			Do()
		if err != nil {
			return nil, err
		}

		if len(result.Records) < 1 {
			break
		}

		for _, record := range result.Records {
			guid := getOptionalString(record.Fields["GUID"], record.ID)
			slug := getOptionalString(record.Fields["Slug"], guid)
			coverImage := getEpisodeAssetUrl(slug, "cover.jpg")
			publishedAt := carbon.Parse(record.Fields["Publish Time"].(string))
			audioURL := getEpisodeAssetUrl(slug, "audio.mp3")

			episode := &pod.Episode{
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
				AudioURL:      audioURL,
				AudioSize:     getFileSizeFromURL(audioURL),
			}

			episodes = append(episodes, episode)
		}

		if result.Offset == "" {
			break
		}
		offset = result.Offset
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

func getEpisodeAssetUrl(slug, file string) string {
	return os.Getenv("ASSETS_BASE_URL") + "/podcast/" + url.PathEscape(slug) + "/" + file
}

func getFileSizeFromURL(url string) int64 {
	res, err := http.Head(url)
	if err != nil {
		return 0
	}
	return res.ContentLength
}

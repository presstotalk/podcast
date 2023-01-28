package feed

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/golang-module/carbon/v2"
	"github.com/mmcdole/gofeed"
	ext "github.com/mmcdole/gofeed/extensions"
	"github.com/ninan.fm/podcast"
)

func ParseURL(url string) (*podcast.Podacst, error) {
	fp := gofeed.NewParser()

	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, err
	}

	return parsePodcast(feed)
}

func parsePodcast(feed *gofeed.Feed) (*podcast.Podacst, error) {
	episodes := make([]*podcast.Episode, len(feed.Items))
	for i, item := range feed.Items {
		episode, err := parseItem(item)
		if err != nil {
			return nil, err
		}
		episodes[i] = episode
	}

	return &podcast.Podacst{
		Title:         feed.Title,
		URL:           feed.Link,
		Description:   feed.Description,
		Authors:       personsToAuthors(feed.Authors),
		CoverImageURL: feed.Image.URL,
		Episodes:      episodes,
		Categories:    parseCategories(feed),
	}, nil
}

func parseCategories(feed *gofeed.Feed) []*podcast.Category {
	if len(feed.ITunesExt.Categories) > 0 {
		categories := make([]*podcast.Category, len(feed.ITunesExt.Categories))
		for i, category := range feed.ITunesExt.Categories {
			categories[i] = parseItunesCategory(category)
		}
		return categories
	}

	categories := make([]*podcast.Category, len(feed.Categories))
	for i, category := range feed.Categories {
		categories[i] = &podcast.Category{
			Name: category,
		}
	}
	return categories
}

func parseItunesCategory(itunesCategory *ext.ITunesCategory) *podcast.Category {
	category := &podcast.Category{
		Name: itunesCategory.Text,
	}
	if itunesCategory.Subcategory != nil {
		category.SubCategory = parseItunesCategory(itunesCategory.Subcategory)
	}
	return category
}

func parseItem(item *gofeed.Item) (*podcast.Episode, error) {
	if len(item.Enclosures) <= 0 {
		return nil, errors.New("Missing audio url")
	}

	if match, err := regexp.MatchString("^audio/.+", item.Enclosures[0].Type); err != nil || !match {
		return nil, errors.New("Unsupported file type: " + item.Enclosures[0].Type)
	}

	audioSize, err := strconv.ParseInt(item.Enclosures[0].Length, 10, 64)
	if err != nil {
		return nil, err
	}

	return &podcast.Episode{
		GUID:          item.GUID,
		Title:         item.Title,
		URL:           item.Link,
		Description:   item.Description,
		Authors:       personsToAuthors(item.Authors),
		AudioURL:      item.Enclosures[0].URL,
		AudioSize:     audioSize,
		CoverImageURL: &item.Image.URL,
		EpisodeNumber: item.ITunesExt.Episode,
		SeasonNumber:  item.ITunesExt.Season,
		PublishedAt:   carbon.Parse(item.Published).ToStdTime(),
	}, nil
}

func personsToAuthors(persons []*gofeed.Person) []*podcast.Author {
	authors := make([]*podcast.Author, len(persons))
	for i, person := range persons {
		authors[i] = &podcast.Author{
			Name:  person.Name,
			Email: person.Email,
		}
	}
	return authors
}

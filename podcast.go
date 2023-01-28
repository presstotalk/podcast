package podcast

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"
)

type Author struct {
	Name  string
	Email string
}

func NewAuthor(author string) Author {
	re := regexp.MustCompile("^(.+?)(?:\\s*\\((.*?)\\))?$")
	match := re.FindStringSubmatch(author)
	return Author{match[1], match[2]}
}

func (a Author) String() string {
	if a.Email != "" {
		return fmt.Sprintf("%s (%s)", a.Name, a.Email)
	}
	return a.Name
}

type Category struct {
	Name        string
	SubCategory *Category
}

func NewCategory(val string) Category {
	list := []string{}
	json.Unmarshal([]byte(val), &list)

	mainCategory := Category{}
	category := &mainCategory

	for _, item := range list {
		category.Name = item
		category.SubCategory = &Category{}
		category = category.SubCategory
	}

	category.SubCategory = nil
	return mainCategory
}

func (c Category) String() string {
	list := make([]string, 0)
	curr := &c

	for true {
		if curr == nil {
			break
		}
		list = append(list, curr.Name)
		curr = curr.SubCategory
	}

	data, _ := json.Marshal(list)
	return string(data)
}

type Podacst struct {
	Title         string
	URL           string
	Description   string
	Authors       []*Author
	CoverImageURL string
	Episodes      []*Episode
	Categories    []*Category
}

type Episode struct {
	GUID          string
	Title         string
	URL           string
	Description   string
	AudioURL      string
	AudioSize     int64
	CoverImageURL *string
	EpisodeNumber string
	SeasonNumber  string
	Authors       []*Author
	PublishedAt   time.Time
}

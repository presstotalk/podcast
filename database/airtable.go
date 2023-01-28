package database

import (
	"os"

	"github.com/mehanizm/airtable"
)

const PODCAST_TABLE = "Podcast"
const EPISODES_TABLE = "Episodes"

func New() *airtable.Client {
	return airtable.NewClient(os.Getenv("AIRTABLE_API_KEY"))
}

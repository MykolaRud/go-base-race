package models

import "time"

type Article struct {
	Id          int
	SourceId    int
	URL         string
	Title       string
	CreatedAt   time.Time
	IsProcessed bool
}

package app

import (
	"time"
)

type ShortURL struct {
	ID       string    `bson:"_id"`
	URL      string    `bson:"url"`
	ExpireAt time.Time `bson:"expireAt,omitempty"`
}

type UrlID struct {
	ID   string `bson:"_id"`
	Used bool   `bson:"used,omitempty"`
}

package data

import "time"

type ShortenedUrlQ interface {
	GetByCode(code string) (*ShortenedUrl, error)
	Insert(value ShortenedUrl) (*ShortenedUrl, error)
}

type ShortenedUrl struct {
	ID        int64     `db:"id" structs:"-"`
	Code      string    `db:"code" structs:"code"`
	LongUrl   string    `db:"long_url" structs:"long_url"`
	CreatedAt time.Time `db:"created_at" structs:"created_at"`
}

package postgres

import (
	"github.com/bohdan-vykhovanets/url-shortener-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

func NewMainQ(db *pgdb.DB) data.MainQ {
	return &mainQ{
		db: db.Clone(),
	}
}

type mainQ struct {
	db *pgdb.DB
}

func (m *mainQ) New() data.MainQ {
	return NewMainQ(m.db)
}

func (m *mainQ) ShortenedUrl() data.ShortenedUrlQ {
	return newShortenedUrlQ(m.db)
}

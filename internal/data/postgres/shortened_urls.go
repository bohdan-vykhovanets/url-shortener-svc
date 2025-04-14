package postgres

import (
	"database/sql"
	er "errors"
	"github.com/fatih/structs"
	"github.com/lib/pq"
	"gitlab.com/distributed_lab/logan/v3/errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/bohdan-vykhovanets/url-shortener-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const tableName = "shortened_urls"

func newShortenedUrlQ(db *pgdb.DB) data.ShortenedUrlQ {
	return &shortenedUrlQ{
		db:  db,
		sql: sq.StatementBuilder,
	}
}

type shortenedUrlQ struct {
	db  *pgdb.DB
	sql sq.StatementBuilderType
}

func (q *shortenedUrlQ) GetByCode(code string) (*data.ShortenedUrl, error) {
	var result data.ShortenedUrl

	err := q.db.Get(&result, q.sql.Select("*").From(tableName).Where(sq.Eq{"code": code}))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get shortened url from database")
	}
	return &result, err
}

func (q *shortenedUrlQ) Insert(value data.ShortenedUrl) (*data.ShortenedUrl, error) {
	clauses := structs.Map(value)
	var result data.ShortenedUrl

	stmt := sq.Insert(tableName).SetMap(clauses).Suffix("RETURNING *")
	err := q.db.Get(&result, stmt)
	if err != nil {
		var pqErr *pq.Error
		if er.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, data.ErrCodeCollision
		}
		return nil, errors.Wrap(err, "failed to insert shortened url to database")
	}
	return &result, nil
}

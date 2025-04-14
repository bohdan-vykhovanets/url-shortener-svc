package data

import "gitlab.com/distributed_lab/logan/v3/errors"

var ErrCodeCollision = errors.New("code already exists")

type MainQ interface {
	New() MainQ
	ShortenedUrl() ShortenedUrlQ
}

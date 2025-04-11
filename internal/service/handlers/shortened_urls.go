package handlers

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/bohdan-vykhovanets/url-shortener-svc/internal/service/requests"
	"github.com/go-chi/chi"
	"github.com/google/jsonapi"
	"github.com/jxskiss/base62"
	"github.com/lib/pq"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/kit/pgdb"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type ShortenedUrl struct {
	ID        int64     `db:"id"`
	Code      string    `db:"code"`
	LongUrl   string    `db:"long_url"`
	CreatedAt time.Time `db:"created_at"`
}

type ShortenedUrls struct {
	db *pgdb.DB
}

func NewShortenedUrls(db *pgdb.DB) *ShortenedUrls {
	return &ShortenedUrls{db: db}
}

func (h *ShortenedUrls) CreateShortenedUrl(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewCreateUrl(r)
	if err != nil {
		Log(r).Error(err.Error())
		ape.RenderErr(w, &jsonapi.ErrorObject{
			Detail: err.Error(),
		})
		return
	}

	for i := 0; i < 5; i++ {

	}

	code, err := insertUniqueShortenedUrl(h.db, req.Url, 5)
	if err != nil {
		Log(r).Error(err.Error())
		ape.RenderErr(w, &jsonapi.ErrorObject{
			Detail: err.Error(),
		})
		return
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	shortenedUrl := fmt.Sprintf("%s://%s/%s", scheme, r.Host, code)

	ape.Render(w, shortenedUrl)
}

func (h *ShortenedUrls) Redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	query := "SELECT * FROM \"shortened_urls\" WHERE code = $1"
	var shortenedUrl ShortenedUrl
	err := h.db.Queryer.GetRaw(&shortenedUrl, query, code)
	if err != nil {
		Log(r).Error(err.Error())
		ape.RenderErr(w, &jsonapi.ErrorObject{
			Detail: fmt.Errorf("could not get shortened url: %w", err).Error(),
		})
		return
	}

	http.Redirect(w, r, shortenedUrl.LongUrl, http.StatusFound)
}

func generateShortenedUrlCode() (string, error) {
	b := make([]byte, 8)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	code := base62.EncodeToString(b)
	return code, nil
}

func insertUniqueShortenedUrl(db *pgdb.DB, longUrl string, attempts int) (string, error) {
	for i := 0; i < attempts; i++ {
		code, err := generateShortenedUrlCode()
		if err != nil {
			return "", err
		}

		query := "INSERT INTO \"shortened_urls\" (code, long_url, created_at) VALUES ($1, $2, $3)"
		err = db.Queryer.ExecRaw(query, code, longUrl, time.Now())
		if err != nil {
			var pgErr *pq.Error
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				continue
			}
		}

		return code, nil
	}

	return "", fmt.Errorf("failed to insert shortened url, code already exists")
}

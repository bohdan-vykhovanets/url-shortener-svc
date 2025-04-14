package handlers

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/bohdan-vykhovanets/url-shortener-svc/internal/data"
	"github.com/bohdan-vykhovanets/url-shortener-svc/internal/service/requests"
	"github.com/go-chi/chi"
	"github.com/google/jsonapi"
	"github.com/jxskiss/base62"
	"gitlab.com/distributed_lab/ape"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

func CreateShortenedUrl(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewCreateUrl(r)
	if err != nil {
		Log(r).Error(err.Error())
		ape.RenderErr(w, &jsonapi.ErrorObject{
			Detail: err.Error(),
		})
		return
	}

	shortenedUrl, err := insertUniqueShortenedUrl(Db(r), req.Url, 5)
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
	resultUrl := fmt.Sprintf("%s://%s/integrations/url-shortener-svc/urls/%s", scheme, r.Host, shortenedUrl.Code)

	ape.Render(w, resultUrl)
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	var shortenedUrl *data.ShortenedUrl
	shortenedUrl, err := Db(r).ShortenedUrl().GetByCode(code)
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

func insertUniqueShortenedUrl(db data.MainQ, longUrl string, attempts int) (*data.ShortenedUrl, error) {
	for i := 0; i < attempts; i++ {
		code, err := generateShortenedUrlCode()
		if err != nil {
			return nil, err
		}

		shortenedUrl := data.ShortenedUrl{
			Code:      code,
			LongUrl:   longUrl,
			CreatedAt: time.Now(),
		}

		res, err := db.ShortenedUrl().Insert(shortenedUrl)
		if errors.Is(err, data.ErrCodeCollision) {
			continue
		}

		return res, nil
	}

	return nil, fmt.Errorf("failed to insert shortened url, code already exists")
}

package service

import (
	"github.com/bohdan-vykhovanets/url-shortener-svc/internal/service/handlers"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"io"
	"net/http"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
		),
	)
	r.Route("/integrations/url-shortener-svc", func(r chi.Router) {
		r.Get("/urls/{code}", s.shortenedUrls.Redirect)
		r.Post("/urls/", s.shortenedUrls.CreateShortenedUrl)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			_, err := io.WriteString(w, "Health is OK")
			if err != nil {
				return
			}
		})
	})

	return r
}

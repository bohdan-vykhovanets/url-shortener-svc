package service

import (
	"github.com/bohdan-vykhovanets/url-shortener-svc/internal/config"
	"github.com/bohdan-vykhovanets/url-shortener-svc/internal/data/postgres"
	"github.com/bohdan-vykhovanets/url-shortener-svc/internal/service/handlers"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"io"
	"net/http"
)

func (s *service) router(cfg config.Config) chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxDb(postgres.NewMainQ(cfg.DB())),
		),
	)
	r.Route("/integrations/url-shortener-svc", func(r chi.Router) {
		r.Get("/urls/{code}", handlers.Redirect)
		r.Post("/urls/", handlers.CreateShortenedUrl)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			_, err := io.WriteString(w, "Health is OKay")
			if err != nil {
				return
			}
		})
	})

	return r
}

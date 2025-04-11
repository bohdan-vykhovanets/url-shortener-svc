package service

import (
	"github.com/bohdan-vykhovanets/url-shortener-svc/internal/service/handlers"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
	"gitlab.com/distributed_lab/ape"
	"io"
	"net/http"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()

	r.Use(cors.New(cors.Options{
		// All origins, or specify a domain: e.g. "http://localhost:3000"
		AllowedOrigins: []string{"*"},

		// Methods your API actually allows
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},

		// Headers you allow (e.g. "Content-Type" for JSON POST)
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},

		// Set to true if you need to send cookies / auth headers in requests
		AllowCredentials: false,

		// Preflight request cache duration
		MaxAge: 300,
	}).Handler)

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

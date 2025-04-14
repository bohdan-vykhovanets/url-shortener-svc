package handlers

import (
	"context"
	"github.com/bohdan-vykhovanets/url-shortener-svc/internal/data"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	dbCtxKey  ctxKey = iota
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxDb(entry data.MainQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context { return context.WithValue(ctx, dbCtxKey, entry) }
}

func Db(r *http.Request) data.MainQ { return r.Context().Value(dbCtxKey).(data.MainQ).New() }

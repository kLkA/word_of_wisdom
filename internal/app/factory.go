package app

import (
	"context"
	"net/http"

	"word_of_wisdom_pow/config"
	"word_of_wisdom_pow/internal/book"
)

type Factory struct {
	ctx context.Context
	cfg *config.Config

	httpServer *http.Server

	bookRepository *book.Repository
	bookService    book.Service
}

func NewFactory(ctx context.Context, cfg *config.Config) *Factory {
	return &Factory{
		ctx: ctx,
		cfg: cfg,
	}
}

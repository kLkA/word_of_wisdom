package app

import (
	"context"

	"github.com/rs/zerolog"
	"word_of_wisdom_pow/config"
	"word_of_wisdom_pow/internal/book"
	"word_of_wisdom_pow/internal/handshake"
	"word_of_wisdom_pow/internal/server"
)

type Application struct {
	cfg *config.Config
	log zerolog.Logger

	bookService *book.Service
	Factory     *Factory
	httpServer  *server.Server
}

func NewApplication(c *config.Config, log zerolog.Logger) (*Application, error) {
	ctx := context.Background()

	factory := NewFactory(ctx, c)
	bookService, err := factory.GetBookService()
	if err != nil {
		return nil, err
	}

	h := handshake.NewHandshake(c)
	httpServer := server.NewServer(c, log, h, bookService)
	app := Application{
		cfg:        c,
		log:        log,
		httpServer: httpServer,

		Factory: factory,
	}

	return &app, nil
}

func (app *Application) Run() error {
	app.log.Info().Msgf("running to listen http handler on %s", app.cfg.Server.ListenAddr)
	if err := app.httpServer.Start(); err != nil {
		return err
	}

	return nil
}

func (app *Application) Stop() error {
	if err := app.httpServer.Close(); err != nil {
		return err
	}
	return nil

}

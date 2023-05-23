package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"word_of_wisdom_pow/config"
	"word_of_wisdom_pow/internal/app"
)

func main() {
	configPath := flag.String("c", "config.toml", "Path to `toml` configuration file")
	flag.Parse()
	zerolog.DurationFieldUnit = time.Millisecond
	log := zerolog.New(&zerolog.ConsoleWriter{Out: os.Stdout}).
		Level(zerolog.TraceLevel).
		With().Timestamp().
		Logger()

	c, err := config.GetConfigFromFile(*configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't load cfg")
		return
	}

	application, err := app.NewApplication(c, log)
	if err != nil {
		log.Fatal().Err(err)
		return
	}

	go func() {
		if err := application.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("service run error")
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	// interrupt signal sent from terminal
	signal.Notify(quit, os.Interrupt)
	// sigterm signal sent from kubernetes
	signal.Notify(quit, syscall.SIGTERM)
	<-quit

	log.Info().Msg("api stop")
	if err := application.Stop(); err != nil {
		log.Err(err).Msg("api stop error")
		return
	}
}

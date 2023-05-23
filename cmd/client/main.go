package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"word_of_wisdom_pow/config"
	"word_of_wisdom_pow/internal/client"
	"word_of_wisdom_pow/internal/handshake"
)

func main() {
	configPath := flag.String("c", "config.toml", "Path to `toml` configuration file")
	flag.Parse()
	zerolog.DurationFieldUnit = time.Millisecond
	log := zerolog.New(&zerolog.ConsoleWriter{Out: os.Stdout}).
		Level(zerolog.TraceLevel).
		With().Timestamp().
		Logger()

	ctx, ctxCancel := context.WithCancel(context.TODO())
	defer ctxCancel()

	c, err := config.GetConfigFromFile(*configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't load cfg")
		return
	}
	// overwrite config via env var when run through docker
	if os.Getenv("LISTEN_ADDR") != "" {
		c.Server.ListenAddr = os.Getenv("LISTEN_ADDR")
	}
	h := handshake.NewHandshake(c)
	cl := client.NewClient(c, log, h)
	cl.StartFetchWorkers(ctx)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh
}

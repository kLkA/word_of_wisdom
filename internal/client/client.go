package client

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/rs/zerolog"
	"word_of_wisdom_pow/config"
	"word_of_wisdom_pow/internal/handshake"
)

type Client struct {
	handshake   *handshake.Handshake
	addr        string
	timeout     time.Duration
	workerCount int
	log         zerolog.Logger
}

func NewClient(c *config.Config, log zerolog.Logger, handshake *handshake.Handshake) *Client {
	timeout := time.Millisecond * time.Duration(c.Client.Timeout)

	return &Client{
		log:         log,
		addr:        c.Server.ListenAddr,
		handshake:   handshake,
		timeout:     timeout,
		workerCount: c.Client.WorkerCount,
	}
}

func (cl *Client) Request(log zerolog.Logger) bool {
	conn, err := net.Dial("tcp", cl.addr)
	if err != nil {
		log.Err(err).Msg("failed to connect")
		return false
	}

	calcDifficulty, calcDuration, err := cl.handshake.Connect(conn)
	if err != nil {
		conn.Close()
		log.Err(err).Msg("error occured")
		return false
	}

	response, err := io.ReadAll(conn)
	conn.Close()
	if err != nil {
		log.Error().Err(err).Msg("failed to read from connection")
		return false
	}
	if len(response) == 0 {
		log.Warn().Msg("empty response")
		return false
	}

	log.Info().
		Bytes("response", response).
		Int("difficulty", int(calcDifficulty)).
		Dur("pow_calc_duration", calcDuration).
		Msg("received response")
	return true
}

func (cl *Client) StartFetchWorkers(ctx context.Context) {
	creationPause := cl.timeout / time.Duration(cl.workerCount)
	for i := 0; i < cl.workerCount; i++ {
		go cl.runWorker(ctx, i)
		time.Sleep(creationPause)
	}
}

func (cl *Client) runWorker(ctx context.Context, workerID int) {
	log := cl.log.With().Int("worker_id", workerID).Logger()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if !cl.Request(log) {
				time.Sleep(cl.timeout)
			}
		}
	}
}

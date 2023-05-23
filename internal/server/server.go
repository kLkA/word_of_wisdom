package server

import (
	"errors"
	"io"
	"net"
	"strings"

	"github.com/rs/zerolog"
	"word_of_wisdom_pow/config"
	"word_of_wisdom_pow/internal/book"
	"word_of_wisdom_pow/internal/handshake"
)

type Server struct {
	addr        string
	listener    net.Listener
	log         zerolog.Logger
	handshake   *handshake.Handshake
	difficulty  int
	bookService book.Service
}

func NewServer(c *config.Config, log zerolog.Logger, handshake *handshake.Handshake, bookService book.Service) *Server {
	return &Server{
		log:         log,
		addr:        c.Server.ListenAddr,
		difficulty:  c.Handshake.Difficulty,
		handshake:   handshake,
		bookService: bookService,
	}
}

func (s *Server) Start() error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.listener = l

	go func() {
		connIdx := 0
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}
				s.log.Warn().Err(err).Msg("failed to listen socket")
				continue
			}
			connIdx++
			go s.serve(conn, connIdx)
		}
	}()
	return nil
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) serve(conn net.Conn, id int) {
	defer conn.Close()

	log := s.log.With().
		Int("id", id).
		Str("addr", conn.RemoteAddr().String()).
		Logger()
	log.Trace().Msg("receive conn")

	d, err := s.handshake.Serve(conn)
	if err != nil {
		log.Warn().Err(err).Dur("check_duration", d).Msg("refuse conn")
		return
	}
	log.Debug().
		Int("difficulty", s.difficulty).
		Dur("check_duration", d).
		Msg("pow is valid")

	s.handler(conn)
}

func (s *Server) handler(conn net.Conn) {
	b := s.bookService.GetRandomQuote()
	r := strings.NewReader(b.Quote)
	_, err := io.Copy(conn, r)
	if err != nil {
		s.log.Warn().Err(err).Msg("failed to write response")
	}
}

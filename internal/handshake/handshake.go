package handshake

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"net"
	"sync"
	"time"

	"word_of_wisdom_pow/config"
)

const (
	powHeaderSize = 3
	nonceSize     = 8
)

type Handshake struct {
	difficulty   byte
	powTokenSize int
	pool         *sync.Pool
}

func NewHandshake(c *config.Config) *Handshake {
	bufPool := &sync.Pool{
		New: func() interface{} {
			b := make([]byte, powHeaderSize+c.Handshake.PowTokenSize+nonceSize)
			return &b
		},
	}
	return &Handshake{
		pool:         bufPool,
		difficulty:   byte(c.Handshake.Difficulty),
		powTokenSize: c.Handshake.PowTokenSize,
	}
}

func (h *Handshake) Serve(conn net.Conn) (checkDuration time.Duration, err error) {
	bufPtr := h.pool.Get().(*[]byte)
	defer h.pool.Put(bufPtr)
	buf := *bufPtr

	resultOffset := powHeaderSize + h.powTokenSize

	puzzleBuf := buf[powHeaderSize:resultOffset]

	// read rand data
	_, err = rand.Read(puzzleBuf)
	if err != nil {
		return 0, errors.New("failed to read crypto rand")
	}

	buf[0] = h.difficulty
	binary.BigEndian.PutUint16(buf[1:], uint16(h.powTokenSize))

	// write puzzle packet
	if _, err = conn.Write(buf[:resultOffset]); err != nil {
		return 0, errors.New("failed to write PoW packet")
	}

	// read nonce packet answer
	_, err = conn.Read(buf[resultOffset:])
	if err != nil {
		return 0, errors.New("failed to read PoW packet")
	}

	// check proof
	beginCheck := time.Now()
	isValid := verifyBufPow(h.difficulty, buf[powHeaderSize:])
	checkDuration = time.Since(beginCheck)
	if !isValid {
		return checkDuration, errors.New("is not valid proof")
	}

	return checkDuration, nil
}

func (h *Handshake) Connect(conn net.Conn) (difficulty byte, duration time.Duration, err error) {
	buf := make([]byte, powHeaderSize)
	_, err = conn.Read(buf)
	if err != nil {
		return 0, 0, errors.New("failed to read PoW header")
	}

	difficulty = buf[0]
	tokenSize := binary.BigEndian.Uint16(buf[1:])

	buf = make([]byte, tokenSize)
	_, err = conn.Read(buf)
	if err != nil {
		err = errors.New("failed to read PoW data")
		return
	}

	beginCalc := time.Now()
	nonce, _, calcErr := calculatePow(difficulty, buf)
	duration = time.Since(beginCalc)
	if calcErr != nil {
		err = calcErr
		return
	}

	_, err = conn.Write(nonce)
	if err != nil {
		err = errors.New("failed to write nonce")
		return
	}

	return
}

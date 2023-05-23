package handshake

import (
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"word_of_wisdom_pow/config"
)

func TestHandshake(t *testing.T) {
	difficulty := 10
	powTokenSize := 64

	h := NewHandshake(&config.Config{
		Handshake: config.HandshakeSettings{
			Difficulty:   difficulty,
			PowTokenSize: powTokenSize,
		}})

	message := "hello"

	clientConn, serverConn := net.Pipe()

	go func() {
		_, receiveErr := h.Serve(serverConn)
		assert.NoError(t, receiveErr)
		serverConn.Write([]byte(message))
		serverConn.Close()
	}()

	receivedDifficulty, _, establishErr := h.Connect(clientConn)
	assert.NoError(t, establishErr)
	assert.Equal(t, difficulty, int(receivedDifficulty))

	receivedData, err := io.ReadAll(clientConn)
	assert.NoError(t, err)
	receivedMessage := string(receivedData)
	assert.Equal(t, message, receivedMessage)
}

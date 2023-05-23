package handshake

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcAndCheckPoW(t *testing.T) {
	var difficulty byte = 10

	data := make([]byte, 64)
	rand.Read(data)

	nonce, _, err := calculatePow(difficulty, data)
	assert.NoError(t, err)

	isValid := verifyPow(difficulty, data, nonce)
	assert.True(t, isValid)

	data[0], data[1] = data[1], data[0]
	isValid = verifyPow(difficulty, data, nonce)
	assert.False(t, isValid)
	data[0], data[1] = data[1], data[0]

	rand.Read(nonce)
	isValid = verifyPow(difficulty, data, nonce)
	assert.False(t, isValid)
}

func BenchmarkCalculatePow(b *testing.B) {
	data := make([]byte, 64)
	rand.Read(data)

	var difficulty byte = 5

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calculatePow(difficulty, data)
	}
}

func BenchmarkVerifyPow(b *testing.B) {
	data := make([]byte, 64)
	rand.Read(data)

	var difficulty byte = 10

	nonce, _, _ := calculatePow(difficulty, data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		verifyPow(difficulty, data, nonce)
	}
}

func BenchmarkVerifyBufPow(b *testing.B) {
	const tokenSize = 64
	data := make([]byte, tokenSize+nonceSize)
	rand.Read(data[:tokenSize])

	var difficulty byte = 10

	nonce, _, _ := calculatePow(difficulty, data)
	copy(data[tokenSize:], nonce)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		verifyBufPow(difficulty, data)
	}
}

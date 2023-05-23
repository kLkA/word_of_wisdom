package handshake

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"math"
)

func calculatePow(difficulty byte, data []byte) (powNonce, powHash []byte, err error) {
	nonceOffset := len(data)
	buf := make([]byte, nonceOffset+nonceSize)
	copy(buf, data)

	var hash [32]byte

	var nonce uint64
	for nonce < math.MaxUint64 {
		binary.BigEndian.PutUint64(buf[nonceOffset:], nonce)

		hash = sha256.Sum256(buf)

		if leadingZerosCount(hash[:]) >= difficulty {
			powNonce = buf[nonceOffset:]
			powHash = hash[:]
			return
		} else {
			nonce++
		}
	}

	err = errors.New("cannot to calculate proof hash")
	return
}

func verifyPow(difficulty byte, data []byte, powNonce []byte) bool {
	nonceOffset := len(data)
	buf := make([]byte, nonceOffset+nonceSize)
	copy(buf, data)
	copy(buf[nonceOffset:], powNonce)

	return verifyBufPow(difficulty, buf)
}

func verifyBufPow(difficulty byte, buf []byte) bool {
	hash := sha256.Sum256(buf)

	return leadingZerosCount(hash[:]) >= difficulty
}

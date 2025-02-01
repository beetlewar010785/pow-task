package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
)

type Nonce string

func (r Nonce) Bytes() []byte {
	return []byte(r)
}

func (r Nonce) String() string {
	return string(r)
}

func NonceFromBytes(b []byte) Nonce {
	return Nonce(b)
}

func NonceFromInt(i int) Nonce {
	return Nonce(strconv.Itoa(i))
}

type Difficulty int

func (r Difficulty) Int() int {
	return int(r)
}

type NonceFinder interface {
	Find(challenge Challenge, difficulty Difficulty) Nonce
}

type IncrementalNonceFinder struct {
}

func NewIncrementalNonceFinder() NonceFinder {
	return &IncrementalNonceFinder{}
}

func (r *IncrementalNonceFinder) Find(challenge Challenge, difficulty Difficulty) Nonce {
	nonce := 0
	for {
		guess := challenge.String() + strconv.Itoa(nonce)
		hash := sha256.Sum256([]byte(guess))
		hashHex := hex.EncodeToString(hash[:])

		if strings.HasPrefix(hashHex, strings.Repeat("0", difficulty.Int())) {
			return NonceFromInt(nonce)
		}
		nonce++
	}
}

package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"strings"
	"time"
)

type Challenge string

func (r Challenge) String() string {
	return string(r)
}

func (r Challenge) Bytes() []byte {
	return []byte(r)
}

func ChallengeFromBytes(b []byte) Challenge {
	return Challenge(b)
}

type ChallengeRandomizer interface {
	Generate(length int) Challenge
}

type SimpleChallengeRandomizer struct {
}

func NewSimpleChallengeRandomizer() *SimpleChallengeRandomizer {
	return &SimpleChallengeRandomizer{}
}

func (r *SimpleChallengeRandomizer) Generate(length int) Challenge {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rnd.Intn(len(charset))]
	}
	return ChallengeFromBytes(b)

}

type ChallengeVerifier interface {
	Verify(challenge Challenge, nonce Nonce, difficulty Difficulty) bool
}

type SimpleChallengeVerifier struct {
}

func NewSimpleChallengeVerifier() *SimpleChallengeVerifier {
	return &SimpleChallengeVerifier{}
}

func (r *SimpleChallengeVerifier) Verify(
	challenge Challenge,
	nonce Nonce,
	difficulty Difficulty,
) bool {
	guess := challenge.String() + nonce.String()
	hash := sha256.Sum256([]byte(guess))
	hashHex := hex.EncodeToString(hash[:])

	return strings.HasPrefix(hashHex, strings.Repeat("0", difficulty.Int()))
}

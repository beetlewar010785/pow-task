package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand/v2"
	"strconv"
	"strings"
)

type Challenge string
type Difficulty int

type ChallengeRandomizer interface {
	Generate() Challenge
}

type ChallengeVerifier interface {
	Verify(challenge Challenge, nonce Nonce, difficulty Difficulty) bool
}

type ASCIIChallengeRandomizer struct {
	challengeLength int
}

func NewASCIIChallengeRandomizer(
	challengeLength int,
) *ASCIIChallengeRandomizer {
	return &ASCIIChallengeRandomizer{
		challengeLength: challengeLength,
	}
}

func (r *ASCIIChallengeRandomizer) Generate() Challenge {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, r.challengeLength)
	for i := range b {
		index := rand.IntN(len(charset))
		b[i] = charset[index]
	}
	return Challenge(b)
}

type SHA256ChallengeVerifier struct {
}

func NewSHA256ChallengeVerifier() *SHA256ChallengeVerifier {
	return &SHA256ChallengeVerifier{}
}

func (r *SHA256ChallengeVerifier) Verify(
	challenge Challenge,
	nonce Nonce,
	difficulty Difficulty,
) bool {
	guess := string(challenge) + strconv.Itoa(int(nonce))
	hash := sha256.Sum256([]byte(guess))
	hashHex := hex.EncodeToString(hash[:])

	return strings.HasPrefix(hashHex, strings.Repeat("0", int(difficulty)))
}

package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Challenge string
type Nonce int

type Difficulty int

type ChallengeRandomizer interface {
	Generate() Challenge
}

type SimpleChallengeRandomizer struct {
	challengeLength int
}

func NewSimpleChallengeRandomizer(challengeLength int) *SimpleChallengeRandomizer {
	return &SimpleChallengeRandomizer{
		challengeLength: challengeLength,
	}
}

func (r *SimpleChallengeRandomizer) Generate() Challenge {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, r.challengeLength)
	for i := range b {
		b[i] = charset[rnd.Intn(len(charset))]
	}
	return Challenge(b)
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
	guess := string(challenge) + strconv.Itoa(int(nonce))
	hash := sha256.Sum256([]byte(guess))
	hashHex := hex.EncodeToString(hash[:])

	return strings.HasPrefix(hashHex, strings.Repeat("0", int(difficulty)))
}

type NonceFinder interface {
	Find(challenge Challenge, difficulty Difficulty) Nonce
}

type IncrementalNonceFinder struct {
	challengeVerifier ChallengeVerifier
}

func NewIncrementalNonceFinder(challengeVerifier ChallengeVerifier) NonceFinder {
	return &IncrementalNonceFinder{
		challengeVerifier,
	}
}

func (r *IncrementalNonceFinder) Find(challenge Challenge, difficulty Difficulty) Nonce {
	var nonce Nonce
	for {
		if r.challengeVerifier.Verify(challenge, nonce, difficulty) {
			return nonce
		}
		nonce++
	}
}

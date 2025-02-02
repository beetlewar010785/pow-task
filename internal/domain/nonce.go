package domain

import (
	"context"
	"fmt"
)

type Nonce int

type NonceFinder interface {
	Find(ctx context.Context, challenge Challenge, difficulty Difficulty) (Nonce, error)
}

type IncrementalNonceFinder struct {
	challengeVerifier ChallengeVerifier
}

func NewIncrementalNonceFinder(challengeVerifier ChallengeVerifier) *IncrementalNonceFinder {
	return &IncrementalNonceFinder{
		challengeVerifier,
	}
}

func (r *IncrementalNonceFinder) Find(ctx context.Context, challenge Challenge, difficulty Difficulty) (Nonce, error) {
	var nonce Nonce
	for {
		select {
		case <-ctx.Done():
			return 0, fmt.Errorf("nonce search timed out: %w", ctx.Err())
		default:
			if r.challengeVerifier.Verify(challenge, nonce, difficulty) {
				return nonce, nil
			}
			nonce++
		}
	}
}

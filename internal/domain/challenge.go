package domain

type Challenge string

func (r Challenge) Bytes() []byte {
	return []byte(r)
}

func ChallengeFromBytes(b []byte) Challenge {
	return Challenge(b)
}

type ChallengeRandomizer interface {
	Generate() Challenge
}

type SimpleChallengeRandomizer struct {
}

func NewSimpleChallengeRandomizer() *SimpleChallengeRandomizer {
	return &SimpleChallengeRandomizer{}
}

func (r *SimpleChallengeRandomizer) Generate() Challenge {
	return "abcdefghijklmnop"
}

type ChallengeVerifier interface {
	Verify(challenge Challenge, nonce Nonce) bool
}

type SimpleChallengeVerifier struct {
}

func NewSimpleChallengeVerifier() *SimpleChallengeVerifier {
	return &SimpleChallengeVerifier{}
}

func (r *SimpleChallengeVerifier) Verify(challenge Challenge, nonce Nonce) bool {
	return true
}

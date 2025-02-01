package domain

type Nonce string

func (r Nonce) Bytes() []byte {
	return []byte(r)
}

func NonceFromBytes(b []byte) Nonce {
	return Nonce(b)
}

type Difficulty int

type NonceFinder interface {
	Find(challenge Challenge, difficulty Difficulty) Nonce
}

type IncrementalNonceFinder struct {
}

func NewIncrementalNonceFinder() NonceFinder {
	return &IncrementalNonceFinder{}
}

func (r *IncrementalNonceFinder) Find(challenge Challenge, difficulty Difficulty) Nonce {
	return "123"
}

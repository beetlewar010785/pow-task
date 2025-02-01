package domain

type Grant string

func (r Grant) Bytes() []byte {
	return []byte(r)
}

type GrantProvider interface {
	Provide(nonce Nonce) Grant
}

type SimpleGrantProvider struct {
}

func NewSimpleGrantProvider() *SimpleGrantProvider {
	return &SimpleGrantProvider{}
}

func (g *SimpleGrantProvider) Provide(_ Nonce) Grant {
	return "some-grant"
}

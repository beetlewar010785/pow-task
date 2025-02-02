package domain

type GrantResult int

const (
	GrantResultSuccess GrantResult = iota
	GrantResultFailure
)

type Grant struct {
	Result GrantResult
	Quote  *Quote
}

func NewGrant(
	result GrantResult,
	quote *Quote,
) Grant {
	return Grant{
		result,
		quote,
	}
}

func SuccessGrant(quote Quote) Grant {
	return NewGrant(GrantResultSuccess, &quote)
}

func FailureGrant() Grant {
	return NewGrant(GrantResultFailure, nil)
}

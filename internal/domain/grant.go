package domain

type GrantResult int

const (
	GrantResultSuccess GrantResult = iota
	GrantResultFailure
)

type Grant struct {
	Result GrantResult
	Quote  *string
}

func NewGrant(
	result GrantResult,
	quote *string,
) Grant {
	return Grant{
		result,
		quote,
	}
}

func SuccessGrant(quote string) Grant {
	return NewGrant(GrantResultSuccess, &quote)
}

func FailureGrant() Grant {
	return NewGrant(GrantResultFailure, nil)
}

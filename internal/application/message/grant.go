package message

type GrantResult int

const (
	GrantResultSuccess GrantResult = iota
	GrantResultFailure
)

// TODO: move json attributes to adapter
type Grant struct {
	Result GrantResult `json:"result"`
	Quote  *string     `json:"quote,omitempty"`
}

func SuccessGrant(quote string) Grant {
	return Grant{
		Result: GrantResultSuccess,
		Quote:  &quote,
	}
}

func FailureGrant() Grant {
	return Grant{
		Result: GrantResultFailure,
	}
}

type GrantWriter interface {
	Write(grant Grant) error
}

type GrantReader interface {
	Read() (Grant, error)
}

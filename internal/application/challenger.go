package application

import "io"

type ChallengerFactory interface {
	Create(reader io.Reader, writer io.Writer) Challenger
}

type Challenger interface {
	Challenge() error
}

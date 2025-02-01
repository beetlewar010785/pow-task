package application

import "github.com/beetlewar010785/pow-task/internal/application/message"

type GrantReceiver interface {
	Receive() (message.Grant, error)
}

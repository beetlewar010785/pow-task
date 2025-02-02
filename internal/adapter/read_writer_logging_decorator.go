package adapter

import (
	"encoding/json"
	"fmt"
	"github.com/beetlewar010785/pow-task/internal/domain"
)

type ReadWriterLoggingDecorator struct {
	next   domain.ReadWriter
	logger domain.Logger
}

func NewReadWriterLoggingDecorator(
	next domain.ReadWriter,
	logger domain.Logger,
) *ReadWriterLoggingDecorator {
	return &ReadWriterLoggingDecorator{
		next,
		logger,
	}
}

func (r *ReadWriterLoggingDecorator) WritePOWRequest(pow domain.POWRequest) error {
	return writeDataAndLog(r.logger, pow, "POWRequest", r.next.WritePOWRequest)
}

func (r *ReadWriterLoggingDecorator) ReadPOWRequest() (domain.POWRequest, error) {
	return readDataAndLog(r.logger, "POWRequest", r.next.ReadPOWRequest)
}

func (r *ReadWriterLoggingDecorator) WritePOWResponse(pow domain.POWResponse) error {
	return writeDataAndLog(r.logger, pow, "POWResponse", r.next.WritePOWResponse)
}

func (r *ReadWriterLoggingDecorator) ReadPOWResponse() (domain.POWResponse, error) {
	return readDataAndLog(r.logger, "POWResponse", r.next.ReadPOWResponse)
}

func (r *ReadWriterLoggingDecorator) WriteGrant(grant domain.Grant) error {
	return writeDataAndLog(r.logger, grant, "Grant", r.next.WriteGrant)
}

func (r *ReadWriterLoggingDecorator) ReadGrant() (domain.Grant, error) {
	return readDataAndLog(r.logger, "Grant", r.next.ReadGrant)
}

type TReadWriterLoggerDecorator struct {
	domain.Logger
}

func readDataAndLog[T any](
	logger domain.Logger,
	dataName string,
	action func() (T, error),
) (T, error) {
	logger.Trace(fmt.Sprintf("%s reading", dataName))
	result, err := action()
	if err == nil {
		j, _ := json.MarshalIndent(result, "", "  ")
		logger.Trace(fmt.Sprintf("%s read succeeded: %s", dataName, j))
	}
	return result, err
}

func writeDataAndLog[T any](
	logger domain.Logger,
	data T,
	dataName string,
	action func(T) error,
) error {
	j, _ := json.MarshalIndent(data, "", "  ")
	logger.Trace(fmt.Sprintf("%s writing: %s", dataName, j))
	err := action(data)
	if err == nil {
		logger.Trace(fmt.Sprintf("%s write succeeded", dataName))
	}
	return err
}

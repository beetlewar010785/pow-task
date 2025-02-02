package domain

type ReadWriter interface {
	WritePOWRequest(pow POWRequest) error
	ReadPOWRequest() (POWRequest, error)

	WritePOWResponse(pow POWResponse) error
	ReadPOWResponse() (POWResponse, error)

	WriteGrant(grant Grant) error
	ReadGrant() (Grant, error)
}

type ReadWriterMock struct {
	POWRequest  POWRequest
	POWResponse POWResponse
	Grant       Grant
}

func (r *ReadWriterMock) WritePOWRequest(pow POWRequest) error {
	r.POWRequest = pow
	return nil
}

func (r *ReadWriterMock) ReadPOWRequest() (POWRequest, error) {
	return r.POWRequest, nil
}

func (r *ReadWriterMock) WritePOWResponse(pow POWResponse) error {
	r.POWResponse = pow
	return nil
}

func (r *ReadWriterMock) ReadPOWResponse() (POWResponse, error) {
	return r.POWResponse, nil
}

func (r *ReadWriterMock) WriteGrant(grant Grant) error {
	r.Grant = grant
	return nil
}

func (r *ReadWriterMock) ReadGrant() (Grant, error) {
	return r.Grant, nil
}

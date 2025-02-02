package domain

type ReadWriter interface {
	WritePowRequest(pow POWRequest) error
	ReadPowRequest() (POWRequest, error)

	WritePowResponse(pow POWResponse) error
	ReadPowResponse() (POWResponse, error)

	WriteGrant(grant Grant) error
	ReadGrant() (Grant, error)
}

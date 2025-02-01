package domain

type POWRequest struct {
	Challenge  Challenge
	Difficulty Difficulty
}

func NewPOWRequest(
	challenge Challenge,
	difficulty Difficulty,
) POWRequest {
	return POWRequest{
		challenge,
		difficulty,
	}
}

type POWResponse struct {
	Challenge Challenge
	Nonce     Nonce
}

func NewPOWResponse(
	challenge Challenge,
	nonce Nonce,
) POWResponse {
	return POWResponse{
		challenge,
		nonce,
	}
}

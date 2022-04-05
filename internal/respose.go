package internal

// CommonResponse is the base response for all the requests.
type CommonResponse struct {
	Err string `json:"err"`
}

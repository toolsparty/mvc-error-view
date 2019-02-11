package error_view

type response struct {
	Error *errResponse `json:"error,omitempty"`
}

type errResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

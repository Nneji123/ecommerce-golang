package leads

type ErrorResponse struct {
	Error string `json:"error"`
}

type JSONResponse struct {
	Message string `json:"message"`
}

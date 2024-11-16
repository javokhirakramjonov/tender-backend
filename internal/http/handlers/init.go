package handlers

type HTTPHandler struct {
	US string
}

func NewHandler() *HTTPHandler {
	return &HTTPHandler{
		US: "not implemented",
	}
}

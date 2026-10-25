package data

type ErrorResponse struct {
	Message string   `json:"message"`
	Items   []string `json:"items"`
}

package api

type errorResponse struct {
	Message string      `json:"message"`
	Code    int         `json:"status_code"`
	Params  interface{} `json:"params"`
}

package types

type APIResponse struct {
	Code    int    `json:"code"`
	Results Result `json:"results"`
	Status  string `json:"status"`
}

type Result struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
package types

type APIResponse struct {
	Code    int    `json:"code"`
	Results any    `json:"results"`
	Status  string `json:"status"`
}

type Result struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type LoginResult struct {
	Message string      `json:"message"`
	Token   interface{} `json:"token"`
}

type LogoutResult struct {
	Message string `json:"message"`
}

type ValidTokenResult struct {
	Message string `json:"message"`
}
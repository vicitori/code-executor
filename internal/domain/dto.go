package domain

type Request struct {
	Program  string `json:"program"`
	Compiler string `json:"compiler"`
}

type IdResponse struct {
	Id string `json:"id"`
}

type StatusResponse struct {
	Status TaskStatus `json:"status"`
}

type ResultResponse struct {
	Result string `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

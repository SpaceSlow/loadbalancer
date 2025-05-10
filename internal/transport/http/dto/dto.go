package dto

type Client struct {
	ClientID string  `json:"client_id"`
	Capacity float64 `json:"capacity"`
	RPS      float64 `json:"rate_per_sec"`
}

type CreateClientRequest struct {
	Client
}

type CreateClientResponse struct {
	ClientID string  `json:"client_id"`
	Capacity float64 `json:"capacity"`
	RPS      float64 `json:"rate_per_sec"`
	Token    string  `json:"token"`
}

type UpdateClientRequest struct {
	Capacity float64 `json:"capacity"`
	RPS      float64 `json:"rate_per_sec"`
}

type UpdateClientResponse struct {
	CreateClientResponse
}

type ListClientsResponse struct {
	Clients []Client `json:"clients"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

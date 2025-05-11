package dto

type Client struct {
	ClientID string  `json:"client_id"`
	Capacity float64 `json:"capacity"`
	RPS      float64 `json:"rate_per_sec"`
	APIKey   string  `json:"api_key"`
}

type CreateClientRequest struct {
	ClientID string  `json:"client_id"`
	Capacity float64 `json:"capacity"`
	RPS      float64 `json:"rate_per_sec"`
}

type CreateClientResponse struct {
	ClientID string  `json:"client_id"`
	Capacity float64 `json:"capacity"`
	RPS      float64 `json:"rate_per_sec"`
	APIKey   string  `json:"api_key"`
}

type UpdateClientRequest struct {
	Capacity float64 `json:"capacity"`
	RPS      float64 `json:"rate_per_sec"`
}

type UpdateClientResponse struct {
	ClientID string  `json:"client_id"`
	Capacity float64 `json:"capacity"`
	RPS      float64 `json:"rate_per_sec"`
	APIKey   string  `json:"api_key"`
}

type ListClientsResponse struct {
	Clients []Client `json:"clients"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

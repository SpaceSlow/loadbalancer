package dto

import "errors"

type Client struct {
	ClientID string  `json:"client_id"`
	Capacity float64 `json:"capacity"`
	RPS      float64 `json:"rate_per_sec"`
	APIKey   string  `json:"api_key"`
}

type CreateClientRequest struct {
	ClientID string   `json:"client_id"`
	Capacity *float64 `json:"capacity,omitempty"`
	RPS      *float64 `json:"rate_per_sec,omitempty"`
}

func (c CreateClientRequest) Validate() error {
	var errs []error
	if c.ClientID == "" {
		errs = append(errs, ErrClientIDValidation)
	}
	if c.Capacity != nil && *c.Capacity < 1 {
		errs = append(errs, ErrCapacityValidation)
	}
	if c.RPS != nil && *c.RPS <= 0 {
		errs = append(errs, ErrRPSValidation)
	}

	return errors.Join(errs...)
}

type CreateClientResponse struct {
	ClientID string  `json:"client_id"`
	Capacity float64 `json:"capacity"`
	RPS      float64 `json:"rate_per_sec"`
	APIKey   string  `json:"api_key"`
}

type UpdateClientRequest struct {
	Capacity *float64 `json:"capacity,omitempty"`
	RPS      *float64 `json:"rate_per_sec,omitempty"`
}

func (c UpdateClientRequest) Validate() error {
	var errs []error
	if c.Capacity == nil && c.RPS == nil {
		errs = append(errs, ErrNoSpecifiedUpdateRequestValidation)
	}
	if c.Capacity != nil && *c.Capacity < 1 {
		errs = append(errs, ErrCapacityValidation)
	}
	if c.RPS != nil && *c.RPS <= 0 {
		errs = append(errs, ErrRPSValidation)
	}

	return errors.Join(errs...)
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

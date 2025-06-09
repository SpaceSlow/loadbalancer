package dto

import "errors"

var ErrCapacityValidation = errors.New("capacity must be >= 1")
var ErrRPSValidation = errors.New("rate_per_sec must be > 0")
var ErrClientIDValidation = errors.New("client_id is not specified")
var ErrNoSpecifiedUpdateRequestValidation = errors.New("capacity and rate_per_sec are not specified")

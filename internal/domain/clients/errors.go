package clients

import "errors"

var ErrCapacityValidation = errors.New("validation error: capacity must be >= 1")
var ErrRPSValidation = errors.New("validation error: RPS must be > 0")
var ErrClientExists = errors.New("client already exists")
var ErrClientNotExists = errors.New("client not exists")
var ErrIncorrectClientAPIKey = errors.New("incorrect api key")

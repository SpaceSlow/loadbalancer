package clients

import "errors"

var ErrClientExists = errors.New("client already exists")
var ErrClientNotExists = errors.New("client not exists")
var ErrIncorrectClientAPIKey = errors.New("incorrect api key")

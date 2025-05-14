package clients

import (
	"fmt"
	"strings"
)

type Client struct {
	ID       string
	APIKey   string
	Capacity float64
	RPS      float64
}

func GenerateClientAPIKey(clientID string, capacity, rps float64) string {
	return fmt.Sprintf("%s-%0.2f-%0.2f", clientID, capacity, rps) // TODO: add strong algorithm
}

func ParseClientIDFromAPIKey(apiKey string) (string, error) {
	parts := strings.Split(apiKey, "-")
	if len(parts) < 3 {
		return "", ErrIncorrectClientAPIKey
	}
	return parts[0], nil
}

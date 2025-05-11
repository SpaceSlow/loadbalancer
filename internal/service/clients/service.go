package clients

import (
	"context"

	"github.com/SpaceSlow/loadbalancer/internal/domain/clients"
)

type Repository interface {
	Create(ctx context.Context, clientID, apiKey string, capacity, rps float64) (*clients.Client, error)
	List(ctx context.Context) ([]clients.Client, error)
	Fetch(ctx context.Context, clientID string) (*clients.Client, error)
	Update(ctx context.Context, clientID, newAPIKey string, newCapacity, newRPS float64) (*clients.Client, error)
	Delete(ctx context.Context, clientID string) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, clientID string, capacity, rps float64) (*clients.Client, error) {
	apiKey := clients.GenerateClientAPIKey(clientID, capacity, rps)
	return s.repo.Create(ctx, clientID, apiKey, capacity, rps)
}

func (s *Service) List(ctx context.Context) ([]clients.Client, error) {
	return s.repo.List(ctx)
}

func (s *Service) Fetch(ctx context.Context, clientID string) (*clients.Client, error) {
	return s.Fetch(ctx, clientID)
}

func (s *Service) Update(ctx context.Context, clientID string, newCapacity, newRPS float64) (*clients.Client, error) {
	newAPIKey := clients.GenerateClientAPIKey(clientID, newCapacity, newRPS)
	return s.repo.Update(ctx, clientID, newAPIKey, newCapacity, newRPS)
}

func (s *Service) Delete(ctx context.Context, clientID string) error {
	return s.repo.Delete(ctx, clientID)
}

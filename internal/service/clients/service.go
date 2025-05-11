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

type BucketService interface {
	AddBucket(clientID string, capacity, rps float64)
	UpdateBucket(clientID string, capacity, rps float64)
	DeleteBucket(clientID string)
}

type Service struct {
	repo          Repository
	bucketService BucketService
}

func NewService(repo Repository, bucketService BucketService) *Service {
	return &Service{repo: repo, bucketService: bucketService}
}

func (s *Service) Create(ctx context.Context, clientID string, capacity, rps float64) (*clients.Client, error) {
	apiKey := clients.GenerateClientAPIKey(clientID, capacity, rps)
	client, err := s.repo.Create(ctx, clientID, apiKey, capacity, rps)
	if err != nil {
		return nil, err
	}
	s.bucketService.AddBucket(client.ID, client.Capacity, client.RPS)
	return client, nil
}

func (s *Service) List(ctx context.Context) ([]clients.Client, error) {
	return s.repo.List(ctx)
}

func (s *Service) Fetch(ctx context.Context, clientID string) (*clients.Client, error) {
	return s.repo.Fetch(ctx, clientID)
}

func (s *Service) Update(ctx context.Context, clientID string, newCapacity, newRPS float64) (*clients.Client, error) {
	newAPIKey := clients.GenerateClientAPIKey(clientID, newCapacity, newRPS)
	client, err := s.repo.Update(ctx, clientID, newAPIKey, newCapacity, newRPS)
	if err != nil {
		return nil, err
	}
	s.bucketService.UpdateBucket(client.ID, client.Capacity, client.RPS)
	return client, nil
}

func (s *Service) Delete(ctx context.Context, clientID string) error {
	err := s.repo.Delete(ctx, clientID)
	if err != nil {
		return err
	}
	s.bucketService.DeleteBucket(clientID)
	return nil
}

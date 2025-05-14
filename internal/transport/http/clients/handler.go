package clients

import (
	"context"
	"net/http"

	"github.com/SpaceSlow/loadbalancer/internal/domain/clients"
	"github.com/SpaceSlow/loadbalancer/internal/transport/http/dto"
)

type Service interface {
	Create(ctx context.Context, clientID string, capacity, rps *float64) (*clients.Client, error)
	List(ctx context.Context) ([]clients.Client, error)
	Fetch(ctx context.Context, clientID string) (*clients.Client, error)
	Update(ctx context.Context, clientID string, newCapacity, newRPS *float64) (*clients.Client, error)
	Delete(ctx context.Context, clientID string) error
}

type Handlers struct {
	service Service
}

func NewHandlers(service Service) *Handlers {
	return &Handlers{service: service}
}

func (h *Handlers) Clients() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreateClient(w, r)
		case http.MethodGet:
			h.ListClients(w, r)
		default:
			dto.WriteErrorResponse(
				w,
				http.StatusMethodNotAllowed,
				"Supports only GET and POST methods",
			)
		}
	})
}

func (h *Handlers) ClientByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientID := r.PathValue("id")

		switch r.Method {
		case http.MethodGet:
			h.FetchClient(w, r, clientID)
		case http.MethodPut:
			h.UpdateClient(w, r, clientID)
		case http.MethodDelete:
			h.DeleteClient(w, r, clientID)
		default:
			dto.WriteErrorResponse(
				w,
				http.StatusMethodNotAllowed,
				"Supports only GET, PUT and DELETE methods",
			)
		}
	})
}

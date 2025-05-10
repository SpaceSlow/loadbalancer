package clients

import (
	"net/http"

	"github.com/SpaceSlow/loadbalancer/internal/domain/clients"
	"github.com/SpaceSlow/loadbalancer/internal/transport/http/dto"
)

type ClientService interface {
	Create(clientID string, capacity, rps float64) (*clients.Client, error)
	List() ([]clients.Client, error)
	Fetch(clientID string) (*clients.Client, error)
	Update(clientID string, newCapacity, newRPS float64) (*clients.Client, error)
	Delete(clientID string) error
}

type ClientHandlers struct {
	service ClientService
}

func NewClientHandlers(service ClientService) *ClientHandlers {
	return &ClientHandlers{service: service}
}

func (h *ClientHandlers) Clients() http.Handler {
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

func (h *ClientHandlers) ClientByID() http.Handler {
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

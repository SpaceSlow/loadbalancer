package clients

import (
	"log/slog"
	"net/http"

	"github.com/SpaceSlow/loadbalancer/internal/transport/http/dto"
	"github.com/SpaceSlow/loadbalancer/pkg/httpjson"
)

func (h *Handlers) ListClients(w http.ResponseWriter, r *http.Request) {
	clientSlice, err := h.service.List(r.Context())

	if err != nil {
		slog.Error(
			"List clients error",
			slog.String("error", err.Error()),
		)
		httpjson.WriteJSON(w, http.StatusInternalServerError, nil)
		return
	}

	responseClients := make([]dto.Client, 0, len(clientSlice))
	for _, c := range clientSlice {
		client := dto.Client{
			ClientID: c.ID,
			APIKey:   c.APIKey,
			Capacity: c.Capacity,
			RPS:      c.RPS,
		}
		responseClients = append(responseClients, client)
	}

	response := dto.ListClientsResponse{Clients: responseClients}

	httpjson.WriteJSON(w, http.StatusOK, response)
}

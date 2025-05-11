package clients

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/SpaceSlow/loadbalancer/internal/transport/http/dto"
	"github.com/SpaceSlow/loadbalancer/pkg/httpjson"
)

func (h *ClientHandlers) ListClients(w http.ResponseWriter, r *http.Request) {
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
			Capacity: c.Capacity,
			RPS:      c.RPS,
		}
		responseClients = append(responseClients, client)
	}

	response := dto.ListClientsResponse{Clients: responseClients}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		slog.Error(
			"Marshal response error",
			slog.String("path", r.RequestURI),
			slog.String("method", r.Method),
			slog.String("error", err.Error()),
		)
		httpjson.WriteJSON(w, http.StatusInternalServerError, nil)
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, responseJSON)
}

package clients

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/SpaceSlow/loadbalancer/internal/domain/clients"
	"github.com/SpaceSlow/loadbalancer/internal/transport/http/dto"
	"github.com/SpaceSlow/loadbalancer/pkg/httpjson"
)

func (h *Handlers) FetchClient(w http.ResponseWriter, r *http.Request, clientID string) {
	client, err := h.service.Fetch(r.Context(), clientID)
	if errors.Is(err, clients.ErrClientNotExists) {
		dto.WriteErrorResponse(w, http.StatusNotFound, "User with this username not exists")
		return
	} else if err != nil {
		slog.Error(
			"Fetch client error",
			slog.String("error", err.Error()),
			slog.String("client_id", clientID),
		)
		httpjson.WriteJSON(w, http.StatusInternalServerError, nil)
		return
	}

	response := dto.Client{
		ClientID: client.ID,
		APIKey:   client.APIKey,
		Capacity: client.Capacity,
		RPS:      client.RPS,
	}
	httpjson.WriteJSON(w, http.StatusOK, response)
}

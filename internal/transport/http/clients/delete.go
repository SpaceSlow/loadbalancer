package clients

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/SpaceSlow/loadbalancer/internal/domain/clients"
	"github.com/SpaceSlow/loadbalancer/internal/transport/http/dto"
	"github.com/SpaceSlow/loadbalancer/pkg/httpjson"
)

func (h *ClientHandlers) DeleteClient(w http.ResponseWriter, r *http.Request, clientID string) {
	err := h.service.Delete(r.Context(), clientID)
	if errors.Is(err, clients.ErrClientNotExists) {
		dto.WriteErrorResponse(w, http.StatusNotFound, "User with this username not exists")
		return
	} else if err != nil {
		slog.Error(
			"Delete client error",
			slog.String("error", err.Error()),
			slog.String("client_id", clientID),
		)
		httpjson.WriteJSON(w, http.StatusInternalServerError, nil)
		return
	}

	httpjson.WriteJSON(w, http.StatusNoContent, nil)
}

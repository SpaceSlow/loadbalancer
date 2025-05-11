package clients

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/SpaceSlow/loadbalancer/internal/domain/clients"
	"github.com/SpaceSlow/loadbalancer/internal/transport/http/dto"
	"github.com/SpaceSlow/loadbalancer/pkg/httpjson"
)

func (h *ClientHandlers) UpdateClient(w http.ResponseWriter, r *http.Request, clientID string) {
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		httpjson.WriteJSON(w, http.StatusInternalServerError, nil)
		return
	}

	var req dto.UpdateClientRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		dto.WriteErrorResponse(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Unmarshal request body error: %s", err.Error()),
		)
		return
	}

	client, err := h.service.Update(r.Context(), clientID, req.Capacity, req.RPS)
	if errors.Is(err, clients.ErrClientNotExists) {
		dto.WriteErrorResponse(w, http.StatusNotFound, "User with this username not exists")
		return
	} else if err != nil {
		slog.Error(
			"Update client error",
			slog.String("error", err.Error()),
			slog.String("client_id", clientID),
		)
		httpjson.WriteJSON(w, http.StatusInternalServerError, nil)
		return
	}

	response := dto.Client{
		ClientID: client.ID,
		Capacity: client.Capacity,
		RPS:      client.RPS,
	}
	httpjson.WriteJSON(w, http.StatusOK, response)
}

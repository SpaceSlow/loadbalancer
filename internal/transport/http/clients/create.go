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

func (h *Handlers) CreateClient(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		httpjson.WriteJSON(w, http.StatusInternalServerError, nil)
		return
	}

	var req dto.CreateClientRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		dto.WriteErrorResponse(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Unmarshal request body error: %s", err.Error()),
		)
		return
	}

	if err = req.Validate(); err != nil {
		dto.WriteErrorResponse(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Validation error: %s", err.Error()),
		)
		return
	}

	client, err := h.service.Create(r.Context(), req.ClientID, req.Capacity, req.RPS)

	if errors.Is(err, clients.ErrClientExists) {
		dto.WriteErrorResponse(w, http.StatusConflict, "Client with this client_id already exists")
		return
	} else if err != nil {
		slog.Error(
			"Client creation error",
			slog.String("error", err.Error()),
			slog.String("client_id", req.ClientID),
		)
		httpjson.WriteJSON(w, http.StatusInternalServerError, nil)
		return
	}

	response := dto.CreateClientResponse{
		ClientID: client.ID,
		APIKey:   client.APIKey,
		Capacity: client.Capacity,
		RPS:      client.RPS,
	}

	httpjson.WriteJSON(w, http.StatusCreated, response)
}

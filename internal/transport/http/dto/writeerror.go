package dto

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/SpaceSlow/loadbalancer/pkg/httpjson"
)

func WriteErrorResponse(writer http.ResponseWriter, code int, message string) {
	errorResponse := ErrorResponse{
		Code:    code,
		Message: message,
	}
	errorResponseJSON, err := json.Marshal(errorResponse)
	if err != nil {
		slog.Error("ErrorResponse json marshaling error", slog.String("error", err.Error()))
		httpjson.WriteJSON(writer, http.StatusInternalServerError, nil)
		return
	}
	httpjson.WriteJSON(writer, code, errorResponseJSON)
}

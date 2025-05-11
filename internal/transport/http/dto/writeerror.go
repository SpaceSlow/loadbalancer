package dto

import (
	"net/http"

	"github.com/SpaceSlow/loadbalancer/pkg/httpjson"
)

func WriteErrorResponse(writer http.ResponseWriter, code int, message string) {
	errorResponse := ErrorResponse{
		Code:    code,
		Message: message,
	}
	httpjson.WriteJSON(writer, code, errorResponse)
}

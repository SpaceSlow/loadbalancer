package httpjson

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

const (
	contentTypeKey  = "Content-Type"
	contentTypeJSON = "application/json"
)

func WriteJSON(writer http.ResponseWriter, code int, data any) {
	writer.Header().Set(contentTypeKey, contentTypeJSON)
	writer.WriteHeader(code)

	if data == nil {
		return
	}

	responseJSON, err := json.Marshal(data)
	if err != nil {
		slog.Error(
			"Marshal response error",
			slog.String("error", err.Error()),
		)
		return
	}

	writer.Write(responseJSON)
}

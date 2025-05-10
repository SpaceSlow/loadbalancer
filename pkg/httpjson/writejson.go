package httpjson

import "net/http"

const (
	contentTypeKey  = "Content-Type"
	contentTypeJSON = "application/json"
)

func WriteJSON(writer http.ResponseWriter, code int, data []byte) {
	writer.Header().Set(contentTypeKey, contentTypeJSON)
	writer.WriteHeader(code)
	writer.Write(data)
}

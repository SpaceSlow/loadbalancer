package router

import (
	"fmt"
	"net/http"

	"github.com/SpaceSlow/loadbalancer/internal/transport/http/clients"
)

func NewRouter(clientService clients.Service) *http.ServeMux {
	const basePathLayout = "/api/v1%s"

	mux := http.NewServeMux()

	clientHandlers := clients.NewHandlers(clientService)

	mux.Handle(fmt.Sprintf(basePathLayout, "/clients/"), clientHandlers.Clients())
	mux.Handle(fmt.Sprintf(basePathLayout, "/clients/{id}/"), clientHandlers.ClientByID())

	return mux
}

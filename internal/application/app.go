package application

import (
	"context"
	"net/http"

	"github.com/dwnGnL/ddos-pow/internal/service"
)

type Core interface {
	GetServer() service.ServerService
	GetClient() service.ClientService
}

func WithApp(app Core, f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ContextApp, app)
		r = r.WithContext(ctx)
		f.ServeHTTP(w, r)
	})
}

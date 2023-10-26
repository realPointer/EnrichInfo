package v1

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/realPointer/EnrichInfo/internal/service"
	"github.com/realPointer/EnrichInfo/pkg/logger"

	_ "github.com/realPointer/EnrichInfo/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(handler chi.Router, l logger.Interface, services *service.Services) {
	handler.Use(middleware.Logger)
	handler.Use(middleware.Recoverer)
	handler.Use(middleware.Timeout(60 * time.Second))

	handler.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong!"))
	})

	handler.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	handler.Route("/v1", func(r chi.Router) {
		r.Mount("/people", NewPeopleRouter(services.Person, l))
	})
}

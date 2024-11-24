package main

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/angrodrigo10/opentelemetry/serviceB/internal/infra/web/handlers"
	otelpkg "github.com/angrodrigo10/opentelemetry/serviceB/internal/pkg/otel"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

func main() {

	cleanup := otelpkg.InitTracer()
	defer cleanup()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(next, "request")
	})

	r.Get("/{cep}", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer("serviceB").Start(r.Context(), "handle-zipcode")
		defer span.End()

		cep := chi.URLParam(r, "cep")

		if cep == "" {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}
		if !handlers.IsValidCep(cep) {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}

		localidade, err := handlers.GetLocalidade(ctx, cep)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		encodedLocalidade := url.QueryEscape(localidade)

		currentWeather, err := handlers.GetTemperature(ctx, encodedLocalidade)
		if err != nil {
			http.Error(w, "Error Weather: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(currentWeather)
	})

	http.ListenAndServe(":8081", otelhttp.NewHandler(r, "server"))
}

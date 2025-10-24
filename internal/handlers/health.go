package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/M0F3/docstore-api/internal/middleware"
)

type HealthzResponse struct {
	Healthy bool `json:"healthy"`
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "Application/json")

	db, _ := middleware.GetDatabaseConnectionFromContext(r.Context())
	err := db.Ping(r.Context())
	if err != nil {
		json.NewEncoder(w).Encode(HealthzResponse{Healthy: false})
	}
	json.NewEncoder(w).Encode(HealthzResponse{Healthy: true})
}

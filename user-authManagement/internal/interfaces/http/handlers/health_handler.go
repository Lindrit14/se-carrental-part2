package handlers

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RabbitHealthcheck interface{ Healthy() bool }

type HealthHandler struct {
	mongo  *mongo.Client
	rabbit RabbitHealthcheck
}

func NewHealthHandler(m *mongo.Client, r RabbitHealthcheck) *HealthHandler {
	return &HealthHandler{mongo: m, rabbit: r}
}

// Liveness: process is running. Kept dependency-free.
func (h *HealthHandler) Liveness(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Readiness: dependencies are reachable.
func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	checks := map[string]string{}
	status := http.StatusOK

	if err := h.mongo.Ping(ctx, nil); err != nil {
		checks["mongo"] = "down: " + err.Error()
		status = http.StatusServiceUnavailable
	} else {
		checks["mongo"] = "ok"
	}
	if h.rabbit.Healthy() {
		checks["rabbitmq"] = "ok"
	} else {
		checks["rabbitmq"] = "down"
		status = http.StatusServiceUnavailable
	}
	writeJSON(w, status, checks)
}

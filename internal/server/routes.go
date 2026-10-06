package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/karbasia/doc-manager/internal/auth"
)

func NewRouter(logger *slog.Logger, tokens *auth.TokenManager) http.Handler {
	mux := http.NewServeMux()
	requireAuth := RequireAuth(tokens)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.Handle("GET /me", requireAuth(http.HandlerFunc(me)))

	return RequestLogger(logger)(mux)
}

func me(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.ClaimsFromContext(r.Context())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"user_id": claims.UserID()})
}

package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/urfave/negroni"
)

func NewServer() *chi.Mux {
	router := chi.NewRouter()
	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			h.ServeHTTP(w, r)
		})
	})
	router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nw := negroni.NewResponseWriter(w)
			start := time.Now()
			h.ServeHTTP(nw, r)
			diff := time.Since(start)
			proto := "http"
			if r.TLS != nil {
				proto = "https"
			}
			status := nw.Status()
			logM := slog.Info
			if status < 300 {
				logM = slog.Debug
			}
			logM(fmt.Sprintf("%s %s://%s%s", r.Method, proto, r.Host, r.RequestURI),
				"From", r.RemoteAddr, "Status", nw.Status(), "Time", fmt.Sprintf("%dms", diff.Milliseconds()))
		})
	})
	return router
}

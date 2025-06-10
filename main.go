package main

import (
	"log"
	"net/http"

	"sashstack/config"
	"sashstack/handlers"
)

var cfg = config.Load()

func main() {
	http.Handle("/", secureHeaders(http.HandlerFunc(handlers.Index)))
	http.Handle("/waitlist", secureHeaders(http.HandlerFunc(handlers.Waitlist)))
	http.Handle("/assets/", http.StripPrefix("/assets/", secureHeaders(http.FileServer(http.Dir("public/assets")))))

	appName := cfg.AppName
	port := cfg.Port
	log.Printf("%s is Running on Port:%s", appName, port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("❌ Server error: %v", err)
	}
}

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Basic hardening
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		w.Header().Set("Permissions-Policy", "geolocation=(), camera=(), microphone=()")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "no-referrer-when-downgrade")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; script-src 'self'")

		next.ServeHTTP(w, r)
	})
}

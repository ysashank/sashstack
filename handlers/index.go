package handlers

import (
	"net/http"
	"path/filepath"

	"sashstack/config"
)

var cfg = config.Load()

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, filepath.Join("templates", "html", "index.html"))
}

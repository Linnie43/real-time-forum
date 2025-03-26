package handlers

import (
	"net/http"
	"os"
	"path/filepath"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Add cache-busting headers
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Serve the index.html file
	wd, err := os.Getwd()
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, filepath.Join(wd, "frontend", "index.html"))
}

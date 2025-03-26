package handlers

import "net/http"

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/index.html")
}

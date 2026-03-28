package handlers

import (
	"encoding/json"
	"net/http"
)

func ok(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func fail(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func readJson(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
package render

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, body any) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, "error JSON encoding", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
}

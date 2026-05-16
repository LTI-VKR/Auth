package errors

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func WriteError(w http.ResponseWriter, r *http.Request, mapped MappedError) {
	if mapped.Body == nil {
		mapped = NewBasicMapped(http.StatusInternalServerError, "INTERNAL", "internal error")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(mapped.Status)

	mapped.Body.SetRequestID(middleware.GetReqID(r.Context()))
	_ = json.NewEncoder(w).Encode(mapped.Body)
}

package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func ParseJsonBody(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		payload := map[string]interface{}{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}

		ctx := context.WithValue(r.Context(), PayloadContextKey, payload)
		r = r.WithContext(ctx)

		next(w, r, ps)
	}
}

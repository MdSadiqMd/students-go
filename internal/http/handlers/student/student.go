package student

import "net/http"

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Creating StudentAPI is live 🥳"))
	}
}

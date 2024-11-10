package student

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/MdSadiqMd/students-go/internal/types"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student // defining type of incoming datga and holding it in a variable

		err := json.NewDecoder(r.Body).Decode(&student) // decoding incoming data
		if errors.Is(err, io.EOF) { // checking incoming data is type of Student or not

		}

		w.Write([]byte("Creating StudentAPI is live 🥳"))
	}
}

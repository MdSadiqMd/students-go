package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/MdSadiqMd/students-go/internal/types"
	"github.com/MdSadiqMd/students-go/internal/utils/response"
	"github.com/go-playground/validator"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student // defining type of incoming data and holding it in a variable

		err := json.NewDecoder(r.Body).Decode(&student) // decoding incoming data
		if errors.Is(err, io.EOF) {                     // checking incoming data is type of Student or not
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("request body is empty")))
			return
		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// validating requests using validator package
		if err := validator.New().Struct(student); err != nil {
			validateError := err.(validator.ValidationErrors) // Type cast validation error
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateError))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]interface{}{"data": student, "success": "OK"})
	}
}

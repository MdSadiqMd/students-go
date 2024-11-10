package response

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
)

type Response struct {
	Status string `json:"status"` // when this gets serialize in json i want the feild name to be status not as Status as if we initially give name as Status we cannot use feild in pther packages
	Error  string `json:"error"`
}

const (
	StatusOK    = "OK"
	StatusError = "ERROR"
)

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidationError(err validator.ValidationErrors) Response {
	var errorMessages []string
	for _, err := range err {
		switch err.ActualTag() {
		case "required":
			errorMessages = append(errorMessages, err.Field()+" is required")
		default:
			errorMessages = append(errorMessages, err.Field()+" is invalid")
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errorMessages, ", "),
	}
}

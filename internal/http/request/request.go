package request

import (
	"encoding/json"
	"net/http"

	"github.com/aveloper/blog/internal/http/response"
	"github.com/aveloper/blog/internal/validator"

	"go.uber.org/zap"
)

// Reader has functions to read, parse nad validate request data
type Reader struct {
	log       *zap.Logger
	jw        *response.JSONWriter
	validator *validator.Validator
}

// NewReader returns a new instance of Reader
func NewReader(log *zap.Logger) *Reader {
	return &Reader{
		log: log,
	}
}

// ReadJSONAndValidate reads a json request body into the given struct and the validates the struct data
func (read *Reader) ReadJSONAndValidate(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	err := read.ReadJSONRequest(r, &v)
	if err != nil {
		parseErr := read.HandleParseError(err)
		read.jw.BadRequest(w, r, parseErr)
		return false
	}

	ve := read.validate(v)
	if ve != nil {
		read.jw.UnprocessableEntity(w, r, ve)
		return false
	}

	return true
}

// ReadJSONRequest reads a json request body into the given struct
func (read *Reader) ReadJSONRequest(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (read *Reader) validate(v interface{}) response.APIError {
	fields := read.validator.IsValidStruct(v)
	if fields == nil || len(fields) == 0 {
		return nil
	}

	return &ValidationError{
		ErrData: fields,
	}
}

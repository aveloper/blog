package response

import (
	"encoding/json"
	"github.com/aveloper/blog/internal/blogcontext"
	"net/http"
	"time"

	"go.uber.org/zap"
)

//JSONWriter has functions for sending JSON response to client
type JSONWriter struct {
	log *zap.Logger
}

//NewJSONWriter creates a new instance of JSONWriter
func NewJSONWriter(log *zap.Logger) *JSONWriter {
	return &JSONWriter{
		log: log,
	}
}

//Ok sends the data to client with http status 200
func (j *JSONWriter) Ok(w http.ResponseWriter, r *http.Request, data interface{}) {
	res := j.buildResponse(r, true, data, nil)
	j.jsonWrite(w, res, http.StatusOK)
}

// Error sends error to client with the given http status
func (j *JSONWriter) Error(w http.ResponseWriter, r *http.Request, apiError APIError, httpStatus int) {
	res := j.buildResponse(r, false, nil, j.buildError(apiError))
	j.jsonWrite(w, res, httpStatus)
}

//NotFound sends error to client with http status 404
func (j *JSONWriter) NotFound(w http.ResponseWriter, r *http.Request, apiError APIError) {
	j.Error(w, r, apiError, http.StatusNotFound)
}

//Unauthorized sends error to client with http status 401
func (j *JSONWriter) Unauthorized(w http.ResponseWriter, r *http.Request, apiError APIError) {
	j.Error(w, r, apiError, http.StatusUnauthorized)
}

//Forbidden sends error to client with http status 403
func (j *JSONWriter) Forbidden(w http.ResponseWriter, r *http.Request, apiError APIError) {
	j.Error(w, r, apiError, http.StatusForbidden)
}

//UnprocessableEntity sends error to client with http status 422
func (j *JSONWriter) UnprocessableEntity(w http.ResponseWriter, r *http.Request, apiError APIError) {
	j.Error(w, r, apiError, http.StatusUnprocessableEntity)
}

//BadRequest sends error to client with http status 400
func (j *JSONWriter) BadRequest(w http.ResponseWriter, r *http.Request, apiError APIError) {
	j.Error(w, r, apiError, http.StatusBadRequest)
}

//Internal sends error to client with http status 500
func (j *JSONWriter) Internal(w http.ResponseWriter, r *http.Request, apiError APIError) {
	j.Error(w, r, apiError, http.StatusInternalServerError)
}

//DefaultError sends unknown error to client with http status 500
func (j *JSONWriter) DefaultError(w http.ResponseWriter, r *http.Request) {
	j.Error(w, r, &defaultErr{}, http.StatusInternalServerError)
}

//buildResponse builds the response object to be sent to client
func (j *JSONWriter) buildResponse(r *http.Request, success bool, data interface{}, ew *errorWrap) *response {
	return &response{
		RequestID: j.getRequestID(r),
		Timestamp: time.Now(),
		URI:       r.RequestURI,
		Success:   success,
		Data:      data,
		Error:     ew,
	}
}

//buildError builds the error to be sent to client
func (j *JSONWriter) buildError(apiError APIError) *errorWrap {
	return &errorWrap{
		Message: apiError.Message(),
		Code:    apiError.Code(),
		Data:    apiError.Data(),
	}
}

//jsonWrite writes the data to client as a JSON
func (j *JSONWriter) jsonWrite(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		j.log.Panic("failed encoding json for HTTP response",
			zap.Error(err),
			zap.Any("data", data),
			zap.Int("status_code", statusCode),
		)
	}
}

func (j *JSONWriter) getRequestID(r *http.Request) string {
	id, err := blogcontext.GetRequestID(r.Context())
	if err != nil {
		j.log.Error("failed to read request-id", zap.Error(err))
		return ""
	}

	return id
}

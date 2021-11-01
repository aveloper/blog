package response

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type JSONWriter struct {
	log *zap.Logger
}

//NewJSONWriter creates a new instance of JSONWriter
func NewJSONWriter(log *zap.Logger) *JSONWriter {
	return &JSONWriter{
		log: log,
	}
}

func (j *JSONWriter) Ok(w http.ResponseWriter, r *http.Request, data interface{}) {
	res := j.buildResponse(r, true, data, nil)
	j.jsonWrite(w, res, http.StatusOK)
}

func (j *JSONWriter) Error(w http.ResponseWriter, r *http.Request, apiError APIError, httpStatus int) {
	res := j.buildResponse(r, false, nil, j.buildError(apiError))
	j.jsonWrite(w, res, httpStatus)
}

func (j *JSONWriter) NotFound(w http.ResponseWriter, r *http.Request, apiError APIError) {
	j.Error(w, r, apiError, http.StatusNotFound)
}

func (j *JSONWriter) Unauthorized(w http.ResponseWriter, r *http.Request, apiError APIError) {
	j.Error(w, r, apiError, http.StatusUnauthorized)
}

func (j *JSONWriter) Forbidden(w http.ResponseWriter, r *http.Request, apiError APIError) {
	j.Error(w, r, apiError, http.StatusForbidden)
}

func (j *JSONWriter) UnprocessableEntity(w http.ResponseWriter, r *http.Request, apiError APIError) {
	j.Error(w, r, apiError, http.StatusUnprocessableEntity)
}

func (j *JSONWriter) BadRequest(w http.ResponseWriter, r *http.Request, apiError APIError) {
	j.Error(w, r, apiError, http.StatusBadRequest)
}

func (j *JSONWriter) Internal(w http.ResponseWriter, r *http.Request, apiError APIError) {
	j.Error(w, r, apiError, http.StatusInternalServerError)
}

func (j *JSONWriter) DefaultError(w http.ResponseWriter, r *http.Request) {
	j.Error(w, r, &defaultErr{}, http.StatusInternalServerError)
}

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

func (j *JSONWriter) buildError(apiError APIError) *errorWrap {
	return &errorWrap{
		Message: apiError.Message(),
		Code:    apiError.Code(),
		Data:    apiError.Data(),
	}
}

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
	return ""
}

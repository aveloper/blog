package response

import (
	"github.com/aveloper/blog/internal/blogcontext"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	log = zap.NewExample()
)

const (
	requestID = "dummy-id"

	testErrorMsg  = "error occurred"
	testErrorCode = DefaultErrorCode
	testErrorData = 1
)

type testError struct{}

func (t *testError) Message() string {
	return testErrorMsg
}

func (t *testError) Code() ErrorCode {
	return testErrorCode
}

func (t *testError) Data() interface{} {
	return testErrorData
}

func TestJSONWriter_Ok(t *testing.T) {
	jw := NewJSONWriter(log)
	rr, r := getResponseRequest()

	jw.Ok(rr, r, nil)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestJSONWriter_BadRequest(t *testing.T) {
	jw := NewJSONWriter(log)
	rr, r := getResponseRequest()

	jw.BadRequest(rr, r, &testError{})

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestJSONWriter_DefaultError(t *testing.T) {
	jw := NewJSONWriter(log)
	rr, r := getResponseRequest()

	jw.DefaultError(rr, r)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestJSONWriter_Forbidden(t *testing.T) {
	jw := NewJSONWriter(log)
	rr, r := getResponseRequest()

	jw.Forbidden(rr, r, &testError{})

	assert.Equal(t, http.StatusForbidden, rr.Code)
}

func TestJSONWriter_Internal(t *testing.T) {
	jw := NewJSONWriter(log)
	rr, r := getResponseRequest()

	jw.Internal(rr, r, &testError{})

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestJSONWriter_NotFound(t *testing.T) {
	jw := NewJSONWriter(log)
	rr, r := getResponseRequest()

	jw.NotFound(rr, r, &testError{})

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestJSONWriter_Unauthorized(t *testing.T) {
	jw := NewJSONWriter(log)
	rr, r := getResponseRequest()

	jw.Unauthorized(rr, r, &testError{})

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestJSONWriter_UnprocessableEntity(t *testing.T) {
	jw := NewJSONWriter(log)
	rr, r := getResponseRequest()

	jw.UnprocessableEntity(rr, r, &testError{})

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestJSONWriter_buildError(t *testing.T) {
	jw := NewJSONWriter(log)

	err := jw.buildError(&defaultErr{})

	de := &defaultErr{}

	assert.NotNil(t, err)
	assert.Equal(t, de.Data(), err.Data)
	assert.Equal(t, de.Message(), err.Message)
	assert.Equal(t, de.Code(), err.Code)
}

func TestJSONWriter_buildResponse(t *testing.T) {
	jw := NewJSONWriter(log)

	t.Run("test success response", func(t *testing.T) {
		r := getRequest()

		s := time.Now()

		res := jw.buildResponse(r, true, 1, nil)

		e := time.Now()

		assert.Nil(t, res.Error)
		assert.Equal(t, true, res.Success)
		assert.Equal(t, 1, res.Data)
		assert.Equal(t, requestID, res.RequestID)
		assert.Equal(t, "/", res.URI)
		assert.True(t, res.Timestamp.After(s) && res.Timestamp.Before(e))
	})

	t.Run("test failure response", func(t *testing.T) {
		r := getRequest()

		s := time.Now()

		res := jw.buildResponse(r, false, nil, jw.buildError(&testError{}))

		e := time.Now()

		assert.Nil(t, res.Data)
		assert.Equal(t, false, res.Success)
		assert.Equal(t, requestID, res.RequestID)
		assert.Equal(t, "/", res.URI)
		assert.True(t, res.Timestamp.After(s) && res.Timestamp.Before(e))
	})

}

func TestJSONWriter_getRequestID(t *testing.T) {
	jw := NewJSONWriter(log)
	_, r := getResponseRequest()

	id := jw.getRequestID(r)

	assert.Equal(t, requestID, id)
}

func TestJSONWriter_jsonWrite(t *testing.T) {
	jw := NewJSONWriter(log)
	rr, _ := getResponseRequest()

	data := map[string]string{
		"hello": "world",
	}

	jw.jsonWrite(rr, data, http.StatusOK)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func getResponseRequest() (*httptest.ResponseRecorder, *http.Request) {
	w := httptest.NewRecorder()
	return w, getRequest()
}

func getRequest() *http.Request {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := blogcontext.AddRequestID(r.Context(), requestID)
	return r.WithContext(ctx)
}

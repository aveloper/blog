package request

import (
	"bytes"
	"encoding/json"
	"github.com/aveloper/blog/internal/http/response"
	"github.com/aveloper/blog/internal/validator"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testRequest struct {
	Field1 string `json:"field_1" validate:"required"`
	Field2 string `json:"field_2" validate:"email"`
}

func TestReader_ReadJSONAndValidate(t *testing.T) {
	log := zaptest.NewLogger(t)
	jw := response.NewJSONWriter(log)
	reader := NewReader(log, jw, validator.New(log))

	t.Run("test that reading request fails with incorrect data", func(t *testing.T) {
		input := &testRequest{
			Field1: "1",
			Field2: "test@wrong-email",
		}

		inputJson, err := json.Marshal(input)
		assert.Nil(t, err)

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/some/endpoint", bytes.NewReader(inputJson))

		output := &testRequest{}

		ok := reader.ReadJSONAndValidate(rr, req, output)
		assert.False(t, ok)
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("test that reading request fails with invalid data", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/some/endpoint", bytes.NewReader([]byte("hello")))

		output := &testRequest{}

		ok := reader.ReadJSONAndValidate(rr, req, output)
		assert.False(t, ok)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("test that reading request fails with no data", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/some/endpoint", bytes.NewReader([]byte("")))

		output := &testRequest{}

		ok := reader.ReadJSONAndValidate(rr, req, output)
		assert.False(t, ok)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("test that reading request fails with invalid field type", func(t *testing.T) {
		input := map[string]interface{}{
			"field_1": 1,
			"field_2": "test@mail.com",
		}

		inputJson, err := json.Marshal(input)
		assert.Nil(t, err)

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/some/endpoint", bytes.NewReader(inputJson))

		output := &testRequest{}

		ok := reader.ReadJSONAndValidate(rr, req, output)
		assert.False(t, ok)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("test that reading request fails with unknown field", func(t *testing.T) {
		input := map[string]interface{}{
			"field_1": "hello",
			"field_3": "test@mail.com",
		}

		inputJson, err := json.Marshal(input)
		assert.Nil(t, err)

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/some/endpoint", bytes.NewReader(inputJson))

		output := &testRequest{}

		ok := reader.ReadJSONAndValidate(rr, req, output)
		assert.False(t, ok)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("test that reading request succeeds with correct data", func(t *testing.T) {
		input := &testRequest{
			Field1: "1",
			Field2: "test@mail.com",
		}

		inputJson, err := json.Marshal(input)
		assert.Nil(t, err)

		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/some/endpoint", bytes.NewReader(inputJson))

		output := &testRequest{}

		ok := reader.ReadJSONAndValidate(rr, req, output)
		assert.True(t, ok)
	})

}

func TestReader_ReadJSONRequest(t *testing.T) {
	log := zaptest.NewLogger(t)
	jw := response.NewJSONWriter(log)
	reader := NewReader(log, jw, validator.New(log))

	input := &testRequest{
		Field1: "1",
		Field2: "2",
	}

	inputJson, err := json.Marshal(input)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/some/endpoint", bytes.NewReader(inputJson))

	output := &testRequest{}

	err = reader.ReadJSONRequest(req, output)
	assert.Nil(t, err)

	assert.Equal(t, input, output)

}

func TestReader_validate(t *testing.T) {
	log := zaptest.NewLogger(t)
	jw := response.NewJSONWriter(log)
	reader := NewReader(log, jw, validator.New(log))

	t.Run("test that validation fails with incorrect data", func(t *testing.T) {
		v := &testRequest{
			Field1: "",
			Field2: "a@dlclec",
		}

		err := reader.validate(v)
		assert.NotNil(t, err)
		assert.Len(t, err.Data(), 2)
		assert.Equal(t, err.Code(), response.ValidationFailed)
		assert.Equal(t, err.Data(), []string{"field_1", "field_2"})
	})

	t.Run("test that validation succeeds with correct data", func(t *testing.T) {
		v := &testRequest{
			Field1: "abcd",
			Field2: "test@mail.com",
		}

		err := reader.validate(v)
		assert.Nil(t, err)
	})
}

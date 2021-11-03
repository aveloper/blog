package users

import (
	"net/http"
	"time"

	"github.com/aveloper/blog/internal/db/query"

	"github.com/aveloper/blog/internal/http/request"
	"github.com/aveloper/blog/internal/http/response"
	"go.uber.org/zap"
)

//Handler has http handler functions for user APIs
type Handler struct {
	log        *zap.Logger
	reader     *request.Reader
	jsonWriter *response.JSONWriter
	repository *Repository
}

//NewHandler creates a new instance of Handler
func NewHandler(log *zap.Logger, reader *request.Reader, jsonWriter *response.JSONWriter, repository *Repository) *Handler {
	return &Handler{log: log, reader: reader, jsonWriter: jsonWriter, repository: repository}
}

func (h *Handler) addUser() http.HandlerFunc {
	// TODO: Add better password validation(verify for special characters, number etc.
	type Request struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,gte=8"`
		Role     string `json:"role" validate:"required,oneof=admin editor author contributor"`
	}

	type Response struct {
		ID            int32     `json:"id"`
		Name          string    `json:"name"`
		Email         string    `json:"email"`
		Role          string    `json:"role"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &Request{}

		ok := h.reader.ReadJSONAndValidate(w, r, req)
		if !ok {
			return
		}

		// FIXME: Hash password before saving
		user, err := h.repository.AddUser(r.Context(), query.AddUserParams{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
			Role:     query.UserRole(req.Role),
		})
		if err != nil {
			// TODO: Handle user constraint errors separately
			h.log.Error("Failed adding user to DB", zap.Error(err))
			h.jsonWriter.DefaultError(w, r)
			return
		}

		resp := &Response{
			ID:            user.ID,
			Name:          user.Name,
			Email:         user.Email,
			Role:          string(user.Role),
			CreatedAt:     user.CreatedAt,
			UpdatedAt:     user.UpdatedAt,
		}

		// TODO: Send verification Email

		h.jsonWriter.Ok(w, r, resp)
	}
}

package users

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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

func (h *Handler) getUser() http.HandlerFunc {
	// TODO: Add better password validation(verify for special characters, number etc.

	type Response struct {
		ID            int32     `json:"id"`
		Name          string    `json:"name"`
		Email         string    `json:"email"`
		Role          string    `json:"role"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			h.log.Error("ID should be int type", zap.Error(err))
			h.jsonWriter.BadRequest(w, r, &UserNotFound{})
			return
		}
		

		user, err := h.repository.GetUser(r.Context(), int32(id))
		if err != nil {
			// TODO: Handle user constraint errors separately
			h.log.Error("Failed get the user from DB", zap.Error(err))
			h.jsonWriter.BadRequest(w, r, &UserNotFound{ID: int32(id)})
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


		h.jsonWriter.Ok(w, r, resp)
	}
}

func (h *Handler) getAllUser() http.HandlerFunc {
	type Response struct {
		ID            int32     `json:"id"`
		Name          string    `json:"name"`
		Email         string    `json:"email"`
		Role          string    `json:"role"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// FIXME: Hash password before saving
		users, err := h.repository.GetAllUsers(r.Context())
		if err != nil {
			// TODO: Handle user constraint errors separately
			h.log.Error("Failed get all the user from DB", zap.Error(err))
			h.jsonWriter.BadRequest(w, r, &UserNotFound{})
			return
		}

		resp := make([]Response, 0)

		for _, user := range users {
			r := Response{
				ID:            user.ID,
				Name:          user.Name,
				Email:         user.Email,
				Role:          string(user.Role),
				CreatedAt:     user.CreatedAt,
				UpdatedAt:     user.UpdatedAt,
			}

			resp = append(resp, r)
		}

		h.jsonWriter.Ok(w, r, resp)
	}
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


		h.jsonWriter.Ok(w, r, resp)
	}
}

func (h *Handler) updateUser() http.HandlerFunc{
	type Request struct {
		ID       int32     	`json:"id" validate:"required"`
		Name     string 		`json:"name" validate:"required"`
		Email    string 		`json:"email" validate:"required,email"`
		Role     string 		`json:"role" validate:"required,oneof=admin editor author contributor"`
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

		user, err := h.repository.UpdateUser(r.Context(), query.UpdateUserParams{
			ID: 			req.ID,
			Name:     req.Name,
			Email:    req.Email,
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


		h.jsonWriter.Ok(w, r, resp)
	}
}

func (h *Handler) updateUserPassword() http.HandlerFunc{
	type Request struct {
		ID       int32     	`json:"id" validate:"required"`
		Password string `json:"password" validate:"required,gte=8"`
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

		user, err := h.repository.UpdateUserPassword(r.Context(), query.UpdateUserPasswordParams{
			ID: 			req.ID,
			Password: req.Password,
		})
		if err != nil {
			// TODO: Handle user constraint errors separately
			h.log.Error("Failed updating user password to DB", zap.Error(err))
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


		h.jsonWriter.Ok(w, r, resp)
	}
}

func (h *Handler) DeleteUser() http.HandlerFunc {
	type Request struct {
		ID            int32     `json:"id" validate:"required"`
	}

	type Response struct {
		Message				string		`json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &Request{}

		ok := h.reader.ReadJSONAndValidate(w, r, req)
		if !ok {
			return
		}

		err := h.repository.DeleteUser(r.Context(), req.ID)
		if err != nil {
			h.log.Error("Failed to delete user from DB", zap.Error(err))
			h.jsonWriter.DefaultError(w, r)
			return
		}

		resp := &Response{
			Message: "User successfully deleted",
		}

		h.jsonWriter.Ok(w, r, resp)
	}
}
package users

import (
	"github.com/aveloper/blog/internal/utils"
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

		hash, valid ,err := utils.ValidateAndHashPassword(req.Password)
		if !valid {
			if err == nil {
				h.log.Error("Password is not meeting the requirement "+ req.Password)
				h.jsonWriter.InvalidPasswordErr(w, r)
			} else {
				//TODO: need to change the DefaultError method
				h.log.Error("Failed to hash the password ", zap.Error(err))
				h.jsonWriter.DefaultError(w, r)
			}

			return
		}

		// FIXME: Hash password before saving
		user, err := h.repository.AddUser(r.Context(), query.AddUserParams{
			Name:     req.Name,
			Email:    req.Email,
			Password: hash,
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

func (h *Handler) updateUserName() http.HandlerFunc{
	type Request struct {
		ID       int32     	`json:"id" validate:"required"`
		Name     string `json:"name" validate:"required"`
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

		user, err := h.repository.UpdateUserName(r.Context(), query.UpdateUserNameParams{
			ID: 			req.ID,
			Name:     req.Name,
		})
		if err != nil {
			h.log.Error("Failed updating user name to DB", zap.Error(err))
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

func (h *Handler) updateUserEmail() http.HandlerFunc{
	type Request struct {
		ID       int32     	`json:"id" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
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

		user, err := h.repository.UpdateUserEmail(r.Context(), query.UpdateUserEmailParams{
			ID:				req.ID,
			Email:    req.Email,
		})
		if err != nil {
			h.log.Error("Failed updating user email to DB", zap.Error(err))
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

		hash, valid ,err := utils.ValidateAndHashPassword(req.Password)
		if !valid {
			if err == nil {
				h.log.Error("Password is not meeting the requirement "+ req.Password)
				h.jsonWriter.InvalidPasswordErr(w, r)
			} else {
				//TODO: need to change the DefaultError method
				h.log.Error("Failed to hash the password ", zap.Error(err))
				h.jsonWriter.DefaultError(w, r)
			}

			return
		}

		user, err := h.repository.UpdateUserPassword(r.Context(), query.UpdateUserPasswordParams{
			ID: 			req.ID,
			Password: hash,
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

func (h *Handler) updateUserRole() http.HandlerFunc{
	type Request struct {
		ID       int32     	`json:"id" validate:"required"`
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

		user, err := h.repository.UpdateUserRole(r.Context(), query.UpdateUserRoleParams{
			ID:				req.ID,
			Role:     query.UserRole(req.Role),
		})
		if err != nil {
			h.log.Error("Failed updating user role to DB", zap.Error(err))
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
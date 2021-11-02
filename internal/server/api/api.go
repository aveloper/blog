package api

import (
	"github.com/aveloper/blog/internal/db"
	"github.com/aveloper/blog/internal/http/request"
	"github.com/aveloper/blog/internal/http/response"
	"github.com/aveloper/blog/internal/users"
	"github.com/aveloper/blog/internal/validator"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func Routes(r *mux.Router, log *zap.Logger, db *db.DB) {
	v := validator.New(log)
	jw := response.NewJSONWriter(log)
	reader := request.NewReader(log, jw, v)

	ur := users.NewRepository(db, log)
	uh := users.NewHandler(log, reader, jw, ur)

	users.UserRoutes(r, uh)

}

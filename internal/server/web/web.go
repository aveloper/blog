package web

import (
	"github.com/aveloper/blog/internal/admin"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func Routes(r *mux.Router, log *zap.Logger) {
	admin.Serve(r, log)
}

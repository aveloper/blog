package users

import (
	"context"

	"go.uber.org/zap"

	"github.com/aveloper/blog/internal/db"
	"github.com/aveloper/blog/internal/db/query"
)

//Repository has CRUD functions for users
type Repository struct {
	q   *query.Queries
	log *zap.Logger
}

//NewRepository creates a new instance of Repository
func NewRepository(db *db.DB, log *zap.Logger) *Repository {
	return &Repository{
		q:   query.New(db),
		log: log,
	}
}

//AddUser adds a new User
func (r *Repository) AddUser(ctx context.Context, user query.AddUserParams) (*query.User, error) {
	u, err := r.q.AddUser(ctx, user)
	if err != nil {
		r.log.Error("failed to add a new user", zap.Error(err))
		return nil, err
	}

	return &u, nil
}

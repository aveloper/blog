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

//GetUser fetch the user from users by given id
func (r *Repository) GetUser(ctx context.Context, id int32) (*query.User, error) {
	u, err := r.q.FetchUser(ctx, id)
	if err != nil {
		r.log.Error("failed to get the user", zap.Error(err))
		return nil, err
	}

	return &u, nil
}

//GetAllUsers fetch all the users from users
func (r *Repository) GetAllUsers(ctx context.Context) ([]query.User, error) {
	u, err := r.q.FetchAllUsers(ctx)
	if err != nil {
		r.log.Error("failed to get the user", zap.Error(err))
		return nil, err
	}

	return u, nil
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

//UpdateUser update the existing user
func (r *Repository) UpdateUser(ctx context.Context, user query.UpdateUserParams) (*query.User, error) {
	u, err := r.q.UpdateUser(ctx, user)
	if err != nil {
		r.log.Error("failed to update user", zap.Error(err))
		return nil, err
	}

	return &u, nil
}

//UpdateUserPassword update the existing user password
func (r *Repository) UpdateUserPassword(ctx context.Context, user query.UpdateUserPasswordParams) (*query.User, error) {
	u, err := r.q.UpdateUserPassword(ctx, user)
	if err != nil {
		r.log.Error("failed to update user password", zap.Error(err))
		return nil, err
	}

	return &u, nil
}

//DeleteUser delete the existing user
func (r *Repository) DeleteUser(ctx context.Context, id int32) error {
	err := r.q.DeleteUser(ctx, id)
	if err != nil {
		r.log.Error("failed to delete the user ", zap.Error(err))
		return  err
	}

	return nil
}


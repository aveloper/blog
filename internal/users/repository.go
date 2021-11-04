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

//UpdateUser update the existing user
func (r *Repository) UpdateUser(ctx context.Context, user query.UpdateUserParams) (*query.User, error) {
	u, err := r.q.UpdateUser(ctx, user)
	if err != nil {
		r.log.Error("failed to update user", zap.Error(err))
		return nil, err
	}

	return &u, nil
}

//UpdateUserName update the existing user username
func (r *Repository) UpdateUserName(ctx context.Context, user query.UpdateUserNameParams) (*query.User, error) {
	u, err := r.q.UpdateUserName(ctx, user)
	if err != nil {
		r.log.Error("failed to update user name", zap.Error(err))
		return nil, err
	}

	return &u, nil
}

//UpdateUserEmail update the existing user email
func (r *Repository) UpdateUserEmail(ctx context.Context, user query.UpdateUserEmailParams) (*query.User, error) {
	u, err := r.q.UpdateUserEmail(ctx, user)
	if err != nil {
		r.log.Error("failed to update user email", zap.Error(err))
		return nil, err
	}

	return &u, nil
}

//UpdateUserRole update the existing user
func (r *Repository) UpdateUserRole(ctx context.Context, user query.UpdateUserRoleParams) (*query.User, error) {
	u, err := r.q.UpdateUserRole(ctx, user)
	if err != nil {
		r.log.Error("failed to update user role", zap.Error(err))
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


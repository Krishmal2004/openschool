package setup

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/openschool-org/openschool/db/sqlc"
	"github.com/openschool-org/openschool/internal/authz"
)

type Repository struct{ queries *db.Queries }

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{queries: db.New(pool)} }

func (r *Repository) countUsersByRole(ctx context.Context, role string) (int64, error) {
	return r.queries.CountUsersByRole(ctx, role)
}

func (r *Repository) createAdmin(ctx context.Context, id uuid.UUID, email, fullName string) (AdminUser, error) {
	user, err := r.queries.CreateUser(ctx, db.CreateUserParams{ID: id, Email: email, FullName: fullName, Role: authz.RoleAdmin})
	return AdminUser{ID: user.ID, Email: user.Email}, err
}

func (r *Repository) deleteUser(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteUser(ctx, id)
}

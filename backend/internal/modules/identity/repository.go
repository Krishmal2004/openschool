package identity

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/openschool-org/openschool/db/sqlc"
)

type userRepository struct{ queries *db.Queries }

func newUserRepository(pool *pgxpool.Pool) *userRepository {
	return &userRepository{queries: db.New(pool)}
}

func (r *userRepository) ensureExists(ctx context.Context, command ensureUserCommand) (provisionedUser, error) {
	user, err := r.queries.EnsureUserExists(ctx, db.EnsureUserExistsParams{
		ID: command.ID, Email: command.Email, FullName: command.FullName, Role: command.Role,
	})
	if err != nil {
		return provisionedUser{}, err
	}
	return provisionedUser{MustChangePassword: user.MustChangePassword}, nil
}

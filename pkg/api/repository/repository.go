package repository

import (
	"context"

	"github.com/ismaeljpv/auth-api/pkg/api/domain"
)

type Repository interface {
	FindByID(context.Context, int64) (domain.User, error)
	FindByEmail(context.Context, string) (domain.User, error)
	Create(context.Context, domain.User) (domain.User, error)
	Update(context.Context, domain.User) (domain.User, error)
	Delete(context.Context, int64) (string, error)
}

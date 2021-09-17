package service

import (
	"context"

	"github.com/ismaeljpv/auth-api/pkg/api/domain"
)

type Service interface {
	FindByID(context.Context, int64) (domain.User, error)
	Create(context.Context, domain.User) (domain.User, error)
	Update(context.Context, domain.User, int64) (domain.User, error)
	Delete(context.Context, int64) (string, error)
	Login(context.Context, string, string) (string, error)
}

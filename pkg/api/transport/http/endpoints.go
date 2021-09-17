package http

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/ismaeljpv/auth-api/pkg/api/domain"
	"github.com/ismaeljpv/auth-api/pkg/api/service"
)

type Endpoints struct {
	FindByID endpoint.Endpoint
	Create   endpoint.Endpoint
	Update   endpoint.Endpoint
	Delete   endpoint.Endpoint
	Login    endpoint.Endpoint
}

func MakeEndpoints(s service.Service) Endpoints {
	return Endpoints{
		FindByID: makeFindByIDEndpoint(s),
		Create:   makeCreateEndpoint(s),
		Update:   makeUpdateEndpoint(s),
		Delete:   makeDeleteEndpoint(s),
		Login:    makeLoginEndpoint(s),
	}
}

func makeFindByIDEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		id := request.(int64)
		user, err := s.FindByID(ctx, id)
		return user, err
	}
}

func makeCreateEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		user := request.(domain.User)
		newUser, err := s.Create(ctx, user)
		return newUser, err
	}
}

func makeUpdateEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UserWithIDRequest)
		updatedUser, err := s.Update(ctx, req.User, req.ID)
		return updatedUser, err
	}
}

func makeDeleteEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		id := request.(int64)
		msg, err := s.Delete(ctx, id)
		return GenericMessageResponse{
			Code:    http.StatusOK,
			Message: msg,
			Status:  http.StatusText(http.StatusOK),
		}, err
	}
}

func makeLoginEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		body := request.(LoginRequest)
		token, err := s.Login(ctx, body.Email, body.Password)
		return token, err
	}
}

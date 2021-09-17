package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/ismaeljpv/auth-api/pkg/api/auth"
	"github.com/ismaeljpv/auth-api/pkg/api/domain"
	"github.com/ismaeljpv/auth-api/pkg/api/repository"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo   repository.Repository
	logger log.Logger
}

func NewService(repo repository.Repository, logger log.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) FindByID(ctx context.Context, id int64) (domain.User, error) {
	reqUser := ctx.Value("user").(jwt.MapClaims)
	reqId := reqUser["id"].(float64)
	if id != int64(reqId) {
		return domain.User{}, errors.New("User can only request for his data")
	}

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (s *service) Create(ctx context.Context, user domain.User) (domain.User, error) {

	_, err := s.repo.FindByEmail(ctx, user.Email)
	if err == nil {
		return domain.User{}, errors.New("Email is already taken")
	}

	if len(user.Password) < 6 {
		return domain.User{}, errors.New("Password min length is 6 characters")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, errors.New("Internal Server Error")
	}
	user.Password = string(hashedPassword)
	user, er := s.repo.Create(ctx, user)
	if er != nil {
		return domain.User{}, er
	}
	return user, nil
}

func (s *service) Update(ctx context.Context, user domain.User, id int64) (domain.User, error) {

	reqUser := ctx.Value("user").(jwt.MapClaims)
	reqId := reqUser["id"].(float64)
	if id != int64(reqId) {
		return domain.User{}, errors.New("User can only update his data")
	}

	if user.ID != id {
		return domain.User{}, errors.New("Bad Request")
	}

	updUser, er := s.repo.FindByID(ctx, user.ID)
	if er != nil {
		return domain.User{}, er
	}

	if strings.Compare(updUser.Firstname, user.Firstname) != 0 {
		updUser.Firstname = user.Firstname
	}
	if strings.Compare(updUser.Lastname, user.Lastname) != 0 {
		updUser.Lastname = user.Lastname
	}
	if strings.Compare(updUser.Email, user.Email) != 0 {
		updUser.Email = user.Email
	}

	user, err := s.repo.Update(ctx, updUser)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (s *service) Delete(ctx context.Context, id int64) (string, error) {
	reqUser := ctx.Value("user").(jwt.MapClaims)
	reqId := reqUser["id"].(float64)
	if id != int64(reqId) {
		return "", errors.New("User can only inactivate his own account")
	}

	msg, err := s.repo.Delete(ctx, id)
	if err != nil {
		return "", err
	}
	return msg, nil
}

func (s *service) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("Invalid Credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		level.Warn(s.logger).Log("msg", err.Error())
		return "", errors.New("Invalid Credentials")
	}

	token, err := auth.EncodeToken(ctx, user)
	if err != nil {
		level.Warn(s.logger).Log("msg", fmt.Sprintf("Error Signing Token => %s", err.Error()))
		return "", errors.New("Internal Server Error")
	}

	return token, nil
}

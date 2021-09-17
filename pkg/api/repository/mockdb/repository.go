package mockdb

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/ismaeljpv/auth-api/pkg/api/domain"
	repo "github.com/ismaeljpv/auth-api/pkg/api/repository"
)

var mockData = []domain.User{
	{
		ID:        1,
		Firstname: "Ismael",
		Lastname:  "Peña",
		Email:     "ip@gmail.com",
		Password:  "$2a$10$jR2qQUJx9b/CiC81KtPbXu5/1rwOA7AjQzTxDXNnzv4/IweV14mj6",
		Status:    "ACTIVE",
		CreatedOn: time.Now().Unix(),
	},
	{
		ID:        2,
		Firstname: "Chris",
		Lastname:  "Vargas",
		Email:     "cjvp@gmail.com",
		Password:  "$2a$10$jR2qQUJx9b/CiC81KtPbXu5/1rwOA7AjQzTxDXNnzv4/IweV14mj6",
		Status:    "ACTIVE",
		CreatedOn: time.Now().Unix(),
	},
}

type repository struct {
	db     []domain.User
	logger log.Logger
}

func NewRepository(log log.Logger) repo.Repository {
	return &repository{
		db:     mockData,
		logger: log,
	}
}

func (r *repository) FindByID(ctx context.Context, id int64) (domain.User, error) {

	for _, user := range r.db {
		if user.ID == id && user.Status == "ACTIVE" {
			return user, nil
		}
	}
	return domain.User{}, errors.New("User Not Found")
}

func (r *repository) FindByEmail(ctx context.Context, email string) (domain.User, error) {

	for _, user := range r.db {
		if strings.Compare(user.Email, email) == 0 && user.Status == "ACTIVE" {
			return user, nil
		}
	}
	return domain.User{}, errors.New("User Not Found")
}

func (r *repository) Create(ctx context.Context, user domain.User) (domain.User, error) {

	for _, prevUser := range r.db {
		if prevUser.Email == user.Email {
			return domain.User{}, errors.New("Email is already registered")
		}
	}

	maxID := r.db[len(r.db)-1].ID
	user.ID = maxID + 1
	r.db = append(r.db, user)
	return user, nil
}

func (r *repository) Update(ctx context.Context, user domain.User) (domain.User, error) {

	for i, prevUser := range r.db {
		if prevUser.ID == user.ID && user.Status == "ACTIVE" {
			r.db[i] = user
			return user, nil
		}
	}
	return domain.User{}, errors.New("User Not Found")
}

func (r *repository) Delete(ctx context.Context, id int64) (string, error) {

	for i, prevUser := range r.db {
		if prevUser.ID == id && prevUser.Status == "ACTIVE" {
			r.db[i].Status = "INACTIVE"
			return "User deleted!", nil
		}
	}
	return "", errors.New("User Not Found")
}

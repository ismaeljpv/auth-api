package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ismaeljpv/auth-api/pkg/api/domain"
	repo "github.com/ismaeljpv/auth-api/pkg/api/repository"
)

type repository struct {
	db     *sql.DB
	logger log.Logger
}

func initDBConnection(ctx context.Context, logger log.Logger, uri string) (*sql.DB, error) {

	db, err := sql.Open("mysql", uri)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db, nil
}

func NewRepository(ctx context.Context, logger log.Logger, uri string) (repo.Repository, error) {
	dabatase, err := initDBConnection(ctx, logger, uri)
	if err != nil {
		return &repository{}, err
	}

	return &repository{
		db:     dabatase,
		logger: logger,
	}, nil
}

func (r *repository) FindByID(ctx context.Context, id int64) (domain.User, error) {
	var user domain.User
	query := "SELECT id, firstname, lastname, email, createdOn, password, status FROM users WHERE id = ? AND status = 'ACTIVE' "

	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Email, &user.CreatedOn, &user.Password, &user.Status)
	if err != nil {
		level.Warn(r.logger).Log("msg", fmt.Sprintf("No user found by ID = %v , error  = %s", id, err.Error()))
		return domain.User{}, errors.New("User Not Found")
	}

	return user, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	query := "SELECT id, firstname, lastname, email, createdOn, password, status FROM users WHERE email = ? AND status = 'ACTIVE' "

	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Email, &user.CreatedOn, &user.Password, &user.Status)
	if err != nil {
		level.Warn(r.logger).Log("msg", fmt.Sprintf("No user found by email = %v , error  = %s", email, err.Error()))
		return domain.User{}, errors.New("User Not Found")
	}

	return user, nil
}

func (r *repository) Create(ctx context.Context, user domain.User) (domain.User, error) {

	query := "INSERT INTO users(firstname, lastname, email, password, status, createdOn) VALUES(?, ?, ?, ?, ?, ?)"

	user.CreatedOn = time.Now().Unix()
	result, err := r.db.Exec(query, user.Firstname, user.Lastname, user.Email, user.Password, user.Status, user.CreatedOn)
	if err != nil {
		level.Warn(r.logger).Log("msg", fmt.Sprintf("Couldn't create the new user in the DB, error  = %s", err.Error()))
		return domain.User{}, errors.New("There was an error processing your request")
	}
	newId, er := result.LastInsertId()
	if er != nil {
		level.Warn(r.logger).Log("msg", fmt.Sprintf("Couldn't create the new user in the DB, error  = %s", er.Error()))
		return domain.User{}, errors.New("There was an error processing your request")
	}
	user.ID = newId
	return user, nil
}

func (r *repository) Update(ctx context.Context, user domain.User) (domain.User, error) {
	query := "UPDATE users SET firstname = ?, lastname = ?, email = ? WHERE id = ? AND status = 'ACTIVE' "

	result, err := r.db.Exec(query, user.Firstname, user.Lastname, user.Email, user.ID)
	if err != nil {
		level.Warn(r.logger).Log("msg", fmt.Sprintf("Couldn't update the the user with ID = %v, error  = %s", user.ID, err.Error()))
		return domain.User{}, errors.New("There was an error processing your request")
	}
	rows, er := result.RowsAffected()
	if rows == 0 || er != nil {
		level.Warn(r.logger).Log("msg", fmt.Sprintf("Couldn't update the the user with ID = %v", user.ID))
		return domain.User{}, errors.New("User Not Found")
	}
	return user, nil
}

func (r *repository) Delete(ctx context.Context, id int64) (string, error) {
	query := "UPDATE users SET status = 'INACTIVE' WHERE id = ?"

	result, err := r.db.Exec(query, id)
	if err != nil {
		level.Warn(r.logger).Log("msg", fmt.Sprintf("Couldn't update the the user with ID = %v, error  = %s", id, err.Error()))
		return "", errors.New("There was an error processing your request")
	}
	rows, er := result.RowsAffected()
	if rows == 0 || er != nil {
		level.Warn(r.logger).Log("msg", fmt.Sprintf("Couldn't update the the user with ID = %v", id))
		return "", errors.New("User Not Found")
	}

	return "User deleted!", nil
}

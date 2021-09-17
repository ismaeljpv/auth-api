package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/ismaeljpv/auth-api/pkg/api/domain"
)

type (
	UserWithIDRequest struct {
		ID   int64       `json:"id"`
		User domain.User `json:"user"`
	}

	LoginRequest struct {
		Email    string `json:"email"  validated:"required"`
		Password string `json:"password"  validated:"required"`
	}

	GenericMessageResponse struct {
		Message string `json:"message"`
		Status  string `json:"status"`
		Code    int64  `json:"code"`
	}
)

var validate *validator.Validate

func ValidateBody(s interface{}) error {
	validate = validator.New()
	errs := validate.Struct(s)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			return errors.New(fmt.Sprintf("Error on field %v, data is %v", err.Field(), err.ActualTag()))
		}
	}
	return nil
}

func DecodeParamIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	param, ok := mux.Vars(r)["id"]
	if !ok {
		return domain.User{}, errors.New("Invalid Request")
	}

	id, er := strconv.ParseInt(param, 0, 64)
	if er != nil {
		return domain.User{}, errors.New("Invalid Request")
	}

	return id, nil
}

func DecodeUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var body domain.User

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return domain.User{}, err
	}

	err = ValidateBody(&body)
	if err != nil {
		return domain.User{}, err
	}

	return body, nil
}

func DecodeUserWithIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var body domain.User
	param, ok := mux.Vars(r)["id"]
	if !ok {
		return domain.User{}, errors.New("Invalid Request")
	}

	id, er := strconv.ParseInt(param, 0, 64)
	if er != nil {
		return domain.User{}, errors.New("Invalid Request")
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return domain.User{}, err
	}

	err = ValidateBody(&body)
	if err != nil {
		return domain.User{}, err
	}

	return UserWithIDRequest{ID: id, User: body}, nil
}

func DecodeLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var body LoginRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	err = ValidateBody(&body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func EncodeLoginResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	token, ok := response.(string)
	if !ok {
		return errors.New("Internal Server Error")
	}

	w.Header().Add("Authorization", token)
	w.WriteHeader(http.StatusOK)
	return nil
}

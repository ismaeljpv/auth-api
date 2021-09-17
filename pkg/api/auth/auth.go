package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ismaeljpv/auth-api/pkg/api/domain"
	"github.com/joho/godotenv"
)

var SecretKey = []byte(getTokenSecret())

func getTokenSecret() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	return os.Getenv("JWT_SECRET")
}

func Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] != nil {

			token, err := jwt.Parse(r.Header["Authorization"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf(("Invalid Token"))
				}

				if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
					return nil, fmt.Errorf(("Expired token"))
				}

				return SecretKey, nil
			})

			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
			}

			user := token.Claims.(jwt.MapClaims)
			if token.Valid {
				ctx := context.WithValue(r.Context(), "user", user)
				h.ServeHTTP(w, r.WithContext(ctx))
			}
		} else {
			http.Error(w, "Access Denied", http.StatusForbidden)
		}
	})
}

func EncodeToken(ctx context.Context, user domain.User) (string, error) {
	newToken := jwt.New(jwt.SigningMethodHS256)
	claims := newToken.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["id"] = user.ID
	claims["user"] = user.Email
	claims["iss"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	token, err := newToken.SignedString(SecretKey)
	if err != nil {
		return "", errors.New("Error Signing Token")
	}

	return token, nil
}

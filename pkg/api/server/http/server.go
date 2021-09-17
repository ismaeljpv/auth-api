package http

import (
	"context"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/ismaeljpv/auth-api/pkg/api/auth"
	transport "github.com/ismaeljpv/auth-api/pkg/api/transport/http"
)

func NewHTTPServer(ctx context.Context, endpoints transport.Endpoints) http.Handler {

	router := mux.NewRouter()
	router.Use(contentMiddleware)

	router.Methods("GET").Path("/api/user/{id}").Handler(auth.Authenticate(httptransport.NewServer(
		endpoints.FindByID,
		transport.DecodeParamIDRequest,
		transport.EncodeResponse,
	)))

	router.Methods("POST").Path("/api/user").Handler(httptransport.NewServer(
		endpoints.Create,
		transport.DecodeUserRequest,
		transport.EncodeResponse,
	))

	router.Methods("PUT").Path("/api/user/{id}").Handler(auth.Authenticate(httptransport.NewServer(
		endpoints.Update,
		transport.DecodeUserWithIDRequest,
		transport.EncodeResponse,
	)))

	router.Methods("DELETE").Path("/api/user/{id}").Handler(auth.Authenticate(httptransport.NewServer(
		endpoints.Delete,
		transport.DecodeParamIDRequest,
		transport.EncodeResponse,
	)))

	router.Methods("POST").Path("/api/login").Handler(httptransport.NewServer(
		endpoints.Login,
		transport.DecodeLoginRequest,
		transport.EncodeLoginResponse,
	))

	return router
}

func contentMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

package main

import (
	"context"
	"flag"
	"fmt"
	golog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/ismaeljpv/auth-api/pkg/api/config"
	"github.com/ismaeljpv/auth-api/pkg/api/repository/mysql"
	httpserver "github.com/ismaeljpv/auth-api/pkg/api/server/http"
	"github.com/ismaeljpv/auth-api/pkg/api/service"
	httpinit "github.com/ismaeljpv/auth-api/pkg/api/transport/http"
)

func main() {

	var httpAddr = flag.String("http", ":8080", "HTTP listen address")
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewSyncLogger(logger)
	logger = log.With(logger,
		"service", "api",
		"time:", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,
	)
	errs := make(chan error)

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	flag.Parse()
	ctx := context.Background()

	dbURI, envErr := config.GetDBConnection()
	if envErr != nil {
		golog.Fatal(envErr)
	}

	// Init with mock repo
	//repo := mockdb.NewRepository(logger)

	// Init with mysql DB connection
	repo, err := mysql.NewRepository(ctx, logger, dbURI)
	if err != nil {
		golog.Fatal(err)
	}

	serv := service.NewService(repo, logger)
	HTTPEndpoints := httpinit.MakeEndpoints(serv)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		HTTPserver := httpserver.NewHTTPServer(ctx, HTTPEndpoints)
		level.Info(logger).Log("msg", fmt.Sprintf("HTTP Server listening on port %v", *httpAddr))
		errs <- http.ListenAndServe(*httpAddr, HTTPserver)
	}()

	level.Error(logger).Log("exit", <-errs)
	close(errs)
}

package main

import (
	"context"
	"sync"

	"github.com/cldmnky/krcrdr/internal/api"
	"github.com/cldmnky/krcrdr/internal/api/handlers/record"
	ctrl "sigs.k8s.io/controller-runtime"
)

func runApiServer(wg *sync.WaitGroup) {
	defer wg.Done()
	// Setup auth
	fakeAuth := &record.FakeAuthenticator{}
	options := api.Options{
		Env:           "dev",
		Addr:          apiAddr,
		ApiLogger:     ctrl.Log.WithName("api"),
		Authenticator: fakeAuth,
	}

	server := api.NewServer(options)
	setupLog.Info("Starting API server", "addr", options.Addr)
	server.Run(context.Background())
}

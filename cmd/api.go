package main

import (
	"context"
	"sync"

	"github.com/cldmnky/krcrdr/internal/api"
	ctrl "sigs.k8s.io/controller-runtime"
)

func runApiServer(wg *sync.WaitGroup) {
	defer wg.Done()
	options := api.Options{
		Env:       "dev",
		Addr:      apiAddr,
		ApiLogger: ctrl.Log.WithName("api"),
	}

	server := api.NewServer(options)
	setupLog.Info("Starting API server", "addr", options.Addr)
	server.Run(context.Background())
}

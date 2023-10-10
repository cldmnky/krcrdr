/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/cldmnky/krcrdr/cmd"
	ctrl "sigs.k8s.io/controller-runtime"
)

var log = ctrl.Log.WithName("main")

func main() {
	// Set up a Ctrl-C context.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err := cmd.New().ExecuteContext(ctx); err != nil {
		log.Error(err, "error running command")
		os.Exit(1)
	}
}

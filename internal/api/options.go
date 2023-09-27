package api

import "github.com/go-logr/logr"

type Options struct {
	Env       string
	Addr      string
	ApiLogger logr.Logger
}

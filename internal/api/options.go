package api

import (
	"github.com/cldmnky/krcrdr/internal/api/handlers/record"
	"github.com/go-logr/logr"
)

type Options struct {
	Env           string
	Addr          string
	ApiLogger     logr.Logger
	Authenticator record.JWSValidator
}

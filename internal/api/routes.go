package api

import (
	"github.com/cldmnky/krcrdr/internal/api/handlers/base"
	"github.com/cldmnky/krcrdr/internal/api/handlers/record"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes register the routes for the server
func (s *Server) RegisterRoutes(r *gin.Engine) error {
	r.Static("/assets", "./assets")
	base.Mount(r, base.NewHandler(base.NewService()))
	if err := record.Mount(r, s.Options.Authenticator); err != nil {
		s.Options.ApiLogger.Error(err, "failed to mount record handler")
		return err
	}
	return nil
}

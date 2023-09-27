package api

import (
	"github.com/cldmnky/krcrdr/internal/api/handlers/base"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes register the routes for the server
func (s *Server) RegisterRoutes(r *gin.Engine) error {
	r.Static("/assets", "./assets")
	base.Mount(r, base.NewHandler(base.NewService()))
	return nil
}

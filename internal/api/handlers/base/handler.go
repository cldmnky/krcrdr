package base

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	// Handler is the home page handler
	Handler interface {
		Root() gin.HandlerFunc
	}

	handler struct {
		service Service
	}
)

func NewHandler(service Service) Handler {
	return &handler{
		service: service,
	}
}

func Mount(r *gin.Engine, h Handler) {
	root := r.Group("/")
	root.GET("/", h.Root())
}

func (h *handler) Root() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, "foo")
	}
}

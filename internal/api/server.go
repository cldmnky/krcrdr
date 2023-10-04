package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Options Options
}

func NewServer(opt Options) *Server {
	return &Server{
		Options: opt,
	}
}

func (s *Server) Run(ctx context.Context) error {
	// Create context that listens for the interrupt signal from the OS.
	serverCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	r := gin.New()
	if err := s.RegisterRoutes(r); err != nil {
		return err
	}
	srv := &http.Server{
		Addr:    s.Options.Addr,
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Options.ApiLogger.Error(err, "failed to listen and serve")
			os.Exit(1)
		}
	}()
	// Listen for the interrupt signal.
	<-serverCtx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	serverCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return srv.Shutdown(serverCtx)
}

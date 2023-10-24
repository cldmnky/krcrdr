package api

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"sigs.k8s.io/controller-runtime/pkg/certwatcher"
)

type Server struct {
	Options        Options
	defaultingOnce sync.Once
	started        bool
	mu             sync.Mutex
}

func NewServer(opt Options) *Server {
	return &Server{
		Options: opt,
	}
}

func (s *Server) setDefaults() {
	s.Options.setDefaults()
}

// Start starts the server and blocks until the context is cancelled.
func (s *Server) Start(ctx context.Context) error {
	s.Options.ApiLogger.Info("Starting api server")
	s.defaultingOnce.Do(s.setDefaults)

	// Setup gin router.
	if s.Options.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	if err := s.RegisterRoutes(r); err != nil {
		return err
	}
	// Setup TLS config.
	// cfg is a pointer to a tls.Config struct that specifies the NextProtos field to be "h2".
	cfg := &tls.Config{ //nolint:gosec
		NextProtos: []string{"h2"},
	}
	for _, op := range s.Options.TLSOpts {
		op(cfg)
	}
	if cfg.GetCertificate == nil {
		// If the GetCertificate field is nil, the server will use the default
		certPath := filepath.Join(s.Options.CertDir, s.Options.CertName)
		keyPath := filepath.Join(s.Options.CertDir, s.Options.KeyName)
		// Create the certificate watcher and set the config's GetCertificate on the TLSConfig.
		certWatcher, err := certwatcher.New(certPath, keyPath)
		if err != nil {
			return err
		}
		cfg.GetCertificate = certWatcher.GetCertificate
		go func() {
			if err := certWatcher.Start(ctx); err != nil {
				s.Options.ApiLogger.Error(err, "certificate watcher error")
			}
		}()
	}

	// Load CA to verify client certificate, if configured.
	if s.Options.ClientCAName != "" {
		certPool := x509.NewCertPool()
		clientCABytes, err := os.ReadFile(filepath.Join(s.Options.CertDir, s.Options.ClientCAName))
		if err != nil {
			return fmt.Errorf("failed to read client CA cert: %w", err)
		}
		ok := certPool.AppendCertsFromPEM(clientCABytes)
		if !ok {
			return fmt.Errorf("failed to append client CA cert to CA pool")
		}
		cfg.ClientCAs = certPool
		cfg.ClientAuth = tls.RequireAndVerifyClientCert
	}
	listener, err := tls.Listen("tcp", net.JoinHostPort(s.Options.Host, strconv.Itoa(s.Options.Port)), cfg)
	if err != nil {
		return err
	}
	srv := http.Server{
		Handler:           r,
		MaxHeaderBytes:    1 << 20,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 32 * time.Second,
	}
	s.Options.ApiLogger.Info("Api server listening", "host", s.Options.Host, "port", s.Options.Port)
	idleConnsClosed := make(chan struct{})
	go func() {
		<-ctx.Done()
		s.Options.ApiLogger.Info("shutting down api server with timeout of 1 minute")
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			s.Options.ApiLogger.Error(err, "failed to shutdown api server")
		}
		close(idleConnsClosed)
	}()

	s.mu.Lock()
	s.started = true
	s.mu.Unlock()
	if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
		return err
	}
	<-idleConnsClosed
	return nil
}

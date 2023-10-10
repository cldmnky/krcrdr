package api

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/cldmnky/krcrdr/cmd/options"
	"github.com/cldmnky/krcrdr/internal/api"
	"github.com/cldmnky/krcrdr/internal/api/handlers/record"
	"github.com/cldmnky/krcrdr/internal/api/store"
	"github.com/cldmnky/krcrdr/internal/api/store/providers/nats"
	"github.com/madflojo/testcerts"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/spf13/cobra"

	ctrl "sigs.k8s.io/controller-runtime"

	logf "sigs.k8s.io/controller-runtime/pkg/log"

	natsgo "github.com/nats-io/nats.go"
)

var apiLog = ctrl.Log.WithName("api")

func Complete(cmd *cobra.Command, args []string, ro *options.RootOptions, o *options.ApiOptions) error {
	// create a comma separated list of nats urls
	natsUrls := ""
	for _, url := range o.NatsUrl {
		natsUrls += url + ","
	}
	natsUrls = natsUrls[:len(natsUrls)-1]

	var defaultNats = natsUrls
	if o.RunNatsServer {
		// Start nats
		apiLog.Info("Starting NATS server")
		dir, err := os.MkdirTemp("", "store")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)
		natsOpts := &server.Options{
			JetStream: true,
			Debug:     true,
			Host:      "127.0.0.1",
			// mktmpdir
			StoreDir: dir,
		}
		ns, err := server.NewServer(natsOpts)
		if err != nil {
			return err
		}
		logf.Log.Info("Starting NATS server")
		ns.Start()
		defer ns.Shutdown()

		defaultNats = fmt.Sprintf("nats://%s:%d", natsOpts.Host, natsOpts.Port)
	}
	// setup nats options
	natsOpts := []natsgo.Option{}
	if o.NatsUserCredentials != "" {
		natsOpts = append(natsOpts, natsgo.UserCredentials(o.NatsUserCredentials))
	}
	// append nats options
	natsOpts = append(natsOpts, natsgo.RetryOnFailedConnect(true), natsgo.MaxReconnects(10), natsgo.ReconnectWait(time.Second))
	// Setup the store
	stream, err := nats.NewStream(
		defaultNats,
		natsOpts..., // append nats options
	)
	if err != nil {
		return err
	}

	kv, err := nats.NewKV(
		defaultNats,
		natsOpts..., // append nats options
	)
	if err != nil {
		return err
	}
	s := store.NewStore(stream, kv)

	if o.GenerateSelfSignedCert {
		// Generate a self-signed cert
		// Make sure the cert dir exists
		if err := os.MkdirAll(o.CertDir, 0755); err != nil {
			return err
		}
		defer os.RemoveAll(o.CertDir)
		cert, key, err := testcerts.GenerateCertsToTempFile(o.CertDir)
		if err != nil {
			return err
		}
		// get filename from path
		cert = filepath.Base(cert)
		key = filepath.Base(key)

		o.CertName = cert
		o.KeyName = key
	}

	// Start the API server
	fa, err := record.NewFakeAuthenticator()
	if err != nil {
		return err
	}
	host, port, err := net.SplitHostPort(o.Addr)
	if err != nil {
		return err
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return err
	}
	opts := &api.Options{
		Host:          host,
		Port:          portInt,
		Authenticator: fa,
		ApiLogger:     logf.Log.WithName("api"),
		Store:         s,
		CertDir:       o.CertDir,
		CertName:      o.CertName,
		KeyName:       o.KeyName,
		Tracer:        ro.Tracer,
		Debug:         o.Debug,
	}
	if err := api.NewServer(*opts).Start(ctrl.SetupSignalHandler()); err != nil {
		return err
	}
	return nil
}

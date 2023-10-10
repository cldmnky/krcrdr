package webhook

import (
	"flag"
	"os"

	recorderv1beta1 "github.com/cldmnky/krcrdr/api/v1beta1"
	apiclient "github.com/cldmnky/krcrdr/internal/api/handlers/record/client"
	"github.com/cldmnky/krcrdr/internal/recorder"
	"github.com/cldmnky/krcrdr/internal/webhook"

	"github.com/cldmnky/krcrdr/cmd/options"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/madflojo/testcerts"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	k8swebhook "sigs.k8s.io/controller-runtime/pkg/webhook"
	ctrladmission "sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	scheme     = runtime.NewScheme()
	webhookLog = ctrl.Log.WithName("webhook")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(recorderv1beta1.AddToScheme(scheme))
}

// Complete sets up and starts the webhook server with the given options.
// It generates certs if requested and registers a recorder webhook.
// It also sets up health and readiness checks.
// Returns an error if there was a problem starting the server.
func Complete(cmd *cobra.Command, args []string, ro *options.RootOptions, o *options.WebhookOptions) error {
	zapOpts := zap.Options{
		Development: o.Debug,
	}
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	zapOpts.BindFlags(flag.CommandLine)
	flag.Parse()
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&zapOpts)))
	// Generate certs if requested
	if o.GenerateCerts {
		webhookLog.Info("Generating certs to: " + o.CertDir)
		// Create CertDir if it doesn't exist
		if err := os.MkdirAll(o.CertDir, 0755); err != nil {
			webhookLog.Error(err, "unable to create cert dir")
			return err
		}
		cert, key, err := testcerts.GenerateCertsToTempFile(o.CertDir)
		if err != nil {
			webhookLog.Error(err, "unable to generate certs")
			return err
		}
		// get filename from path
		cert = cert[len(o.CertDir)+1:]
		key = key[len(o.CertDir)+1:]
		webhookLog.Info("Generated certs", "cert", cert, "key", key)
		o.CertName = cert
		o.KeyName = key
	}
	webhookLog.Info("Setting up webhook")
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		Metrics:                metricsserver.Options{BindAddress: o.MetricsAddr},
		HealthProbeBindAddress: o.ProbeAddr,
		WebhookServer: k8swebhook.NewServer(
			k8swebhook.Options{
				CertDir:  o.CertDir,
				CertName: o.CertName,
				KeyName:  o.KeyName,
			},
		),
	})
	if err != nil {
		webhookLog.Error(err, "unable to start manager")
		return err
	}
	wh := mgr.GetWebhookServer()
	dec := ctrladmission.NewDecoder(scheme)
	apiClient, err := apiclient.NewApiClient(o.ApiRemoteAddr, o.ApiToken, false)
	if err != nil {
		webhookLog.Error(err, "unable to create api client")
		return err
	}

	r := recorder.NewRecorder(apiClient, ro.Tracer)
	wh.Register("/recorder", &k8swebhook.Admission{
		Handler: &webhook.RecorderWebhook{
			Client:   mgr.GetClient(),
			Decoder:  dec,
			Recorder: r,
			Tracer:   ro.Tracer,
		},
	})

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		webhookLog.Error(err, "unable to set up health check")
		return err
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		webhookLog.Error(err, "unable to set up ready check")
		return err
	}
	webhookLog.Info("Starting webhook server")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		webhookLog.Error(err, "problem running webhook server")
		return err
	}
	return nil
}

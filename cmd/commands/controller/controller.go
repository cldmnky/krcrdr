package controller

import (
	"flag"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	recorderv1beta1 "github.com/cldmnky/krcrdr/api/v1beta1"
	"github.com/cldmnky/krcrdr/cmd/options"
	"github.com/cldmnky/krcrdr/internal/controller"
)

var (
	scheme        = runtime.NewScheme()
	controllerLog = ctrl.Log.WithName("controller")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(recorderv1beta1.AddToScheme(scheme))
}

func Complete(cmd *cobra.Command, args []string, o *options.ControllerOptions) error {
	zapOpts := zap.Options{
		Development: o.Debug,
	}
	// bind flags to zapOpts here
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	zapOpts.BindFlags(flag.CommandLine)
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&zapOpts)))
	controllerLog.Info("Setting up controller")
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		Metrics:                metricsserver.Options{BindAddress: o.MetricsAddr},
		HealthProbeBindAddress: o.ProbeAddr,
		LeaderElection:         o.EnableLeaderElection,
		LeaderElectionID:       "controller.krcrdr.blahoonga.me",
	})
	if err != nil {
		controllerLog.Error(err, "unable to start manager")
		return err
	}
	if err = (&controller.ConfigReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		controllerLog.Error(err, "unable to create controller", "controller", "Config")
		return err
	}
	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		controllerLog.Error(err, "unable to set up health check")
		return err
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		controllerLog.Error(err, "unable to set up ready check")
		return err
	}
	controllerLog.Info("Starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		controllerLog.Error(err, "problem running manager")
		return err
	}
	return nil
}

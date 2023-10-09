package options

import (
	"github.com/spf13/cobra"
)

type ControllerOptions struct {
	MetricsAddr          string
	EnableLeaderElection bool
	ProbeAddr            string
	Debug                bool
}

var _ Interface = (*ControllerOptions)(nil)

func (o *ControllerOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.MetricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	cmd.Flags().StringVar(&o.ProbeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	cmd.Flags().BoolVar(&o.EnableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	cmd.Flags().BoolVar(&o.Debug, "debug", false, "Enable debug logging")

}

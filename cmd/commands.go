package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/cldmnky/krcrdr/cmd/options"
	"github.com/cldmnky/krcrdr/internal/tracer"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	ro      = &options.RootOptions{}
	Version = "1.0.0"
)

func New() *cobra.Command {
	ro.Version = Version
	cmd := &cobra.Command{
		Use:               "krcrdr",
		Short:             "krcrdr is a Kubernetes controller that records events to a database",
		Long:              `krcrdr is a Kubernetes controller that records events to a database`,
		DisableAutoGenTag: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			traceExporter, err := tracer.NewExporter(ro.OTLPExporter, ro.OTLPAddr, cmd.OutOrStdout())
			cobra.CheckErr(err)
			traceProvider, err := tracer.NewProvider(cmd.Context(), "version", traceExporter)
			cobra.CheckErr(err)
			ro.Tracer = traceProvider.Tracer("krcrdr")
			go func() {
				err = tracer.StartTracer(cmd.Context(), traceExporter)
				cobra.CheckErr(err)
			}()
			return bindViper(cmd, args, "KRCRDR")
		},
		Run: func(cmd *cobra.Command, args []string) {
			// print help
			cmd.Help()
		},
	}
	ro.AddFlags(cmd)
	cmd.AddCommand(Controller())
	cmd.AddCommand(Webhook())
	cmd.AddCommand(Api())

	opts := zap.Options{
		Development: ro.Debug,
	}
	// Add global flags
	opts.BindFlags(flag.CommandLine)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	return cmd
}

func bindViper(cmd *cobra.Command, args []string, prefix string) error {
	v := viper.New()

	v.SetEnvPrefix(prefix)
	v.AutomaticEnv()
	bindFlags(cmd, v, prefix)

	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper, prefix string) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent keys with underscores
		if strings.Contains(f.Name, "-") {
			v.BindEnv(f.Name, flagToEnvVar(f.Name, prefix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

func flagToEnvVar(flag, prefix string) string {
	envVarSuffix := strings.ToUpper(strings.ReplaceAll(flag, "-", "_"))
	return fmt.Sprintf("%s_%s", prefix, envVarSuffix)
}

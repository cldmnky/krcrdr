package cmd

import (
	"fmt"
	"strings"

	"github.com/cldmnky/krcrdr/cmd/options"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
	}
	ro.AddFlags(cmd)
	cmd.AddCommand(Controller())
	cmd.AddCommand(Webhook())
	cmd.AddCommand(Api())
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

package cmd

import (
	"flag"

	"github.com/cldmnky/krcrdr/cmd/commands/api"
	"github.com/cldmnky/krcrdr/cmd/options"
	"github.com/spf13/cobra"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var cmdLog = ctrl.Log.WithName("cmd")

func Api() *cobra.Command {
	o := &options.ApiOptions{}
	cmd := &cobra.Command{
		Use:   "api",
		Short: "Run the api server for krcrdr",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return bindViper(cmd, args, "API")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := zap.Options{
				Development: o.Debug,
			}
			opts.BindFlags(flag.CommandLine)
			flag.Parse()

			ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
			return api.Complete(cmd, args, o)
		},
	}
	o.AddFlags(cmd)
	return cmd
}

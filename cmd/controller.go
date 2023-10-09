package cmd

import (
	"github.com/cldmnky/krcrdr/cmd/commands/controller"
	"github.com/cldmnky/krcrdr/cmd/options"
	"github.com/spf13/cobra"
)

// Controller returns a pointer to a cobra.Command that runs the controller for krcrdr.
// It takes no arguments and returns an error if the command fails to complete.
func Controller() *cobra.Command {
	o := &options.ControllerOptions{}
	cmd := &cobra.Command{
		Use:   "controller",
		Short: "Run the controller for krcrdr",
		RunE: func(cmd *cobra.Command, args []string) error {
			return controller.Complete(cmd, args, o)
		},
	}
	o.AddFlags(cmd)
	return cmd
}

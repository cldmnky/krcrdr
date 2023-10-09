package cmd

import (
	"github.com/cldmnky/krcrdr/cmd/commands/webhook"
	"github.com/cldmnky/krcrdr/cmd/options"
	"github.com/spf13/cobra"
)

func Webhook() *cobra.Command {
	o := &options.WebhookOptions{}
	cmd := &cobra.Command{
		Use:   "webhook",
		Short: "Run the webhook server for krcrdr",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return bindViper(cmd, args, "WEBHOOK")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return webhook.Complete(cmd, args, ro, o)
		},
	}
	o.AddFlags(cmd)
	return cmd
}

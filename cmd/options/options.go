package options

import "github.com/spf13/cobra"

// Interface defines the methods that an options struct should implement.
type Interface interface {
	// AddFlags adds this options' flags to the cobra command.
	AddFlags(cmd *cobra.Command)
}

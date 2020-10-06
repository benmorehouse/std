package help

import "github.com/spf13/cobra"

// Command will provide key editing capabilities
var Command = &cobra.Command{
	Use:   "help",
	Short: "Show the list of commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

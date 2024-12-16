package cli

import "github.com/spf13/cobra"

// rootCmd represents the base command when called without any subcommands
var periodicCmd = &cobra.Command{
	Use:   "pt",
	Short: "Periodic task subcommands ",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

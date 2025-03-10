package cmd

import "github.com/spf13/cobra"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	SilenceUsage: true,
	Short:        "Task Manager Demo",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}

func init() {
	serverCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to listen on")
	rootCmd.AddCommand(serverCmd)
}

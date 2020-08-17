package cmd

import (
	"github.com/renatomdg/rodai/pkg"

	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute the specified flow",
	Run: pkg.Executor,
}

func init() {
	rootCmd.AddCommand(execCmd)
}

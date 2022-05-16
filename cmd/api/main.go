package main

import (
	"fmt"
	"os"

	cmd "github.com/icaroribeiro/new-go-code-challenge-template-2/tools/api/commands"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "api",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func main() {
	rootCmd.AddCommand(cmd.VersionCmd)
	rootCmd.AddCommand(cmd.RunCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

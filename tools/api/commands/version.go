package commands

import (
	"fmt"

	"github.com/icaroribeiro/new-go-code-challenge-template-2/tools/api"
	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the API version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("api v%s\n", api.Version)
	},
}

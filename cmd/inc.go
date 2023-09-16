package cmd

import (
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var IncCmd = &cobra.Command{
	Use:   "inc",
	Short: "inc test code",
	Long:  `this is an INC test code`,
        Run: func(cmd *cobra.Command, args []string) {

        },


}

func init() {

	RootCmd.AddCommand(IncCmd)

}

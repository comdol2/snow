package cmd

import (
        "fmt"
        "log"
	"os"

        "github.com/spf13/cobra"
)

// branchCreateCmd represents the branch get command
var IncGetCmd = &cobra.Command{
        Use:   "get",
        Short: "inc get test code",
        Long:  `this is an INC get test code`,
        Run: func(cmd *cobra.Command, args []string) {

                if strings.HasPrefix(strings.ToUpper(pIncidentNumber), "INC") {

                        fmt.Println("\n===== INCIDENT ======\n")

                        resp, respstr, err := sClient.GetIncident(pIncidentNumber, "", "")
                        if err != nil {
                                log.Fatalf("ERROR: %v", err)
                        }
                        if resp == nil {
                                log.Fatalf("ERROR: No " + pIncidentNumber + " found!!")
                        } else {
                                fmt.Println(respstr)
                        }

                } else {

                        log.Fatalf("ERROR: Mandator parameter(s) missing. Check with --help")

                }

        },
}

func init() {

        IncCmd.AddCommand(IncGetCmd)

        IncGetCmd.PersistentFlags().StringVarP(&pIncidentNumber, "number", "n", "",   "INC1234567")

}

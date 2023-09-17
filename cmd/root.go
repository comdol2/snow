package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"

	snow "github.com/comdol2/snow/api"
	"github.com/spf13/cobra"
)

var sClient *snow.Client
var version, eSNOW_INSTANCE, eSNOW_USERNAME, eSNOW_PASSWORD string
var snowInstanceURL, pIncidentNumber string
var debug bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "snow",
	Version: version,
	Short:   "test code",
	Long:    `this is a test code`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "To turn-on debugging")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	strMyOS := "Linux"
	if runtime.GOOS == "windows" {
		strMyOS = "Windows"
	}
	if debug {
		fmt.Println("My OS : ", strMyOS)
	}

	eSNOW_INSTANCE = os.Getenv("SNOW_INSTANCE")
	eSNOW_USERNAME = os.Getenv("SNOW_USERNAME")
	eSNOW_PASSWORD = os.Getenv("SNOW_PASSWORD")
	snowInstanceURL = "https://" + eSNOW_INSTANCE

	sClient = snow.NewClient(eSNOW_USERNAME, eSNOW_PASSWORD, snowInstanceURL, debug)
	if sClient == nil {
		log.Fatalf("ERROR: Can't create SNOW client")
	}

}

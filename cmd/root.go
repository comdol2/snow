package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v3"

	"github.com/rhysd/abspath"
	"github.com/spf13/cobra"
	snow "github.com/comdol2/snow/api"
)

var version string

var strWhoAmI, strWhoAmI_sys_id, pDownload, pGithubInstance, eVaultInstance, pVaultInstance, pSnowInstance string
var eGITHUB_INSTANCE, eGITHUB_TOKEN, eGITHUB_TOKEN_str, pGITHUB_INSTANCE, eVAULT_INSTANCE, pVAULT_INSTANCE, pVaultTokenFile, pGithubTokenFile, SNOW_Table string
var eSNOW_INSTANCE, pSNOW_INSTANCE, eSNOW_USERNAME, pSNOW_USERNAME, eSNOW_PASSWORD, eSNOW_PASSWORD_str, pSNOW_PASSWORD string
var vaultInstanceURL, snowInstanceURL, githubInstanceURL string
var debug, dryrun, pVersion, pGetWithCtasks, pGetInRaw, pTimezone, pGetMyChanges, pGetMyCtasks, pGetMyIncidents, pOEPipeline bool
var sClient *snow.Client
var githubTokenFile, vaultTokenFile, vaultSecretPath string
var pYamlFilename, pBappIdCheck, pCICheck, pTaskCheck, pRelatedCheck, pBappIdTable, pTemplateName, pWorkNotes, pChangeRequestNumber, pChangeTaskNumber, pIncidentNumber, pRitmNumber, pYamlCheck, pCloseState, pCloseCodes, pCauseCodesArea, pCauseCodesSubArea, pCloseNotes, pPhaseState, pState string
var pBuildYAMLForOEPipeine, pBuildYAML, pTemplateNameWithDetails bool
var pBuildYAMLFromCHG, pBuildYAMLFromCTASK, pBuildYAMLFromINC, pBuildYAMLEnv, pBuildYAMLBappid, pBuildYAMLTemplate, pBuildYAMLOutput string
var pOEAction, pOEClosecode string
var pGetListOptions string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "snow",
	Version: version,
	Short:   "test code",
	Long: 	 `this is a test code`,
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

	eVAULT_INSTANCE = os.Getenv("VAULT_INSTANCE")
	eSNOW_INSTANCE = os.Getenv("SNOW_INSTANCE")

	eSNOW_USERNAME = os.Getenv("SNOW_USERNAME")
	eSNOW_PASSWORD = os.Getenv("SNOW_PASSWORD")

	snowInstanceURL = "https://" + eSNOW_INSTANCE

	sClient = snow.NewClient(eSNOW_USERNAME, eSNOW_PASSWORD, snowInstanceURL, debug)
	if sClient == nil {
		log.Fatalf("ERROR: Can't create SNOW client")
	}

}

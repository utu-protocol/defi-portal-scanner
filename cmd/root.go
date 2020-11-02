package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/utu-crowdsale/defi-portal-scanner/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	environment string
	debug       bool
	settings    config.Schema
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "defi-portal-scanner",
	Short: "A brief description of your application",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(v string) {
	rootCmd.Version = v
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is /etc/defi-portal-scanner/config.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Print debug logs")
	rootCmd.PersistentFlags().StringVar(&environment, "env", "develop", "Set the runtime environment name")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if debug {
		// log.SetReportCaller(true)
		// Output to stdout instead of the default stderr
		// Can be any io.Writer, see below for File example
		// log.SetOutput(os.Stdout)

		// Only log the warning severity or above.
		log.SetLevel(log.DebugLevel)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		// home, err := homedir.Dir()
		// if err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }
		// Search config in home directory with name ".thenewsroom" (without extension).
		viper.AddConfigPath("/etc/defi-portal-scanner")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match
	config.Defaults()

	// If a config file is found, read it in, else use the defaults
	if err := viper.ReadInConfig(); err == nil {
		viper.Unmarshal(&settings)
		if errors := config.Validate(&settings); len(errors) > 0 {
			log.Fatal(errors)
		}
		log.Println("Using config file at ", viper.ConfigFileUsed())
	} else {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			viper.Unmarshal(&settings)
		}
	}
	// make the version available via settings
	settings.RuntimeVersion = rootCmd.Version
	settings.RuntimeEnvironment = environment
	settings.RuntimeName = "defi-portal-scanner"
	log.Debug(fmt.Sprintf("config %#v", settings))
}

package cmd

import (
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/makasim/sentryhook"
	log "github.com/sirupsen/logrus"
	"github.com/utu-crowdsale/defi-portal-scanner/collector"
	"github.com/utu-crowdsale/defi-portal-scanner/utils"

	"github.com/spf13/cobra"
)

var (
	fromFile string
	toFile   string
)

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "A brief description of your command",
	Long:  ``,
	Run:   listen,
}

func init() {
	rootCmd.AddCommand(listenCmd)
}

func listen(cmd *cobra.Command, args []string) {

	log.SetFormatter(&utils.EmojiLogFormatter{})
	if debug {
		// Only log the warning severity or above.
		log.SetLevel(log.DebugLevel)
	}
	// enable sentry logging
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         settings.Services.GlitchtipDsn,
		Release:     fmt.Sprint(settings.RuntimeName, "@", settings.RuntimeEnvironment),
		Environment: settings.RuntimeEnvironment,
	})
	if err != nil {
		log.Warn("Sentry will be disabled - sentry.Init: ", err)
	}
	// add hook to sentry for logging events
	log.AddHook(sentryhook.New([]log.Level{
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel}))
	// start the service
	if err := collector.Start(settings); err != nil {
		log.Fatal(err)
	}
}

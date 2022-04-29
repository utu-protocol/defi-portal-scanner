package cmd

import (
	"fmt"
	"sync"

	"github.com/getsentry/sentry-go"
	"github.com/makasim/sentryhook"
	log "github.com/sirupsen/logrus"
	"github.com/utu-crowdsale/defi-portal-scanner/collector"
	"github.com/utu-crowdsale/defi-portal-scanner/wallet"

	"github.com/spf13/cobra"
)

var (
	dryRun              bool
	httpEnabled         bool
	scanEnabled         bool
	protocolsDescriptor string
	tokensDescriptor    string
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
	listenCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Enable dry-run for the utu api")
	listenCmd.Flags().BoolVar(&httpEnabled, "http", false, "Enable http API to submit addresses")
	listenCmd.Flags().BoolVar(&scanEnabled, "scan", false, "Enable defi protocols subscription scanning")
	listenCmd.Flags().StringVarP(&protocolsDescriptor, "protocols", "p", "", "Override the protocols file description location")
	listenCmd.Flags().StringVarP(&tokensDescriptor, "tokens", "t", "", "Override the tokens file description location")
}

func listen(cmd *cobra.Command, args []string) {
	//log.SetFormatter(&utils.EmojiLogFormatter{})
	if debug {
		// Only log the warning severity or above.
		log.SetLevel(log.DebugLevel)
	}
	log.Info("starting listener")
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

	// set the dryrun option
	settings.UTUTrustAPI.DryRun = settings.UTUTrustAPI.DryRun || dryRun
	if protocolsDescriptor != "" {
		settings.DefiSourcesFile = protocolsDescriptor
	}

	if tokensDescriptor != "" {
		settings.TokensDataFile = tokensDescriptor
	}

	collector.Ready(settings)
	collector.BalanceCollectorReady(settings)
	wallet.Ready(settings)
	// synchronize services
	var wg sync.WaitGroup
	// start the service
	if scanEnabled {
		log.Info("scanning mode enabled")
		wg.Add(1)
		go func() {
			if err := collector.Start(settings); err != nil {
				log.Fatal(err)
			}
		}()
	}

	if httpEnabled {
		log.Info("http mode enabled")
		wg.Add(1)
		go func() {
			if err := collector.Serve(settings); err != nil {
				log.Fatal(err)
			}
		}()
	}
	wg.Wait()
}

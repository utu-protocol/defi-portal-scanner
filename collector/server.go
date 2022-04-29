package collector

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/utu-crowdsale/defi-portal-scanner/config"
	"github.com/utu-crowdsale/defi-portal-scanner/wallet"
)

// Serve - serve the web interface
func Serve(cfg config.Schema) (err error) {
	// echo start
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.GET("/", func(c echo.Context) (err error) {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "ok",
			"version": config.Settings.RuntimeVersion,
		})
	})
	// health check :)
	e.GET("/status", func(c echo.Context) (err error) {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "ok",
			"version": config.Settings.RuntimeVersion,
		})
	})

	e.POST("/subscribe/:address", func(c echo.Context) (err error) {
		// _ := NewAddressFromString(c.Param("address"))
		// err = Scan(cfg, address)
		// if err != nil {
		// 	log.Error(err)
		// 	return c.JSON(http.StatusTeapot, map[string]string{})
		// }
		tokensParam := c.QueryParam("tokens")
		var tokens []string
		if len(tokensParam) > 0 {
			tokens = strings.Split(tokensParam, ",")
		}
		wallet.ScanTokensBalances(cfg, c.Param("address"), tokens)
		if err != nil {
			log.Error(err)
			return c.JSON(http.StatusTeapot, map[string]string{})
		}
		return c.JSON(http.StatusOK, map[string]string{})
	})
	err = e.Start(cfg.Server.ListenAddress)
	if err != nil {
		log.Errorf("error starting server for %s: %v", cfg.RuntimeName, err)
	}
	return
}

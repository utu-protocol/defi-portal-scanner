package collector

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/utu-crowdsale/defi-portal-scanner/config"
)

func Serve(cfg config.Schema) (err error) {
	// echo start
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	// health check :)
	e.GET("/status", func(c echo.Context) (err error) {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "ok",
			"version": config.Settings.RuntimeVersion,
		})
	})

	e.POST("/subscribe/:address", func(c echo.Context) (err error) {
		address := c.Param("address")
		err = Scan(cfg, address)
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

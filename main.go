package main

import (
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io"
	"log"
	"main/config"
	"main/game"
	mapValidator "main/game/map-validator"
	"net/http"
	"os"
	"time"
)

func main() {
	// load config
	err := godotenv.Load()
	if err != nil {
		log.Println("Warn: .env file not found! Using OS values or fallback values.")
	}
	err = os.Mkdir("storage", 0750)
	if err != nil {
		panic(err)
	}
	config.GetGameConfig()
	restoreOrCreateMap()
	startCronJob()

	// start http server
	httpConfig := config.GetHttpConfig()
	e := echo.New()
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: time.Duration(httpConfig.TimeoutSeconds) * time.Second, // timeout requests longer than 4 seconds
	}))
	e.Use(middleware.BodyLimit(httpConfig.MaxSize)) // limit request size
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{httpConfig.FrontendHost},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{echo.POST, echo.GET},
	}))
	e.GET("/start", httpStartGame)
	e.GET("/config", httpConfigGame)
	e.POST("/validate", httpValidateGame)
	e.Logger.Fatal(e.Start("0.0.0.0:8080"))
}

func restoreOrCreateMap() {
	existMapJson := game.ReadFile()
	if len(existMapJson) > 0 {
		game.MapJson = string(existMapJson)
	} else {
		//Generate and store new.
		game.NewMap()
	}
}

func startCronJob() {
	schedule := os.Getenv("CRON_MAP_GENERATE")
	if schedule != "" { //Only enable cron if schedule set.
		clock := gocron.NewScheduler(time.UTC)
		_, err := clock.Cron(schedule).Do(game.NewMap)
		if err != nil {
			panic(err)
		}
		clock.StartAsync()
	}
}

func httpStartGame(c echo.Context) error {
	c.Response().Header().Add("Content-Type", "application/json; charset=utf-8")
	return c.String(http.StatusOK, game.MapJson)
}

func httpConfigGame(c echo.Context) error {
	return c.JSON(http.StatusOK, config.Game)
}

func httpValidateGame(c echo.Context) error {
	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	defer c.Request().Body.Close()
	if len(bodyBytes) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	unpackedMap, err := game.HttpInputToMap(bodyBytes)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	validation, err := mapValidator.Validate(unpackedMap)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	return c.JSON(http.StatusOK, validation)
}

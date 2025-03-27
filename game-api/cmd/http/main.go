package main

import (
	"megastrongbow-api/pkg/adapter/driven/psql"
	"megastrongbow-api/pkg/adapter/driver/rest"
	"megastrongbow-api/pkg/application/game"
	"net/http"
	"os"

	valid "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func main() {
	validator := valid.New()

	db := psql.InitDBConn(os.Getenv("DB_STRING"))

	// driven
	datastore := psql.NewStore(db)

	// driver
	gameService := game.NewGameService(datastore, validator)

	gameAPI := rest.NewGameHandler(gameService)
	server := rest.NewServer()

	server.GET("/health", func(c echo.Context) error { return c.String(http.StatusOK, "OK") })

	server.POST("/game/new", gameAPI.NewGame)
	server.GET("/game/:gameId", gameAPI.GetGamePlay)
	server.GET("/game/:gameId/score", gameAPI.GetScore)
	server.POST("/game/:gameId/score", gameAPI.SubmitScore)

	// not fully implemented
	server.GET("/me", gameAPI.GetMe)                      // not fully implemented
	server.POST("/player/new", gameAPI.NewPlayer)         // not fully implemented
	server.POST("/game/:gameId/start", gameAPI.StartGame) // not fully implemented
	server.POST("/game/:gameId/join", gameAPI.JoinGame)   // not fully implemented

	log.Logger.Fatal().Err(server.Start(":8080")).Msg("Error")
}

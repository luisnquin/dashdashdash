package main

import (
	"github.com/labstack/echo/v4"
	"github.com/luisnquin/dashdashdash/internal/config"
	"github.com/luisnquin/dashdashdash/internal/core"
	"github.com/luisnquin/dashdashdash/internal/storage"
	"github.com/luisnquin/go-log"
)

func main() {
	e := echo.New()

	config := config.New()

	db, err := storage.ConnectToTursoDB(config)
	if err != nil {
		panic(err)
	}

	core.InitControllers(e, db)

	if err := e.Start("localhost:8700"); err != nil {
		log.Err(err).Send()
	}
}

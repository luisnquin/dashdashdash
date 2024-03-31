package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/luisnquin/dashdashdash/internal/core"
)

func main() {
	e := echo.New()

	db, err := sqlx.Open("", "")
	if err != nil {
		log.Fatal(err)
	}

	core.InitControllers(e, db)

	if err := e.Start("localhost:8700"); err != nil {
		log.Panic(err)
	}
}

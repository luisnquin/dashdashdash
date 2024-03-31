package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	echo_log "github.com/labstack/gommon/log"
	"github.com/luisnquin/dashdashdash/internal/config"
	"github.com/luisnquin/dashdashdash/internal/core"
	"github.com/luisnquin/dashdashdash/internal/storage"
	"github.com/luisnquin/go-log"
)

func main() {
	e := echo.New()
	e.Logger.SetOutput(os.Stderr)
	e.Logger.SetLevel(echo_log.DEBUG)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	config := config.New()

	db, err := storage.ConnectToTursoDB(config)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	closers, err := core.Init(ctx, e, db)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := e.Start(":8700"); err != nil {
			log.Err(err).Send()
		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	for i, closer := range closers {
		if err := closer.Close(); err != nil {
			log.Err(err).Msgf("unable to close closer %d", i)
		}
	}

	if err := e.Shutdown(ctx); err != nil {
		log.Err(err).Msg("unable to shutdown server :<")
	}
}

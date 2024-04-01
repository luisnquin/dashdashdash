package main

import (
	"context"
	"io"
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
	log.Init(os.Stderr)

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

	cache, err := storage.NewRedisClient(ctx, config.Cache.GetRedisTrustedURL())
	if err != nil {
		panic(err)
	}

	closers, err := core.Init(ctx, e, config, db, cache)
	if err != nil {
		panic(err)
	}

	closers = append([]io.Closer{cache, db}, closers...)

	defer func() {
		for i, closer := range closers {
			if err := closer.Close(); err != nil {
				log.Err(err).Msgf("unable to close closer %d", i)
			}
		}
	}()

	for _, route := range e.Routes() {
		log.Debug().Msgf("%s - %s", route.Method, route.Path)
	}

	go func() {
		if err := e.Start(":8700"); err != nil {
			log.Err(err).Send()
		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Err(err).Msg("unable to shutdown server :<")
	}
}

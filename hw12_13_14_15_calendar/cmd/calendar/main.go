package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/elak/golang_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/elak/golang_home_work/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/elak/golang_home_work/hw12_13_14_15_calendar/internal/server/http"
	internalstorage "github.com/elak/golang_home_work/hw12_13_14_15_calendar/internal/storage/common"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	ret := run() //nolint:ifshort

	if ret != 0 {
		os.Exit(ret)
	}
}

func run() int {
	config, err := NewConfig()
	if err != nil {
		return 1
	}

	err = logger.Start(config.Logger.Level, config.Logger.Path)
	if err != nil {
		os.Exit(1)
	}
	defer logger.Stop()

	storage := internalstorage.New(config.Storage.Type)

	calendar := app.New(logger.GetDefaultLogger(), storage)

	server := internalhttp.NewServer(calendar, config.Server.Host, config.Server.Port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGHUP)

		select {
		case <-ctx.Done():
			return
		case <-signals:
		}

		signal.Stop(signals)
		cancel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logger.Error("failed to stop http server: " + err.Error())
		}
	}()

	err = storage.Connect(ctx, config.Storage.URI)
	if err != nil {
		logger.Error("failed to connect storage: " + err.Error())
		return 1
	}

	defer storage.Close(ctx)

	logger.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logger.Error("failed to start http server: " + err.Error())
		return 1
	}

	return 0
}

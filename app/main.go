package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-pkgz/lgr"

	"github.com/parMaster/logserver/app/server"
	"github.com/parMaster/logserver/app/store"
	"github.com/parMaster/logserver/config"
	"github.com/umputun/go-flags"
)

var Options struct {
	Config string `long:"config" env:"CONFIG" default:"config/logserver.toml" description:"toml config file name"`
	Cmd    string `long:"cmd" env:"CMD" description:"command to run (server, migrate)"`
}

func main() {
	if _, err := flags.Parse(&Options); err != nil {
		os.Exit(1)
	}

	config, err := config.NewConfig(Options.Config)
	if err != nil {
		log.Fatalf("error loading config: %e", err)
	}

	logOpts := []lgr.Option{lgr.LevelBraces, lgr.StackTraceOnError}
	if config.LogLevel == "debug" {
		logOpts = append(logOpts, lgr.Debug)
	}
	lgr.SetupStdLogger(logOpts...)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if x := recover(); x != nil {
			log.Printf("[WARN] run time panic:\n%v", x)
			panic(x)
		}

		// catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Println("Shutdown signal received\n*********************************")
		cancel()
	}()

	switch Options.Cmd {
	case "migrate":
		store.Migrate(ctx)
	case "server":
	default:
		if err := server.NewLogServer(ctx, *config).Start(); err != nil {
			log.Fatalf("Can't start logserver %e", err)
		}
	}
}

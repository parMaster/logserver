package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-pkgz/lgr"

	"github.com/parMaster/logserver/app/api"
	"github.com/parMaster/logserver/app/config"
	"github.com/parMaster/logserver/app/store"
	"github.com/umputun/go-flags"
)

var Options struct {
	Config string `long:"config" env:"CONFIG" default:"config.yml" description:"YAML config file name"`
	Cmd    string `long:"cmd" env:"CMD" description:"command to run (server, migrate)"`
	Dbg    bool   `long:"dbg" env:"DBG" description:"debug mode, overrides config Serve.Dbg"`
}

func main() {
	if _, err := flags.Parse(&Options); err != nil {
		log.Fatalf("error parsing flags: %e", err)
	}

	config, err := config.NewConfig(Options.Config)
	if err != nil {
		log.Fatalf("error loading config: %e", err)
	}

	logOpts := []lgr.Option{lgr.LevelBraces, lgr.StackTraceOnError}
	if Options.Dbg || config.Server.Dbg {
		logOpts = append(logOpts, lgr.Debug, lgr.Format(lgr.ShortDebug))
	}
	lgr.SetupStdLogger(logOpts...)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Println("Shutdown signal received\n*********************************")
		cancel()
	}()

	defer func() {
		if x := recover(); x != nil {
			log.Printf("[WARN] run time panic: %+v", x)
		}
	}()

	switch Options.Cmd {
	case "migrate":
		store.Migrate(ctx)
	case "service":
		RunService(ctx, *config)
	default:

		go RunService(ctx, *config)

		if err := api.NewApiServer(ctx, *config).Start(); err != nil {
			log.Fatalf("Can't start logserver %e", err)
		}
	}
}

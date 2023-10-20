package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-pkgz/lgr"
	"github.com/parMaster/logserver/app/server"
	"github.com/parMaster/logserver/config"
)

var (
	configPath string
)

func main() {
	flag.StringVar(&configPath, "config", "config/logserver.toml", "path to config file")
	flag.Parse()
	config, err := config.NewConfig(configPath)
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

	if err := server.NewLogServer(ctx, *config).Start(); err != nil {
		log.Fatalf("Can't start logserver %e", err)
	}
}

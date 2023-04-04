package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/go-pkgz/lgr"
	"github.com/parMaster/logserver/config"
	"github.com/parMaster/logserver/internal/app/logserver"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "config/logserver.toml", "path to config file")
}

func main() {
	config := config.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatalf(err.Error())
		os.Exit(1)
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
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop
		log.Println("Shutdown signal received\n*********************************")
		cancel()
	}()

	if err := logserver.NewLogServer(ctx, *config).Start(); err != nil {
		log.Fatalf("Can't start logserver %e", err)
	}
}

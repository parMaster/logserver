package main

import (
	"flag"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/go-pkgz/lgr"
	"github.com/parMaster/logserver/internal/app/logserver"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/logserver.toml", "path to config file")
}

func main() {

	logOpts := []lgr.Option{lgr.Debug, lgr.CallerFile, lgr.CallerFunc, lgr.Msec, lgr.LevelBraces, lgr.StackTraceOnError}
	lgr.SetupStdLogger(logOpts...)
	lgr.Setup(logOpts...)

	config := logserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		lgr.Fatalf(err.Error())
		os.Exit(1)
	}

	if err := logserver.Start(config); err != nil {
		lgr.Fatalf("Can't start logserver %s", err.Error())
		os.Exit(1)
	}
}

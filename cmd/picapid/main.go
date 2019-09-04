package main

import (
	"context"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"

	"github.com/xaionaro-go/picapi/config"
	"github.com/xaionaro-go/picapi/httpserver"
	"github.com/xaionaro-go/picapi/imageprocessor"
)

func fatalIfError(err error) {
	if err == nil {
		return
	}
	log.Fatal(err)
}

func getConfig() *config.Config {
	var cfg config.Config

	fatalIfError(envconfig.Process("picapi", &cfg))
	fatalIfError(cfg.Validate())

	return &cfg
}

func main() {
	cfg := getConfig()

	// TODO: uses a logger with native support of multi-leveling ("fatal", "error", "info" etc)

	var accessLogger *log.Logger
	var handlerLogger *log.Logger
	if cfg.LoggingLevel == `debug` {
		accessLogger = log.New(os.Stderr, `[picapid] [access]`, 0)
		handlerLogger = log.New(os.Stderr, `[picapid] [handler]`, 0)
	}

	srv, err := httpserver.NewHTTPServer(
		imageprocessor.NewImageProcessor(),
		accessLogger,
		handlerLogger,
	)
	fatalIfError(err)

	fatalIfError(srv.Start(
		context.Background(),
		cfg.ListenAddress,
	))
	fatalIfError(srv.Wait())
}

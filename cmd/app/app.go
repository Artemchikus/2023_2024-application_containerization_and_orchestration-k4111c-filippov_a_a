package app

import (
	"context"
	"find-ship/cmd/app/api"
	"find-ship/config"
	"find-ship/storage"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"
)

var confFile string

func init() {
	flag.StringVar(&confFile, "conf", "", "path to config file")
}

func Execute() (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("service exited with error: %v", err)
		}
	}()

	flag.Parse()

	log.Printf("confFile: %v\n", confFile)

	conf, err := config.Get(confFile)
	if err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	storage, err := storage.MustPostgresStorage(ctx, conf.DB)
	if err != nil && err != context.Canceled {
		panic(err)
	}

	if err := api.Start(ctx, storage, conf.App); err != nil && err != context.Canceled {
		panic(err)
	}

	return err
}

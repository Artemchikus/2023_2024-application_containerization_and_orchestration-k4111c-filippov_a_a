package main

import (
	"context"
	vkclient "find-ship/cmd/app/vk_client"
	appclient "find-ship/cmd/cron_lobs/update_urls/app_client"
	"find-ship/cmd/cron_lobs/update_urls/config"
	"flag"
	"log"
	"os/signal"
	"syscall"
)

var confFile string

func init() {
	flag.StringVar(&confFile, "conf", "", "path to config file")
}

func main() {
	var err error

	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("service exited with error: %v", err)
		}
	}()

	flag.Parse()

	conf, err := config.Get(confFile)
	if err != nil {
		panic(err)
	}

	app := appclient.MustAppClient(conf)
	vk := vkclient.MustVKClient()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	channels, err := app.GetVKChannels(ctx)
	if err != nil {
		panic(err)
	}

	for _, channel := range channels {
		select {
		case <-ctx.Done():
			return
		default:
			siteUrl, err := vk.GetChannelSiteUrl(ctx, *channel.ChannelURL)
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}

			if siteUrl != *channel.SiteURL {
				channel.SiteURL = &siteUrl
				err = app.UpdateVKChannel(ctx, channel)
				if err != nil {
					log.Printf("error: %v", err)
					continue
				}
			}
		}
	}
}

package vkclient

import (
	"context"

	"github.com/gocolly/colly"
)

type (
	VK interface {
		GetChannelSiteUrl(ctx context.Context, channelUrl string) (url string, err error)
	}
	vk struct {
		collector *colly.Collector
	}
)

func MustVKClient() VK {
	collector := colly.NewCollector()
	collector.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"
	return &vk{
		collector: collector,
	}
}

func (v *vk) GetChannelSiteUrl(ctx context.Context, channelUrl string) (url string, err error) {
	v.collector.OnHTML("a[data-group-track-type]", func(e *colly.HTMLElement) {
		url = e.Text
	})

	if err := v.collector.Visit(channelUrl); err != nil && err != colly.ErrAlreadyVisited {
		return url, err
	}

	return url, nil
}

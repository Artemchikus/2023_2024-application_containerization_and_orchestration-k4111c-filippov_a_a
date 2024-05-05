package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	Storage interface {
		VKStorage
		Reconnect(ctx context.Context) error
		Disconnect()
		Stat() *pgxpool.Stat
	}
	VKStorage interface {
		GetVKChannel(ctx context.Context, channelId int64) (channel VKChannel, err error)
		GetVKChannelByName(ctx context.Context, channelName string) (channel VKChannel, err error)
		GetVKChannels(ctx context.Context) (channels []VKChannel, err error)
		GetVKChannelsByType(ctx context.Context, channelType string) (channels []VKChannel, err error)
		UpdateVKChannel(ctx context.Context, channel VKChannel) (err error)
		DeleteVKChannel(ctx context.Context, channelId int64) (err error)
		CreateVKChannel(ctx context.Context, channel VKChannel) (channelId int64, err error)
	}
)

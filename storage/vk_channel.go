package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

func (p PostgresStorage) GetVKChannel(ctx context.Context, channelId int64) (channel VKChannel, err error) {
	conn, err := p.db.Acquire(ctx)
	if err != nil {
		return channel, err
	}
	defer conn.Release()

	args := pgx.NamedArgs{
		"channel_id": channelId,
	}

	rows, err := conn.Query(ctx, "SELECT * FROM vk_channels WHERE id = @channel_id", args)
	if err != nil {
		return channel, err
	}

	channel, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[VKChannel])
	if err != nil {
		return channel, err
	}

	return channel, nil
}
func (p PostgresStorage) GetVKChannelByName(ctx context.Context, channelName string) (channel VKChannel, err error) {
	conn, err := p.db.Acquire(ctx)
	if err != nil {
		return channel, err
	}
	defer conn.Release()

	args := pgx.NamedArgs{
		"channel_name": channelName,
	}

	rows, err := conn.Query(ctx, "SELECT * FROM vk_channels WHERE channel_name = @channel_name", args)
	if err != nil {
		return channel, err
	}

	channel, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[VKChannel])
	if err != nil {
		return channel, err
	}

	return channel, nil
}
func (p PostgresStorage) GetVKChannels(ctx context.Context) (channels []VKChannel, err error) {
	conn, err := p.db.Acquire(ctx)
	if err != nil {
		return channels, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, "SELECT * FROM vk_channels")
	if err != nil {
		return channels, err
	}

	channels, err = pgx.CollectRows(rows, pgx.RowToStructByName[VKChannel])
	if err != nil {
		return channels, err
	}

	if len(channels) == 0 {
		return channels, fmt.Errorf("no channels found")
	}

	return channels, nil
}
func (p PostgresStorage) GetVKChannelsByType(ctx context.Context, channelType string) (channels []VKChannel, err error) {
	conn, err := p.db.Acquire(ctx)
	if err != nil {
		return channels, err
	}
	defer conn.Release()

	args := pgx.NamedArgs{
		"channel_type": channelType,
	}

	rows, err := conn.Query(ctx, "SELECT * FROM vk_channels WHERE channel_type = @channel_type", args)
	if err != nil {
		return channels, err
	}
	defer rows.Close()

	channels, err = pgx.CollectRows(rows, pgx.RowToStructByName[VKChannel])
	if err != nil {
		return channels, err
	}

	if len(channels) == 0 {
		return channels, fmt.Errorf("no channels found")
	}

	return channels, nil
}
func (p PostgresStorage) UpdateVKChannel(ctx context.Context, channel VKChannel) (err error) {
	conn, err := p.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	query := "UPDATE vk_channels SET "
	qParts := make([]string, 0, 4)
	args := pgx.NamedArgs{}

	if channel.ChannelName != nil {
		qParts = append(qParts, "channel_name = @channel_name")
		args["channel_name"] = *channel.ChannelName
	}
	if channel.ChannelURL != nil {
		qParts = append(qParts, "channel_url = @channel_url")
		args["channel_url"] = *channel.ChannelURL
	}
	if channel.ChannelType != nil {
		qParts = append(qParts, "channel_type = @channel_type")
		args["channel_type"] = *channel.ChannelType
	}
	if channel.SiteURL != nil {
		qParts = append(qParts, "site_url = @site_url")
		args["site_url"] = *channel.SiteURL
	}

	if len(qParts) == 0 {
		return fmt.Errorf("nothing to update")
	}

	query += strings.Join(qParts, ", ") + " WHERE id = @channel_id"
	args["channel_id"] = channel.Id

	fmt.Printf("query: %v\n", query)

	res, err := conn.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("no rows updated")
	}

	return nil
}
func (p PostgresStorage) DeleteVKChannel(ctx context.Context, channelId int64) (err error) {
	conn, err := p.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	args := pgx.NamedArgs{
		"channel_id": channelId,
	}

	res, err := conn.Exec(ctx, "DELETE FROM vk_channels WHERE id = @channel_id", args)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("no rows deleted")
	}

	return nil
}
func (p PostgresStorage) CreateVKChannel(ctx context.Context, channel VKChannel) (channelId int64, err error) {
	conn, err := p.db.Acquire(ctx)
	if err != nil {
		return channelId, err
	}
	defer conn.Release()

	args := pgx.NamedArgs{
		"channel_name": channel.ChannelName,
		"channel_url":  channel.ChannelURL,
		"channel_type": channel.ChannelType,
		"site_url":     channel.SiteURL,
	}

	err = conn.QueryRow(ctx, "INSERT INTO vk_channels (channel_name, channel_url, channel_type, site_url) VALUES (@channel_name, @channel_url, @channel_type, @site_url) RETURNING id", args).Scan(&channelId)
	if err != nil {
		return channelId, err
	}

	return channelId, nil
}

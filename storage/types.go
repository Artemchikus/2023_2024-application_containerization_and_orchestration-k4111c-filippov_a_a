package storage

import "time"

type (
	VKChannel struct {
		Id          int64     `db:"id" json:"id"`
		ChannelName *string   `db:"channel_name" json:"channel_name,omitempty"`
		ChannelURL  *string   `db:"channel_url" json:"channel_url,omitempty"`
		ChannelType *string   `db:"channel_type" json:"channel_type,omitempty"`
		SiteURL     *string   `db:"site_url" json:"site_url,omitempty"`
		CreatedAt   time.Time `db:"created_at" json:"created_at"`
		UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	}
)

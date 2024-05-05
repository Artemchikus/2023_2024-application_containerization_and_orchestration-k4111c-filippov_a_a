package api

import (
	"find-ship/cmd/app/metrix"
	"find-ship/storage"
)

type (
	ctxKey int8

	ErrorResponse struct {
		Error string `json:"error"`
	}

	IDResponse struct {
		ID int64 `json:"id"`
	}

	CreateVKChannelRequest struct {
		ChannelName string `json:"channel_name"`
		ChannelURL  string `json:"channel_url"`
		ChannelType string `json:"channel_type"`
	}

	PatchVKChannelRequest struct {
		ChannelName *string `json:"channel_name,omitempty"`
		ChannelURL  *string `json:"channel_url,omitempty"`
		ChannelType *string `json:"channel_type,omitempty"`
		SiteURL     *string `json:"site_url,omitempty"`
	}

	MetricsResponse struct {
		metrix.AllMetrixData
	}

	ChannelResponse struct {
		storage.VKChannel
	}

	ChannelsResponse struct {
		Data []storage.VKChannel `json:"data"`
	}
)

const (
	CtxKeyRequestID ctxKey = iota
)

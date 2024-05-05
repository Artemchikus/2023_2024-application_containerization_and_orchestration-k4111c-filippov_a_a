package appclient

import (
	"bytes"
	"context"
	"encoding/json"
	"find-ship/cmd/app/api"
	"find-ship/cmd/cron_lobs/update_urls/config"
	"find-ship/storage"
	"fmt"
	"net/http"
	"net/url"
)

type (
	App interface {
		GetVKChannels(ctx context.Context) (channels []storage.VKChannel, err error)
		UpdateVKChannel(ctx context.Context, channel storage.VKChannel) (err error)
	}
	app struct {
		client     *http.Client
		appUrl     string
		apiVersion string
		appPort    string
	}
)

func MustAppClient(conf *config.Config) App {
	return &app{
		client:     http.DefaultClient,
		appUrl:     conf.ServiceUrl,
		apiVersion: conf.SefviceApi,
		appPort:    conf.ServicePort,
	}
}

func (app *app) GetVKChannels(ctx context.Context) (channels []storage.VKChannel, err error) {
	url, err := url.JoinPath(app.appUrl, app.appPort, app.apiVersion, "ships")
	if err != nil {
		return channels, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return channels, err
	}

	resp, err := app.client.Do(req)
	if err != nil {
		return channels, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(resp.Body).Decode(&channels); err != nil {
			return channels, err
		}

		return channels, nil
	default:
		result := api.ErrorResponse{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return channels, fmt.Errorf("app http status: %d. error: %v", resp.StatusCode, err)
		}

		return channels, fmt.Errorf("app http status: %d. error: %v", resp.StatusCode, result.Error)
	}
}

func (app *app) UpdateVKChannel(ctx context.Context, channel storage.VKChannel) (err error) {
	url, err := url.JoinPath(app.appUrl, app.apiVersion, "ships")
	if err != nil {
		return err
	}

	body, err := json.Marshal(channel)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	resp, err := app.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	default:
		result := api.ErrorResponse{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return fmt.Errorf("app http status: %d. error: %v", resp.StatusCode, err)
		}

		return fmt.Errorf("app http status: %d. error: %v", resp.StatusCode, result.Error)
	}
}

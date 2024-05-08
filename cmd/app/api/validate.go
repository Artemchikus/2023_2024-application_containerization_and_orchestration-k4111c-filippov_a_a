package api

import (
	"fmt"
	"regexp"
)

func (s *server) validate(request any) error {
	switch object := request.(type) {
	case CreateVKChannelRequest:
		if object.ChannelName == "" {
			return fmt.Errorf("empty channel name")
		}
		if object.ChannelURL == "" {
			return fmt.Errorf("empty channel url")
		}
		if !isVKURL(object.ChannelURL) {
			return fmt.Errorf("invalid channel url")
		}
		if object.ChannelType == "" {
			return fmt.Errorf("empty channel type")
		}
		return nil
	case PatchVKChannelRequest:
		if object.ChannelName == nil && object.ChannelURL == nil && object.ChannelType == nil && object.SiteURL == nil {
			return fmt.Errorf("nothing to update")
		}
		if object.ChannelURL != nil && !isVKURL(*object.ChannelURL) {
			return fmt.Errorf("invalid channel url")
		}
		if object.SiteURL != nil && !isURL(*object.SiteURL) {
			return fmt.Errorf("invalid site url")
		}
		return nil
	default:
		return fmt.Errorf("unknown request type")
	}
}

func isVKURL(url string) bool {
	regex := regexp.MustCompile(`^https://vk\.com/[a-zA-Z\-_0-9]+$`)
	return regex.MatchString(url)
}

func isURL(url string) bool {
	regex := regexp.MustCompile(`^(http|https)://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/[a-zA-Z0-9-._~:?#@!$&'()*+,;=]*)*$`)
	return regex.MatchString(url)
}

package client

import "time"

type Option func(client *Client)

func WithBatchingInterval(interval time.Duration) Option {
	return func(client *Client) {
		client.interval = interval
	}
}

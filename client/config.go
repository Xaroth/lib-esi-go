package client

import "net/http"

type Option func(*ESIClient)

func WithApplication(name, version string, contact ...string) Option {
	return func(c *ESIClient) {
		c.applicationName = name
		c.applicationVersion = version
		c.applicationContact = contact
	}
}

func WithClient(client *http.Client) Option {
	return func(c *ESIClient) {
		c.Client = client
	}
}

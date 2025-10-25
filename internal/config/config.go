package config

import (
	"net/url"

	defaults "github.com/xaroth/lib-esi-go"
)

type Defaults struct {
	host *url.URL

	userAgent         string
	compatibilityDate string
}

func (d *Defaults) Host() *url.URL {
	return d.host
}

func (d *Defaults) UserAgent() string {
	return d.userAgent
}

func (d *Defaults) CompatibilityDate() string {
	return d.compatibilityDate
}

type Option func(*Defaults)

func NewDefaults(opts ...Option) *Defaults {
	host, err := url.Parse(defaults.Host)
	if err != nil {
		panic(err)
	}
	defaults := &Defaults{
		host:              host,
		userAgent:         defaults.UserAgent,
		compatibilityDate: defaults.CompatibilityDate,
	}
	for _, opt := range opts {
		opt(defaults)
	}
	return defaults
}

func WithHost(url *url.URL) Option {
	return func(d *Defaults) {
		d.host = url
	}
}

func WithUserAgent(userAgent string) Option {
	return func(d *Defaults) {
		d.userAgent = userAgent
	}
}

func WithCompatibilityDate(compatibilityDate string) Option {
	return func(d *Defaults) {
		d.compatibilityDate = compatibilityDate
	}
}

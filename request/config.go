package request

import (
	"context"
	"fmt"
	"net/http"
	"time"

	defaults "github.com/xaroth/lib-esi-go"
	"github.com/xaroth/lib-esi-go/internal/config"
)

type Token interface {
	Owner() int64
	Token() string
}

type RefreshableToken interface {
	Token
	RefreshIfNeeded(ctx context.Context) error
}

type DefaultsProvider interface {
	RequestSender

	DefaultRequestConfig() *config.Defaults
}

type Config struct {
	*config.Defaults

	language          string
	scheduleTimeout   time.Duration
	requestTimeout    time.Duration
	token             Token
	additionalHeaders http.Header
}

func NewConfig(client RequestSender, opts ...Option) *Config {
	internalDefaults := config.NewDefaults()
	if client != nil {
		if esiClient, ok := client.(DefaultsProvider); ok {
			internalDefaults = esiClient.DefaultRequestConfig()
		}
	}

	cfg := &Config{
		Defaults: internalDefaults,
		language: defaults.Language,

		scheduleTimeout: defaults.ScheduleTimeout,
		requestTimeout:  defaults.RequestTimeout,

		token: nil,

		additionalHeaders: make(http.Header, 0),
	}
	cfg.apply(opts...)
	return cfg
}

type Option func(*Config)

func (c *Config) ScheduleTimeout() time.Duration {
	return c.scheduleTimeout
}

func (c *Config) RequestTimeout() time.Duration {
	return c.requestTimeout
}

func (c *Config) Token() Token {
	return c.token
}

func (c *Config) Language() string {
	return c.language
}

func (c *Config) apply(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
}

func (c *Config) ApplyRequest(req *http.Request) {
	c.ApplyHeaders(req.Header)
}

func (c *Config) ApplyHeaders(headers http.Header) {
	for key, values := range c.additionalHeaders {
		headers.Del(key)
		for _, value := range values {
			headers.Add(key, value)
		}
	}
	headers.Set("Accept-Language", c.Language())
	headers.Set("User-Agent", c.UserAgent())
	headers.Set("X-Compatibility-Date", c.CompatibilityDate())

	if c.token != nil {
		headers.Set("Authorization", fmt.Sprintf("Bearer %s", c.token.Token()))
	}
}

func WithToken(token Token) Option {
	return func(c *Config) {
		c.token = token
	}
}

func WithLanguage(language string) Option {
	return func(c *Config) {
		c.language = language
	}
}

func WithHeader(key, value string) Option {
	return func(c *Config) {
		c.additionalHeaders.Add(key, value)
	}
}

func WithETag(etag string) Option {
	return func(c *Config) {
		c.additionalHeaders.Set("If-None-Match", etag)
	}
}

func WithLastModified(lastModified time.Time) Option {
	return func(c *Config) {
		c.additionalHeaders.Set("If-Modified-Since", lastModified.Format(http.TimeFormat))
	}
}

func WithScheduleTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.scheduleTimeout = timeout
	}
}

func WithRequestTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.requestTimeout = timeout
	}
}

func WithDefaultOption(opt config.Option) Option {
	return func(c *Config) {
		opt(c.Defaults)
	}
}

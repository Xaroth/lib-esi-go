package client

import (
	"fmt"
	"net/http"
	"strings"

	defaults "github.com/xaroth/lib-esi-go"
	"github.com/xaroth/lib-esi-go/client/ratelimiter"
	"github.com/xaroth/lib-esi-go/internal/config"
	"github.com/xaroth/lib-esi-go/request"
)

type ESIClient struct {
	*http.Client
	*ratelimiter.Ratelimiter

	applicationName    string
	applicationVersion string
	applicationContact []string
}

var _ request.DefaultsProvider = (*ESIClient)(nil)
var _ request.SchedulingSender = (*ESIClient)(nil)

func New(opts ...Option) *ESIClient {
	client := &ESIClient{
		Client:      http.DefaultClient,
		Ratelimiter: ratelimiter.New(),
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func buildUserAgent(applicationName, applicationVersion string, applicationContact []string) string {
	sb := strings.Builder{}

	if applicationName != "" {
		sb.WriteString(applicationName)
		if applicationVersion != "" {
			sb.WriteString(fmt.Sprintf("/%s", applicationVersion))
		}
		if len(applicationContact) > 0 {
			sb.WriteString(fmt.Sprintf(" (%s)", strings.Join(applicationContact, "; ")))
		}
	} else {
		sb.WriteString("Unconfigured")
	}

	sb.WriteRune(' ')
	sb.WriteString(defaults.UserAgent)
	return sb.String()
}

func (c *ESIClient) DefaultRequestConfig() *config.Defaults {
	return config.NewDefaults(
		config.WithUserAgent(buildUserAgent(c.applicationName, c.applicationVersion, c.applicationContact)),
	)
}

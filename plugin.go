package caddy_pirsch_plugin

import (
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	pirsch "github.com/pirsch-analytics/pirsch-go-sdk/v2/pkg"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(PirschPlugin{})
}

type PirschPlugin struct {
	ClientId     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	BaseURL      string `json:"base_url,omitempty"`

	logger *zap.Logger
	client *pirsch.Client
}

func (m PirschPlugin) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.pirsch",
		New: func() caddy.Module { return new(PirschPlugin) },
	}
}

func (m *PirschPlugin) Provision(ctx caddy.Context) (err error) {
	m.client = pirsch.NewClient(m.ClientId, m.ClientSecret, &pirsch.ClientConfig{
		BaseURL: strings.TrimSpace(m.BaseURL),
	})
	m.logger = ctx.Logger(m)
	return err
}

func (m *PirschPlugin) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	request := r.Clone(r.Context())
	go func(r *http.Request) {
		clientIP, ok := caddyhttp.GetVar(r.Context(), caddyhttp.ClientIPVarKey).(string)

		if !ok || clientIP == "" {
			clientIP = r.RemoteAddr
		}

		if err := m.client.PageView(r, &pirsch.PageViewOptions{
			IP: clientIP,
		}); err != nil {
			m.logger.Error("failed sending page view to pirsch: %v", zap.Error(err))
		}
	}(request)
	return next.ServeHTTP(w, r)
}

var _ caddyhttp.MiddlewareHandler = (*PirschPlugin)(nil)

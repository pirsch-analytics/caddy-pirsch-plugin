package caddy_pirsch_plugin

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	httpcaddyfile.RegisterHandlerDirective("pirsch", parseCaddyfile)
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	p := new(PirschPlugin)

	for h.Next() {
		// configuration should be in a block
		for h.NextBlock(0) {
			switch h.Val() {
			case "client_id":
				var clientId string

				if !h.AllArgs(&clientId) {
					return nil, h.ArgErr()
				}

				p.ClientId = clientId
			case "client_secret":
				var clientSecret string

				if !h.AllArgs(&clientSecret) {
					return nil, h.ArgErr()
				}

				p.ClientSecret = clientSecret
			case "base_url":
				var baseUrl string

				if !h.AllArgs(&baseUrl) {
					return nil, h.ArgErr()
				}

				urlRegex := regexp.MustCompile(`https?://(www\.)?[-a-zA-Z0-9@:%._+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)

				if !urlRegex.MatchString(baseUrl) {
					return nil, h.Errf("'%s' is not a valid url", baseUrl)
				}

				p.BaseURL = baseUrl
			default:
				return nil, h.Errf("unrecognized option '%s'", h.Val())
			}
		}
	}

	if p.ClientId == "" && p.ClientSecret != "" && !strings.HasPrefix(p.ClientSecret, "pa_") {
		return nil, fmt.Errorf("invalid access key, it must start with 'pa_'")
	} else if p.ClientId == "" && p.ClientSecret == "" {
		return nil, fmt.Errorf("the client ID or secret is missing (oAuth)")
	}

	return p, nil
}

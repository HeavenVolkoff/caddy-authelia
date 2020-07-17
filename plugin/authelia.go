/* Copyright (c) 2020 Vítor Vasconcellos. All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package plugin

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/HeavenVolkoff/caddy-authelia/plugin/headers"
	"github.com/HeavenVolkoff/caddy-authelia/plugin/internalized/traefik"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"

	"go.uber.org/zap"
)

const (
	autheliaVerifyPath = "/api/verify"
)

func init() {
	caddy.RegisterModule(Authelia{})
	httpcaddyfile.RegisterHandlerDirective("authelia", parseCaddyfile)
}

// Authelia implements a plugin for securing routes with authentication
type Authelia struct {
	// If true, the connection to the authelia backend will use TLS
	TLS bool `json:"tls,omitempty"`
	// Port which the authelia backend is exposed
	Port uint16 `json:"port,omitempty"`
	// Host where the authelia backend can be reached
	Domain string `json:"domain,omitempty"`
	// URL to redirect unauthorized requests (Optional)
	RedirectURL string `json:"redirect_url,omitempty"`

	client http.Client
	logger *zap.Logger
}

var ( // Interface guards
	_ caddy.Validator             = (*Authelia)(nil)
	_ caddy.Provisioner           = (*Authelia)(nil)
	_ caddyfile.Unmarshaler       = (*Authelia)(nil)
	_ caddyhttp.MiddlewareHandler = (*Authelia)(nil)
)

func (Authelia) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.authelia",
		New: func() caddy.Module { return new(Authelia) },
	}
}

func (a Authelia) Validate() error {
	err := validateDomain(a.Domain)
	if err != nil {
		return err
	}

	if a.RedirectURL != "" {
		_, err = url.Parse(a.RedirectURL)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *Authelia) Provision(ctx caddy.Context) error {
	a.logger = ctx.Logger(a)
	a.logger.Info("Provisioning Authelia plugin instance")

	if a.Port == 0 {
		// Authelia default port
		a.Port = 9091
	}

	a.client = http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 30 * time.Second,
	}

	return nil
}

// Caddyfile Syntax:
//
//	authelia [<matcher>] <domain>:<port> [<port>] {
//		tls
//		port 		 <uint16>
//		domain 		 <string>
//  	redirect_url <string>
//	}
func (a *Authelia) UnmarshalCaddyfile(d *caddyfile.Dispenser) (err error) {
	for d.Next() {
		args := d.RemainingArgs()
		switch len(args) {
		case 1:
			a.Domain, a.Port, err = splitDomainPort(args[0])
			if err != nil {
				return err
			}
		case 2:
			a.Domain = args[0]
			if strings.Contains(a.Domain, ":") {
				return d.Err("port should only be defined once")
			}

			a.Port, err = parsePortNum(args[1])
			if err != nil {
				return err
			}
		default:
			return d.ArgErr()
		}

		var portArg = ""
		for d.NextBlock(0) {
			switch d.Val() {
			case "tls":
				if a.TLS {
					return d.Err("tls already specified")
				}
				a.TLS = true
			case "port":
				if portArg != "" {
					return d.Err("port already specified")
				}

				if !d.AllArgs(&portArg) {
					return d.ArgErr()
				}

				a.Port, err = parsePortNum(portArg)
			case "domain":
				if a.Domain != "" {
					return d.Err("domain already specified")
				}
				if !d.AllArgs(&a.Domain) {
					return d.ArgErr()
				}
			case "redirect_url":
				if a.RedirectURL != "" {
					return d.Err("redirect_url already specified")
				}
				if !d.AllArgs(&a.RedirectURL) {
					return d.ArgErr()
				}
			}
		}
	}

	return nil
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var a Authelia
	err := a.UnmarshalCaddyfile(h.Dispenser)
	return a, err
}

func (a Authelia) ServeHTTP(writer http.ResponseWriter, request *http.Request, nextHandler caddyhttp.Handler) (err error) {
	scheme := "http"
	if a.TLS {
		scheme = "https"
	}

	autheliaUrl := url.URL{
		Path:   autheliaVerifyPath,
		Host:   fmt.Sprintf("%s:%d", a.Domain, a.Port),
		Scheme: scheme,
	}

	if a.RedirectURL != "" {
		autheliaUrl.RawQuery = url.Values{"rd": []string{a.RedirectURL}}.Encode()
	}

	forwardRequest, err := http.NewRequest(http.MethodGet, autheliaUrl.String(), nil)
	if err != nil {
		return err
	}

	traefik.AssignForwardHeaders(request, forwardRequest)

	forwardResponse, err := a.client.Do(forwardRequest)
	if err != nil {
		return err
	}

	defer func() {
		// Deal with possible error during body closure
		innerErr := forwardResponse.Body.Close()
		if innerErr != nil {
			if err == nil {
				err = innerErr
			} else {
				err = fmt.Errorf("%v: %w", innerErr, err)
			}
		}
	}()

	body, err := ioutil.ReadAll(forwardResponse.Body)
	if err != nil {
		return err
	}

	if forwardResponse.StatusCode < http.StatusOK || forwardResponse.StatusCode >= http.StatusMultipleChoices {
		responseHeaders := http.Header{}
		headers.CopyHeadersWithoutHop(responseHeaders, forwardResponse.Header)

		return caddyhttp.StaticResponse{
			Body:       string(body),
			Close:      true,
			Headers:    responseHeaders,
			StatusCode: caddyhttp.WeakString(strconv.Itoa(forwardResponse.StatusCode)),
		}.ServeHTTP(writer, request, nextHandler)
	}

	remoteUser := forwardResponse.Header.Get(headers.RemoteUserHeader)
	remoteGroups := forwardResponse.Header.Get(headers.RemoteGroupsHeader)
	if remoteUser == "" || remoteGroups == "" {
		return caddyhttp.Error(
			http.StatusInternalServerError,
			fmt.Errorf("authelia failed to return a valid user"),
		)
	}

	// Setup authentication success the same way as CaddyAuth's API
	// https://github.com/caddyserver/caddy/blob/829e36d535cf5bbff7cf0f510608e6fca956cec4/modules/caddyhttp/caddyauth/caddyauth.go#L81-L85
	repl := request.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer)
	repl.Set("http.auth.user.id", remoteUser)
	repl.Set("http.auth.user.groups", remoteGroups)

	return nextHandler.ServeHTTP(writer, request)
}

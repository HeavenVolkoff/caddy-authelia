// This source code is a modified portion of traefik's fowardAuth Middleware
// https://github.com/containous/traefik/blob/73ca7ad0c195cb752147a36fba7c715aaa09e2e8/pkg/middlewares/auth/go#L161-L217
// Copyright (c) 2016-2020 Containous SAS
// Licensed under MIT, see ./LICENSE for more information

package traefik

import (
	"net"
	"net/http"
	"strings"

	"github.com/HeavenVolkoff/caddy-authelia/plugin/headers"
)

func AssignForwardHeaders(req, forwardReq *http.Request) {
	headers.CopyHeadersWithoutHop(forwardReq.Header, req.Header)

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		if prior, ok := req.Header[headers.XForwardedFor]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		forwardReq.Header.Set(headers.XForwardedFor, clientIP)
	}

	xMethod := req.Header.Get(headers.XForwardedMethod)
	switch {
	case xMethod != "":
		forwardReq.Header.Set(headers.XForwardedMethod, xMethod)
	case req.Method != "":
		forwardReq.Header.Set(headers.XForwardedMethod, req.Method)
	default:
		forwardReq.Header.Del(headers.XForwardedMethod)
	}

	xfp := req.Header.Get(headers.XForwardedProto)
	switch {
	case xfp != "":
		forwardReq.Header.Set(headers.XForwardedProto, xfp)
	case req.TLS != nil:
		forwardReq.Header.Set(headers.XForwardedProto, "https")
	default:
		forwardReq.Header.Set(headers.XForwardedProto, "http")
	}

	if xfp := req.Header.Get(headers.XForwardedPort); xfp != "" {
		forwardReq.Header.Set(headers.XForwardedPort, xfp)
	}

	xfh := req.Header.Get(headers.XForwardedHost)
	switch {
	case xfh != "":
		forwardReq.Header.Set(headers.XForwardedHost, xfh)
	case req.Host != "":
		forwardReq.Header.Set(headers.XForwardedHost, req.Host)
	default:
		forwardReq.Header.Del(headers.XForwardedHost)
	}

	xfURI := req.Header.Get(headers.XForwardedURI)
	switch {
	case xfURI != "":
		forwardReq.Header.Set(headers.XForwardedURI, xfURI)
	case req.URL.RequestURI() != "":
		forwardReq.Header.Set(headers.XForwardedURI, req.URL.RequestURI())
	default:
		forwardReq.Header.Del(headers.XForwardedURI)
	}
}

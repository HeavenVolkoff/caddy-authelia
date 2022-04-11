/* Copyright (c) 2020 Vítor Vasconcellos. All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package headers

import (
	"net/http"

	"github.com/HeavenVolkoff/caddy-authelia/plugin/internalized/golang"
	"github.com/HeavenVolkoff/caddy-authelia/plugin/internalized/oxy"
)

const (
	XForwardedURI      = "X-Forwarded-Uri"
	XForwardedFor      = "X-Forwarded-For"
	XForwardedHost     = "X-Forwarded-Host"
	XForwardedPort     = "X-Forwarded-Port"
	XForwardedProto    = "X-Forwarded-Proto"
	RemoteUserHeader   = "Remote-User"
	XForwardedMethod   = "X-Forwarded-Method"
	RemoteGroupsHeader = "Remote-Groups"
	RemoteEmailHeader  = "Remote-Email"
	RemoteNameHeader   = "Remote-Name"
)

func CopyHeadersWithoutHop(dst http.Header, src http.Header) {
	oxy.CopyHeaders(dst, src)
	oxy.RemoveHeaders(dst, golang.HopHeaders...)
}

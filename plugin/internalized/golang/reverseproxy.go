// This source code is a modified portion of Go's net/http/httputil/reverseproxy
// https://github.com/golang/go/blob/21898524f66c075d7cfb64a38f17684140e57675/src/net/http/httputil/reverseproxy.go#L169-L184
// Copyright (c) 2009 The Go Authors. All rights reserved.
// Licensed under BSD-3-Clause with a patent grant, see ./LICENSE and ./PATENTS for more information.

package golang

// Hop-by-hop headers. These are removed when sent to the backend.
// As of RFC 7230, hop-by-hop headers are required to appear in the
// Connection header field. These are the headers defined by the
// obsoleted RFC 2616 (section 13.5.1) and are used for backward
// compatibility.
var HopHeaders = []string{
	"Connection",
	"Proxy-Connection", // non-standard but still sent by libcurl and rejected by e.g. google
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",      // canonicalized version of "TE"
	"Trailer", // not Trailers per URL above; https://www.rfc-editor.org/errata_search.php?eid=4522
	"Transfer-Encoding",
	"Upgrade",
}

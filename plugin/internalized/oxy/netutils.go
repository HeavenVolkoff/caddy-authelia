// This source code is a modified portion of vulcand/oxy's utils/netutils
// https://github.com/vulcand/oxy/blob/f0cbb9d6b797d92d168b95b5c443a31dfa67ccd0/utils/netutils.go#L168-L174
// https://github.com/vulcand/oxy/blob/f0cbb9d6b797d92d168b95b5c443a31dfa67ccd0/utils/netutils.go#L186-L191
// Licensed under Apache-2.0, see ./LICENSE for more information.

package oxy

import "net/http"

// CopyHeaders copies http headers from source to destination, it
// does not overide, but adds multiple headers
func CopyHeaders(dst http.Header, src http.Header) {
	for k, vv := range src {
		dst[k] = append(dst[k], vv...)
	}
}

// RemoveHeaders removes the header with the given names from the headers map
func RemoveHeaders(headers http.Header, names ...string) {
	for _, h := range names {
		headers.Del(h)
	}
}

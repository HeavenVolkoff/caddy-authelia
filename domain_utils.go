/* Copyright (c) 2020 Vítor Vasconcellos. All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package authelia

import (
	"fmt"
	"github.com/HeavenVolkoff/caddy-authelia/external/Go"
	"strconv"
	"strings"
)

type DomainError struct {
	Domain string
}

var ( // Interface guards
	_ error = (*DomainError)(nil)
)

func (d DomainError) Error() string {
	return fmt.Sprintf("%s: Invalid domain", d.Domain)
}

func parsePortNum(s string) (uint16, error) {
	port, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return 0, err
	}

	return uint16(port), nil
}

func validateDomain(domain string) error {
	if Go.IsDomainName(domain) {
		return nil
	}
	return DomainError{Domain: domain}
}

func splitDomainPort(s string) (string, uint16, error) {
	if len(s) == 0 {
		return s, 0, nil
	}

	slice := strings.SplitN(s, ":", 2)

	if len(slice) == 1 {
		return slice[0], 0, nil
	}

	port, err := parsePortNum(slice[1])
	return slice[0], port, err
}

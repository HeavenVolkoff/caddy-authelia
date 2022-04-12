# caddy-authelia

[![project_badge](https://img.shields.io/badge/HeavenVolkoff/caddy2--authelia-black.svg?style=for-the-badge&logo=github "Project Badge")](https://github.com/HeavenVolkoff/caddy-authelia)
[![version_badge](https://img.shields.io/github/tag/HeavenVolkoff/caddy-authelia.svg?label=version&style=for-the-badge "Version Badge")](https://github.com/HeavenVolkoff/caddy-authelia/releases/latest)
[![license_badge](https://img.shields.io/github/license/HeavenVolkoff/caddy-authelia.svg?style=for-the-badge "License Badge")](https://www.mozilla.org/en-US/MPL/2.0/)

Caddy 2 plugin for integration with Authelia

> This plugin is still a work in progress.
> Use it in production at your own risk

## Example

The following is an example of using the plugin inside a Caddyfile:

```caddyfile
whoami.example.com {
    route {
        # Authelia URL : internal URL preferred
        authelia http://authelia:9091 {
            # redirect_url: optional
            # Public auth URL, used when not authenticated on whoami.example.com
            redirect_url authelia.example.com
        }
        request_header Remote-User {http.auth.user.id}
        request_header Remote-Groups {http.auth.user.groups}
        request_header Remote-Email {http.auth.user.email}
        request_header Remote-Name {http.auth.user.name}
        reverse_proxy http://whoami
    }
}

authelia.example.com {
    # Needed for 2FA to work - authelia checks X-Forwarded-Host, and Caddy doesn't fill it by default
    request_header X-Forwarded-Host {http.request.host}
    reverse_proxy http://authelia:9091
}
```

## License

This project is available under the Mozilla Public License 2.0 (MPL),
excepted where otherwise explicitly noted.

## Copyright

Copyright (c) 2020 Vítor Vasconcellos. All rights reserved.

I am not affiliated with Caddy or Authelia.

Caddy® is a registered trademark of Light Code Labs, LLC.

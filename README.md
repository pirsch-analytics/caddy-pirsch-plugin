# Caddy Plugin for Pirsch Analytics

A Caddy v2 plugin to track requests in [Pirsch Analytics](https://pirsch.io).

## Usage

```
pirsch [<matcher>] {
    client_id <pirsch client ID>
    client_secret <pirsch client secret or access key>
    base_url <alternative-api-url>
}
```

You can obtain these parameters from the [Settings](https://dashboard.pirsch.io/settings) section of your Pirsch dashboard. We recommend using an access key, as it doesn't require a roundtrip to get an oAuth token first. If you use the access key, the client ID must be left empty.

Because this directive does not come standard with Caddy, you need to [put the directive in order](https://caddyserver.com/docs/caddyfile/options). The correct place is up to you, but usually putting it near the end works if no other terminal directives match the same requests. It's common to pair a Pirsch handler with a `file_server`, so ordering it just before is often a good choice:

```
{
	order pirsch before file_server
}
```

Alternatively, you may use `route` to order it the way you want. For example:

```
localhost
root * /srv
route {
	pirsch * {
		[...]
	}
	file_server
}
```

### Example

Track all requests to HTML pages in Pirsch. You might want to extend the matcher regexp to also include `/` or, alternatively, match everything but assets (like `.css`, `.js`, ...) since usually you wouldn't want to track those.

`client ID` is optional when using an access key (recommended).

```
{
    order pirsch before file_server
}

http://localhost:8080 {
    @html path_regexp .*\.html$

    pirsch @html {
        client_id <client ID>
        client_secret <client secret or access key>
    }

    file_server
}
```

## Building Caddy

To compile Caddy with the Pirsch plugin, install and use [xcaddy](https://github.com/caddyserver/xcaddy):

```
xcaddy build --with github.com/pirsch-analytics/caddy-pirsch-plugin@plugin-tag
```

Replace `plugin-tag` with the plugin version, e.g. `v1.0.2`.

## License

Apache 2.0

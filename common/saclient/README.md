# sacloud-sdk-go/common/saclient

Golang binding of Sakura Cloud API client.

## Features

- Automatic authorization / authentication handling
- Automatic retry handling
- Client side rate limitter
- `httptest.Server` integration for better testing experience

## Quick start

The library provides `Client` struct as its sole public API.  All operations are against it.  You have to first allocate it:

```golang
import "github.com/sacloud/sacloud-sdk-go/common/saclient"

var theClient saclient.Client
```

You can pass various flags to the allocated client using environment variables, command line flags, and others.
Let's say we are going to accept those two here:

```golang
import "os"

func main() {
	fs := theClient.FlagSet()

	// Load settings from the command line arguments
	err := fs.Parse(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	// Load settings from the environment variables
	err = theClient.SetEnviron(os.Environ())
	if err != nil {
		os.Exit(1)
	}

	// This is optional (done automatically)
	// but it is a bit polite to do so explicitly
	err = theClient.Populate()
	if err != nil {
		os.Exit(1)
	}
}
```

After these set up you can use the client like other HTTP clients, by using its `Do` method.

```golang
import (
	"net/http"
	"strings"
)

func yourLogic() error {
	// Say you want to post something...
	api := "https://secure.sakura.ad.jp/cloud/zone/is1a/api/cloud/1.1/..."
	io := strings.NewReader(`{
	  "parameter": [
	    "values", "values", "..."
	  ]
	}`)

	// Create a request
	req, err := http.NewRequest("POST", api, io)
	if err != nil {
		return err
	}

	// then call it
	res, err := theClient.Do(req)
	if err != nil {
		return err
	}

	// ...
}
```

## Trusting a private CA

When an API endpoint is served with a certificate signed by a private CA — for
example a development mock server that generates its own self-signed
certificate, or a corporate CA in front of a staging endpoint — tell the client
to trust that CA. The configured CAs are added to (not replacing) the system
roots, so the client keeps talking to normally-trusted endpoints too.

```golang
// from PEM bytes in memory
client.SetWith(saclient.WithCACertPEM(pemBytes))

// or from a file
client.SetWith(saclient.WithCACertPath("/path/to/ca.pem"))
```

The same settings are also available as environment variables, following the
`SAKURA_PRIVATE_KEY` / `SAKURA_PRIVATE_KEY_PATH` convention:

| environment variable | value |
|----------------------|-------|
| `SAKURA_CA_CERT` | the PEM content itself |
| `SAKURA_CA_CERT_PATH` | the path of a PEM file |

When both are given, the direct PEM value (`WithCACertPEM` / `SAKURA_CA_CERT`)
takes precedence over the path (`WithCACertPath` / `SAKURA_CA_CERT_PATH`). The
client reads and validates the PEM when its settings are populated, so a missing
or broken file fails early.


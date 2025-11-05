# sacloud/saclient-go

Golang binding of Sakura Cloud API client.
An HTTP request doer implementation that communicates with Sakura cloud.

## Features

- Automatic authorization / authentication handling
- Automatic retry handling
- Client side rate limitter
- `httptest.Server` integration for better testing experience

## Quick start

The library provides `Client` struct as its sole public API.  All operations are against it.  You have to first allocate it:

```golang
import saht "github.com/sacloud/saclient-go"

var theClient saht.Client
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
	// but it is a bit polite to do be explicit
	err = theClient.Popuate()
	if err != nil {
		os.Exit(1)
	}
}
```

After these set up you can use the client like other HTTP clients, by using its `Do` method.

```golang
import (
	"strings"
	"http"
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
	req, err := http.NewRequest("POST", api, io);
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
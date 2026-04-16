// Copyright 2025- The sacloud/saclient-go Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/sacloud/saclient-go"
)

var theClient saclient.Client

func main() {
	fs := theClient.FlagSet(flag.PanicOnError)

	for _, arg := range os.Args {
		switch arg {
		case "--help", "-h":
			fmt.Printf("Usage: %s [options] <url>\n\n", os.Args[0])
			fs.PrintDefaults()
			os.Exit(0)
		}
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse flags: %v\n", err)
		os.Exit(1)
	} else if err := theClient.SetEnviron(os.Environ()); err != nil {
		fmt.Fprintf(os.Stderr, "failed to set environment: %v\n", err)
		os.Exit(1)
	} else if req, err := http.NewRequest("GET", fs.Args()[0], nil); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create request: %v\n", err)
		os.Exit(1)
	} else if res, err := theClient.Do(req); err != nil {
		fmt.Fprintf(os.Stderr, "request failed: %v\n", err)
		os.Exit(1)
	} else if _, err := io.Copy(os.Stdout, res.Body); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read response body: %v\n", err)
		os.Exit(1)
	} else if err := res.Body.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to close response body: %v\n", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

// Copyright 2021-2024 The sacloud/apprun-api-go authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"log"
	"net/http"
	"os"

	v1 "github.com/sacloud/apprun-api-go/apis/v1"
	"github.com/sacloud/apprun-api-go/fake"
)

var _ v1.ServerInterface = (*Server)(nil)

type Server struct {
	Engine *fake.Engine
}

// Handler returns an http.Handler built from the generated server wrapper.
// It wires simple recovery and optional logging middlewares based on env vars.
func (s *Server) Handler() http.Handler {
	var middlewares []v1.MiddlewareFunc

	// recovery middleware
	middlewares = append(middlewares, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	})

	// optional logging middleware
	if os.Getenv("APPRUN_SERVER_LOGGING") != "" {
		middlewares = append(middlewares, func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path) //nolint:gosec
				next.ServeHTTP(w, r)
			})
		})
	}

	// create base mux so we can register non-API endpoints like /ping
	m := http.NewServeMux()
	m.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pong"))
	})

	return v1.HandlerWithOptions(s, v1.StdHTTPServerOptions{Middlewares: middlewares, BaseRouter: m})
}

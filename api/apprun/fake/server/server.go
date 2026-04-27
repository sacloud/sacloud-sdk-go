// Copyright 2021-2026 The sacloud/apprun-api-go authors
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
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	v1 "github.com/sacloud/apprun-api-go/apis/v1"
	"github.com/sacloud/apprun-api-go/fake"
)

type Server struct {
	Engine *fake.Engine
}

type Middleware func(http.Handler) http.Handler

// Handler returns an http.Handler built from the generated server wrapper.
// It wires simple recovery and optional logging middlewares based on env vars.
func (s *Server) Handler() http.Handler {
	var middlewares []Middleware

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

	apiHandler := chain(middlewares, http.HandlerFunc(s.route))
	m.Handle("/", apiHandler)

	return m
}

func chain(middlewares []Middleware, next http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		next = middlewares[i](next)
	}
	return next
}

func (s *Server) route(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	if path == "" {
		http.NotFound(w, r)
		return
	}
	parts := strings.Split(path, "/")

	if parts[0] == "user" && len(parts) == 1 {
		switch r.Method {
		case http.MethodGet:
			s.GetUser(w, r)
		case http.MethodPost:
			s.PostUser(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	}

	if parts[0] != "applications" {
		http.NotFound(w, r)
		return
	}

	if len(parts) == 1 {
		switch r.Method {
		case http.MethodGet:
			params, err := listApplicationsParamsFromQuery(r.URL.Query())
			if err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			s.ListApplications(w, r, params)
		case http.MethodPost:
			s.PostApplication(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	}

	appID := parts[1]

	if len(parts) == 2 {
		switch r.Method {
		case http.MethodGet:
			s.GetApplication(w, r, appID)
		case http.MethodPatch:
			s.PatchApplication(w, r, appID)
		case http.MethodDelete:
			s.DeleteApplication(w, r, appID)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	}

	if len(parts) == 3 {
		switch parts[2] {
		case "status":
			if r.Method == http.MethodGet {
				s.GetApplicationStatus(w, r, appID)
				return
			}
		case "traffics":
			switch r.Method {
			case http.MethodGet:
				s.ListApplicationTraffics(w, r, appID)
			case http.MethodPut:
				s.PutApplicationTraffic(w, r, appID)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
			return
		case "packet_filter":
			if _, err := uuid.Parse(appID); err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			switch r.Method {
			case http.MethodGet:
				s.GetPacketFilter(w, r, appID)
			case http.MethodPatch:
				s.PatchPacketFilter(w, r, appID)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
			return
		case "versions":
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			params, err := listApplicationVersionsParamsFromQuery(r.URL.Query())
			if err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			s.ListApplicationVersions(w, r, appID, params)
			return
		}
	}

	if len(parts) == 4 && parts[2] == "versions" {
		versionID := parts[3]
		switch r.Method {
		case http.MethodGet:
			s.GetApplicationVersion(w, r, appID, versionID)
		case http.MethodDelete:
			s.DeleteApplicationVersion(w, r, appID, versionID)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	}

	if len(parts) == 5 && parts[2] == "versions" && parts[4] == "status" {
		if r.Method == http.MethodGet {
			s.GetApplicationVersionStatus(w, r, appID, parts[3])
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	http.NotFound(w, r)
}

func listApplicationsParamsFromQuery(q url.Values) (v1.ListApplicationsParams, error) {
	params := v1.ListApplicationsParams{}
	if v := q.Get("page_num"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return params, err
		}
		params.PageNum = v1.NewOptInt(parsed)
	}
	if v := q.Get("page_size"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return params, err
		}
		params.PageSize = v1.NewOptInt(parsed)
	}
	if v := q.Get("sort_field"); v != "" {
		params.SortField = v1.NewOptString(v)
	}
	if v := q.Get("sort_order"); v != "" {
		switch strings.ToLower(v) {
		case "asc":
			params.SortOrder = v1.NewOptListApplicationsSortOrder(v1.ListApplicationsSortOrderAsc)
		case "desc":
			params.SortOrder = v1.NewOptListApplicationsSortOrder(v1.ListApplicationsSortOrderDesc)
		default:
			return params, fmt.Errorf("invalid sort_order: %s", v)
		}
	}
	return params, nil
}

func listApplicationVersionsParamsFromQuery(q url.Values) (v1.ListApplicationVersionsParams, error) {
	params := v1.ListApplicationVersionsParams{}
	if v := q.Get("page_num"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return params, err
		}
		params.PageNum = v1.NewOptInt(parsed)
	}
	if v := q.Get("page_size"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return params, err
		}
		params.PageSize = v1.NewOptInt(parsed)
	}
	if v := q.Get("sort_field"); v != "" {
		params.SortField = v1.NewOptString(v)
	}
	if v := q.Get("sort_order"); v != "" {
		switch strings.ToLower(v) {
		case "asc":
			params.SortOrder = v1.NewOptListApplicationVersionsSortOrder(v1.ListApplicationVersionsSortOrderAsc)
		case "desc":
			params.SortOrder = v1.NewOptListApplicationVersionsSortOrder(v1.ListApplicationVersionsSortOrderDesc)
		default:
			return params, fmt.Errorf("invalid sort_order: %s", v)
		}
	}
	return params, nil
}

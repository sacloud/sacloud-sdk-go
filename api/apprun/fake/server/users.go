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
	"net/http"

	v1 "github.com/sacloud/sacloud-sdk-go/api/apprun/apis/v1"
)

// ReadUser returns user information
// (GET /user)
func (s *Server) ReadUser(w http.ResponseWriter, r *http.Request) {
	if err := s.Engine.GetUser(); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, &v1.HandlerReadUser{
		Limit: v1.HandlerReadUserLimit{ApplicationCount: 0},
	})
}

// CreateUser creates a user
// (POST /user)
func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	if err := s.Engine.CreateUser(); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusCreated, &v1.HandlerCreateUser{
		Limit: v1.HandlerCreateUserLimit{ApplicationCount: 0},
	})
}

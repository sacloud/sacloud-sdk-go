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

	v1 "github.com/sacloud/apprun-api-go/apis/v1"
)

// GetUser returns user information
// (GET /user)
func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	if err := s.Engine.GetUser(); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, &v1.HandlerGetUser{
		Limit: v1.HandlerGetUserLimit{ApplicationCount: 0},
	})
}

// PostUser creates a user
// (POST /user)
func (s *Server) PostUser(w http.ResponseWriter, r *http.Request) {
	if err := s.Engine.CreateUser(); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusCreated, &v1.HandlerPostUser{
		Limit: v1.HandlerPostUserLimit{ApplicationCount: 0},
	})
}

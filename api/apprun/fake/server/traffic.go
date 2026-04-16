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
	"encoding/json"
	"net/http"

	v1 "github.com/sacloud/apprun-api-go/apis/v1"
)

// ListApplicationTraffics returns traffics for application
// (GET /applications/{id}/traffics)
func (s *Server) ListApplicationTraffics(w http.ResponseWriter, r *http.Request, id string) {
	ts, err := s.Engine.ListTraffics(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, ts)
}

// PutApplicationTraffic updates traffics for application
// (PUT /applications/{id}/traffics)
func (s *Server) PutApplicationTraffic(w http.ResponseWriter, r *http.Request, id string) {
	paramJSON := &v1.PutTrafficsBody{}
	if err := json.NewDecoder(r.Body).Decode(paramJSON); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	ut, err := s.Engine.UpdateTraffic(id, paramJSON)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, &ut)
}

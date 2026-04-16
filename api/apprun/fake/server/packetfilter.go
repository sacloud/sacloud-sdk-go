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

	openapi_types "github.com/oapi-codegen/runtime/types"
	v1 "github.com/sacloud/apprun-api-go/apis/v1"
)

// GetPacketFilter returns packet filter
// (GET /applications/{id}/packet_filter)
func (s *Server) GetPacketFilter(w http.ResponseWriter, r *http.Request, id openapi_types.UUID) {
	pf, err := s.Engine.ReadPacketFilter(id.String())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, pf)
}

// PatchPacketFilter partially updates packet filter
// (PATCH /applications/{id}/packet_filter)
func (s *Server) PatchPacketFilter(w http.ResponseWriter, r *http.Request, id openapi_types.UUID) {
	paramJSON := &v1.PatchPacketFilter{}
	if err := json.NewDecoder(r.Body).Decode(paramJSON); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	ut, err := s.Engine.UpdatePacketFilter(id.String(), paramJSON)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, &ut)
}

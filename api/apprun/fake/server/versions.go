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
	"net/http"

	v1 "github.com/sacloud/apprun-api-go/apis/v1"
)

var (
	versionDefaultPageNum   = 1
	versionDefaultPageSize  = 50
	versionDefaultSortField = "created_at"
	versionDefaultSortOrder = v1.ListApplicationVersionsParamsSortOrderDesc
)

// ListApplicationVersions returns versions for application
// (GET /applications/{id}/versions)
func (s *Server) ListApplicationVersions(w http.ResponseWriter, r *http.Request, id string, params v1.ListApplicationVersionsParams) {
	if params.PageNum == nil {
		params.PageNum = &versionDefaultPageNum
	}
	if params.PageSize == nil {
		params.PageSize = &versionDefaultPageSize
	}
	if params.SortField == nil {
		params.SortField = &versionDefaultSortField
	}
	if params.SortOrder == nil {
		params.SortOrder = &versionDefaultSortOrder
	}

	versions, err := s.Engine.ListVersions(id, params)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, versions)
}

// GetApplicationVersion returns a specific version
// (GET /applications/{id}/versions/{version_id})
func (s *Server) GetApplicationVersion(w http.ResponseWriter, r *http.Request, id string, versionId string) {
	v, err := s.Engine.ReadVersion(id, versionId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, v)
}

// GetApplicationVersionStatus returns status for a specific version
// (GET /applications/{id}/versions/{version_id}/status)
func (s *Server) GetApplicationVersionStatus(w http.ResponseWriter, r *http.Request, id string, versionId string) {
	status, err := s.Engine.ReadVersionStatus(id, versionId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, status)
}

// DeleteApplicationVersion deletes a version
// (DELETE /applications/{id}/versions/{version_id})
func (s *Server) DeleteApplicationVersion(w http.ResponseWriter, r *http.Request, id string, versionId string) {
	if err := s.Engine.DeleteVersion(id, versionId); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

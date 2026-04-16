// Copyright 2022-2025 The sacloud/webaccel-service-go authors
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

package acl

import (
	"context"
	"net/http"

	"github.com/sacloud/services/helper"
	"github.com/sacloud/webaccel-api-go"
)

func (s *Service) Read(req *ReadRequest) (*webaccel.ACLResult, error) {
	return s.ReadWithContext(context.Background(), req)
}

func (s *Service) ReadWithContext(ctx context.Context, req *ReadRequest) (*webaccel.ACLResult, error) {
	if err := helper.ValidateStruct(s, req); err != nil {
		return nil, err
	}

	client := webaccel.NewOp(s.client)
	acl, err := client.ReadACL(ctx, req.SiteId)
	if err != nil {
		return nil, err
	}
	if acl != nil && acl.ACL != "" {
		return acl, nil
	}

	return nil, webaccel.NewAPIError(http.MethodGet, nil, http.StatusNotFound, nil)
}

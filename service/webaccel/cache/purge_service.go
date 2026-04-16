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

package cache

import (
	"context"
	"fmt"

	"github.com/sacloud/services/helper"
	"github.com/sacloud/webaccel-api-go"
)

func (s *Service) Purge(req *PurgeRequest) error {
	return s.PurgeWithContext(context.Background(), req)
}

func (s *Service) PurgeWithContext(ctx context.Context, req *PurgeRequest) error {
	if err := helper.ValidateStruct(s, req); err != nil {
		return err
	}
	client := webaccel.NewOp(s.client)

	switch {
	case len(req.URL) > 0:
		_, err := client.DeleteCache(ctx, &webaccel.DeleteCacheRequest{URL: req.URL})
		return err
	case req.Domain != "":
		return client.DeleteAllCache(ctx, &webaccel.DeleteAllCacheRequest{Domain: req.Domain})
	default:
		return fmt.Errorf("invalid request: %v", req) // 到達しない
	}
}

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

package site

import (
	"context"

	"github.com/sacloud/webaccel-api-go"
)

func (s *Service) Find(req *FindRequest) ([]*webaccel.Site, error) {
	return s.FindWithContext(context.Background(), req)
}

func (s *Service) FindWithContext(ctx context.Context, _ *FindRequest) ([]*webaccel.Site, error) {
	client := webaccel.NewOp(s.client)
	found, err := client.List(ctx)
	if err != nil {
		return nil, err
	}

	return found.Sites, nil
}

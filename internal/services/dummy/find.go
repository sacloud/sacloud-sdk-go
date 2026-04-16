// Copyright 2022-2025 The sacloud/services Authors
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

package dummy

import (
	"context"

	"github.com/sacloud/services/meta"
)

func (s *Service) Find(req *FindRequest) ([]*FindResult, error) {
	return s.FindWithContext(context.Background(), req)
}

func (s *Service) FindWithContext(ctx context.Context, req *FindRequest) ([]*FindResult, error) {
	return []*FindResult{
		{Dummy: "result1"},
		{Dummy: "result2"},
		{Dummy: "result3"},
	}, nil
}

type FindRequest struct {
	Field1 string `validate:"required"`
	Field2 string `validate:"omitempty,option2" meta:",options=option2"`
}

func init() {
	optionDefs = append(optionDefs,
		&meta.Option{Key: "option2", Values: []string{"o3", "o4"}},
	)
}

func (req *FindRequest) Initialize() {
	// 初期値はここで設定する
	req.Field1 = "init"
	req.Field2 = "init"
}

type FindResult struct {
	Dummy string
}

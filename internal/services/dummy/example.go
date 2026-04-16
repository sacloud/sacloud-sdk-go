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

	"github.com/sacloud/services/examples"
	"github.com/sacloud/services/meta"
)

func (s *Service) Example(req *ExampleRequest) ([]*ExampleRequest, error) {
	return s.ExampleWithContext(context.Background(), req)
}

func (s *Service) ExampleWithContext(ctx context.Context, req *ExampleRequest) ([]*ExampleRequest, error) {
	return []*ExampleRequest{req}, nil
}

type ExampleRequest struct {
	Field1  string
	Field2  string
	Options string `meta:",options=example_options"` // example_optionsはなるべく同一ファイル内でoptionDefsに定義する
}

var exampleOptions = &meta.Option{Key: "example_options", Values: []string{"example1", "example2"}}

func init() {
	optionDefs = append(optionDefs, exampleOptions)
}

func (req *ExampleRequest) Examples() interface{} {
	return &ExampleRequest{
		Field1:  examples.Id,
		Field2:  examples.Description,
		Options: exampleOptions.String(), // オプションを持つ場合、meta.OptionのString()またはexamples.OptionsString()で指定する
	}
}

// Copyright 2022-2025 The sacloud/iaas-service-go Authors
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

package esme

import (
	"context"
	"fmt"

	"github.com/sacloud/iaas-api-go"
)

func (s *Service) SendMessage(req *SendMessageRequest) (*iaas.ESMESendMessageResult, error) {
	return s.SendMessageWithContext(context.Background(), req)
}

func (s *Service) SendMessageWithContext(ctx context.Context, req *SendMessageRequest) (*iaas.ESMESendMessageResult, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	client := iaas.NewESMEOp(s.caller)
	_, err := client.Read(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("reading ESME[%s] failed: %s", req.ID, err)
	}

	params := req.ToRequestParameter()
	switch p := params.(type) {
	case *iaas.ESMESendMessageWithGeneratedOTPRequest:
		return client.SendMessageWithGeneratedOTP(ctx, req.ID, p)
	case *iaas.ESMESendMessageWithInputtedOTPRequest:
		return client.SendMessageWithInputtedOTP(ctx, req.ID, p)
	default:
		return nil, fmt.Errorf("processing request parameter failed")
	}
}

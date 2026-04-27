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

package authstatus

import (
	"context"

	"github.com/sacloud/iaas-api-go"
)

// Deprecated: Read is deprecated. Use AuthOp.GetAuth in iam-api-go instead.
// See: https://pkg.go.dev/github.com/sacloud/iam-api-go/apis/auth#AuthAPI
func (s *Service) Read() (*iaas.AuthStatus, error) {
	return s.ReadWithContext(context.Background())
}

// Deprecated: ReadWithContext is deprecated. Use AuthOp.GetAuth in iam-api-go instead.
// See: https://pkg.go.dev/github.com/sacloud/iam-api-go/apis/auth#AuthAPI
func (s *Service) ReadWithContext(ctx context.Context) (*iaas.AuthStatus, error) {
	client := iaas.NewAuthStatusOp(s.caller)
	return client.Read(ctx)
}

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

package coupon

import (
	"context"
	"errors"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iam-api-go"
	"github.com/sacloud/iam-api-go/apis/auth"
)

func (s *Service) List() ([]*iaas.Coupon, error) {
	return s.ListWithContext(context.Background())
}

func (s *Service) ListWithContext(ctx context.Context) ([]*iaas.Coupon, error) {
	iamClient, err := iam.NewClient(s.saclient)
	if err != nil {
		return nil, err
	}

	// check auth status
	authOp := auth.NewAuthOp(iamClient)
	authContext, err := authOp.ReadAuthContext(ctx)
	if err != nil {
		return nil, err
	}
	if authContext.LimitedToProjectID.IsNull() {
		return nil, errors.New("LimitedToProjectID is nil")
	}
	projectId := types.ID(authContext.LimitedToProjectID.Value)

	couponOp := iaas.NewCouponOp(s.caller)
	found, err := couponOp.Find(ctx, projectId)
	if err != nil {
		return nil, err
	}
	return found.Coupons, nil
}

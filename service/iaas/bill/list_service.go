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

package bill

import (
	"context"
	"errors"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iam-api-go"
	"github.com/sacloud/iam-api-go/apis/auth"
)

func (s *Service) List(req *ListRequest) ([]*iaas.Bill, error) {
	return s.ListWithContext(context.Background(), req)
}

func (s *Service) ListWithContext(ctx context.Context, req *ListRequest) ([]*iaas.Bill, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	billOp := iaas.NewBillOp(s.caller)
	iamClient, err := iam.NewClient(s.saclient)
	if err != nil {
		return nil, err
	}
	authOp := auth.NewAuthOp(iamClient)

	authContext, err := authOp.ReadAuthContext(ctx)
	if err != nil {
		return nil, err
	}
	if authContext.LimitedToProjectID.IsNull() {
		return nil, errors.New("LimitedToProjectID is nil")
	}
	projectId := types.ID(authContext.LimitedToProjectID.Value)

	// get bills
	var bills []*iaas.Bill
	switch {
	case req.Month > 0:
		res, err := billOp.ByContractYearMonth(ctx, projectId, req.Year, req.Month)
		if err != nil {
			return nil, err
		}
		bills = res.Bills
	case req.Year > 0:
		res, err := billOp.ByContractYear(ctx, projectId, req.Year)
		if err != nil {
			return nil, err
		}
		bills = res.Bills
	default:
		res, err := billOp.ByContract(ctx, projectId)
		if err != nil {
			return nil, err
		}
		bills = res.Bills
	}
	return bills, nil
}

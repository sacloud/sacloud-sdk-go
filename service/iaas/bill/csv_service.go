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
)

func (s *Service) Csv(req *CsvRequest) (*iaas.BillDetailCSV, error) {
	return s.CsvWithContext(context.Background(), req)
}

func (s *Service) CsvWithContext(ctx context.Context, req *CsvRequest) (*iaas.BillDetailCSV, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	billOp := iaas.NewBillOp(s.caller)
	authOp := iaas.NewAuthStatusOp(s.caller)

	// check auth status
	auth, err := authOp.Read(ctx)
	if err != nil {
		return nil, err
	}
	if auth.AccountID.IsEmpty() {
		return nil, errors.New("invalid account id")
	}

	// get latest bill ID if empty
	billID := req.ID
	if billID.IsEmpty() {
		bills, err := billOp.ByContract(ctx, auth.AccountID)
		if err != nil {
			return nil, err
		}
		if len(bills.Bills) == 0 {
			return nil, iaas.NewNoResultsError()
		}
		billID = bills.Bills[0].ID
	}

	return billOp.DetailsCSV(ctx, auth.MemberCode, billID)
}

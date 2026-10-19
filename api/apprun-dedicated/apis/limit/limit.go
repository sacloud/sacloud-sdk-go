// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package limit

import (
	"context"

	v1 "github.com/sacloud/sacloud-sdk-go/api/apprun-dedicated/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/api/apprun-dedicated/common"
)

type LimitAPI interface {
	Read(ctx context.Context) (*v1.Limit, error)
}

type LimitOp struct{ *v1.Client }

func NewLimitOp(client *v1.Client) *LimitOp { return &LimitOp{Client: client} }

func (op *LimitOp) Read(ctx context.Context) (limit *v1.Limit, err error) {
	res, err := common.ErrorFromDecodedResponse("Limit.Read", func() (*v1.ReadLimitsResponse, error) {
		return op.Client.ReadLimits(ctx)
	})
	if res != nil {
		limit = &res.Limit
	}
	return
}

var _ LimitAPI = (*LimitOp)(nil)

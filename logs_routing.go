// Copyright 2025- The sacloud/monitoring-suite-api-go Authors
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

package monitoringsuite

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
)

type LogRoutingAPI interface {
	List(ctx context.Context, params LogsRoutingsListParams) ([]v1.LogRouting, error)
	Create(ctx context.Context, params LogsRoutingCreateParams) (*v1.LogRouting, error)
	Read(ctx context.Context, id uuid.UUID) (*v1.LogRouting, error)
	Update(ctx context.Context, id uuid.UUID, params LogsRoutingUpdateParams) (*v1.LogRouting, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

var _ LogRoutingAPI = (*logRoutingOp)(nil)

type logRoutingOp struct {
	client *v1.Client
}

func NewLogRoutingOp(client *v1.Client) LogRoutingAPI {
	return &logRoutingOp{client: client}
}

type LogsRoutingsListParams struct {
	Count         *int
	From          *int
	PublisherCode *string
	ResourceID    *int64
	Variant       *string
}

func (op *logRoutingOp) List(ctx context.Context, p LogsRoutingsListParams) (ret []v1.LogRouting, err error) {
	res, err := errorFromDecodedResponse("LogRouting.List", func() (*v1.PaginatedLogRoutingList, error) {
		return op.client.LogsRoutingsList(ctx, v1.LogsRoutingsListParams{
			Count:         intoOpt[v1.OptInt](p.Count),
			From:          intoOpt[v1.OptInt](p.From),
			PublisherCode: intoOpt[v1.OptString](p.PublisherCode),
			ResourceID:    intoOpt[v1.OptInt64](p.ResourceID),
			Variant:       intoOpt[v1.OptString](p.Variant),
		})
	})
	if err == nil {
		ret = res.GetResults()
	}
	return
}

type LogsRoutingCreateParams struct {
	PublisherCode string
	ResourceID    *string
	Variant       string
	LogStorageID  string
}

func (op *logRoutingOp) Create(ctx context.Context, params LogsRoutingCreateParams) (*v1.LogRouting, error) {
	res, err := errorFromDecodedResponse("LogRouting.Create", func() (*v1.WrappedLogRouting, error) {
		rid, err := fromStringPtr[v1.OptNilInt64, int64](params.ResourceID)
		if err != nil {
			return nil, fmt.Errorf("LogsRoutingCreateParams.ResourceID: %w", err)
		}
		lid, err := strconv.ParseInt(params.LogStorageID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("LogsRoutingCreateParams.LogStorageID: %w", err)
		}
		request := v1.LogRoutingRequest{
			PublisherCode: params.PublisherCode,
			ResourceID:    rid,
			Variant:       params.Variant,
			LogStorageID:  v1.NewNilInt64(lid),
		}

		return op.client.LogsRoutingsCreate(ctx, &request)
	})
	return unwrapE[*v1.LogRouting](res, err)
}

func (op *logRoutingOp) Read(ctx context.Context, id uuid.UUID) (*v1.LogRouting, error) {
	res, err := errorFromDecodedResponse("LogRouting.Read", func() (*v1.WrappedLogRouting, error) {
		return op.client.LogsRoutingsRetrieve(ctx, v1.LogsRoutingsRetrieveParams{UID: id})
	})
	return unwrapE[*v1.LogRouting](res, err)
}

type LogsRoutingUpdateParams struct {
	PublisherCode *string
	ResourceID    *string
	Variant       *string
	LogStorageID  *string
}

func (op *logRoutingOp) Update(ctx context.Context, id uuid.UUID, params LogsRoutingUpdateParams) (*v1.LogRouting, error) {
	res, err := errorFromDecodedResponse("LogRouting.Update", func() (*v1.WrappedLogRouting, error) {
		rid, err := fromStringPtr[v1.OptNilInt64, int64](params.ResourceID)
		if err != nil {
			return nil, fmt.Errorf("LogsRoutingUpdateParams.ResourceID: %w", err)
		}
		lid, err := fromStringPtr[v1.OptNilInt64, int64](params.LogStorageID)
		if err != nil {
			return nil, fmt.Errorf("LogsRoutingUpdateParams.LogStorageID: %w", err)
		}
		return op.client.LogsRoutingsPartialUpdate(ctx, v1.NewOptPatchedLogRoutingRequest(v1.PatchedLogRoutingRequest{
			PublisherCode: intoOpt[v1.OptString](params.PublisherCode),
			ResourceID:    rid,
			Variant:       intoOpt[v1.OptString](params.Variant),
			LogStorageID:  lid,
		}), v1.LogsRoutingsPartialUpdateParams{UID: id})
	})
	return unwrapE[*v1.LogRouting](res, err)
}

func (op *logRoutingOp) Delete(ctx context.Context, id uuid.UUID) error {
	return errorFromDecodedResponse1("LogRouting.Delete", func() error {
		return op.client.LogsRoutingsDestroy(ctx, v1.LogsRoutingsDestroyParams{UID: id})
	})
}

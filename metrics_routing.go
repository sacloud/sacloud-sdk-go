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

type MetricsRoutingAPI interface {
	List(ctx context.Context, params MetricsRoutingsListParams) ([]v1.MetricsRouting, error)
	Create(ctx context.Context, params MetricsRoutingCreateParams) (*v1.MetricsRouting, error)
	Read(ctx context.Context, id uuid.UUID) (*v1.MetricsRouting, error)
	Update(ctx context.Context, id uuid.UUID, params MetricsRoutingUpdateParams) (*v1.MetricsRouting, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

var _ MetricsRoutingAPI = (*metricsRoutingOp)(nil)

type metricsRoutingOp struct {
	client *v1.Client
}

func NewMetricsRoutingOp(client *v1.Client) MetricsRoutingAPI {
	return &metricsRoutingOp{client: client}
}

type MetricsRoutingsListParams struct {
	Count         *int
	From          *int
	PublisherCode *string
	ResourceID    *int64
	Variant       *string
}

func (op *metricsRoutingOp) List(ctx context.Context, p MetricsRoutingsListParams) (ret []v1.MetricsRouting, err error) {
	res, err := errorFromDecodedResponse("MetricsRouting.List", func() (*v1.PaginatedMetricsRoutingList, error) {
		return op.client.MetricsRoutingsList(ctx, v1.MetricsRoutingsListParams{
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

type MetricsRoutingCreateParams struct {
	PublisherCode    string
	ResourceID       *string
	Variant          string
	MetricsStorageID string
}

func (op *metricsRoutingOp) Create(ctx context.Context, params MetricsRoutingCreateParams) (*v1.MetricsRouting, error) {
	res, err := errorFromDecodedResponse("MetricsRouting.Create", func() (*v1.WrappedMetricsRouting, error) {
		if rid, err := fromStringPtr[v1.OptNilInt64, int64](params.ResourceID); err != nil {
			return nil, fmt.Errorf("MetricsRoutingCreateParams.ResourceID: %w", err)
		} else if mid, err := strconv.ParseInt(params.MetricsStorageID, 10, 64); err != nil {
			return nil, fmt.Errorf("MetricsRoutingCreateParams.MetricsStorageID: %w", err)
		} else {
			req := v1.MetricsRouting{
				PublisherCode:    intoOpt[v1.OptString](&params.PublisherCode),
				ResourceID:       rid,
				Variant:          params.Variant,
				MetricsStorageID: intoOptNil[v1.OptNilInt64](&mid),
			}

			// prevent ogen error (encoder is not accepting empty struct)
			req.Publisher.SetFake()
			req.MetricsStorage.SetFake()
			req.Publisher.SetVariants(make([]v1.PublisherVariant, 0))
			req.MetricsStorage.SetTags(make([]string, 0))

			return op.client.MetricsRoutingsCreate(ctx, &req)
		}
	})
	return unwrapE[*v1.MetricsRouting](res, err)
}

func (op *metricsRoutingOp) Read(ctx context.Context, id uuid.UUID) (*v1.MetricsRouting, error) {
	res, err := errorFromDecodedResponse("MetricsRouting.Read", func() (*v1.WrappedMetricsRouting, error) {
		return op.client.MetricsRoutingsRetrieve(ctx, v1.MetricsRoutingsRetrieveParams{UID: id})
	})
	return unwrapE[*v1.MetricsRouting](res, err)
}

type MetricsRoutingUpdateParams struct {
	PublisherCode    *string
	ResourceID       *string
	Variant          *string
	MetricsStorageID *string
}

func (op *metricsRoutingOp) Update(ctx context.Context, id uuid.UUID, params MetricsRoutingUpdateParams) (*v1.MetricsRouting, error) {
	res, err := errorFromDecodedResponse("MetricsRouting.Update", func() (*v1.WrappedMetricsRouting, error) {
		if rid, err := fromStringPtr[v1.OptNilInt64, int64](params.ResourceID); err != nil {
			return nil, fmt.Errorf("MetricsRoutingUpdateParams.ResourceID: %w", err)
		} else if mid, err := fromStringPtr[v1.OptNilInt64, int64](params.MetricsStorageID); err != nil {
			return nil, fmt.Errorf("MetricsRoutingUpdateParams.MetricsStorageID: %w", err)
		} else {
			return op.client.MetricsRoutingsPartialUpdate(ctx, v1.NewOptPatchedMetricsRouting(v1.PatchedMetricsRouting{
				PublisherCode:    intoOpt[v1.OptString](params.PublisherCode),
				ResourceID:       rid,
				Variant:          intoOpt[v1.OptString](params.Variant),
				MetricsStorageID: mid,
			}), v1.MetricsRoutingsPartialUpdateParams{UID: id})
		}
	})
	return unwrapE[*v1.MetricsRouting](res, err)
}

func (op *metricsRoutingOp) Delete(ctx context.Context, id uuid.UUID) error {
	return errorFromDecodedResponse1("MetricsRouting.Delete", func() error {
		return op.client.MetricsRoutingsDestroy(ctx, v1.MetricsRoutingsDestroyParams{UID: id})
	})
}

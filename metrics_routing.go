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
	"net/http"
	"strconv"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	ogen "github.com/ogen-go/ogen/validate"
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

func (op *metricsRoutingOp) List(ctx context.Context, p MetricsRoutingsListParams) ([]v1.MetricsRouting, error) {
	params := v1.MetricsRoutingsListParams{
		Count:         intoOpt[v1.OptInt](p.Count),
		From:          intoOpt[v1.OptInt](p.From),
		PublisherCode: intoOpt[v1.OptString](p.PublisherCode),
		ResourceID:    intoOpt[v1.OptInt64](p.ResourceID),
		Variant:       intoOpt[v1.OptString](p.Variant),
	}
	resp, err := op.client.MetricsRoutingsList(ctx, params)
	if err != nil {
		return nil, NewAPIError("MetricsRouting.List", 0, err)
	}
	return resp.Results, nil
}

type MetricsRoutingCreateParams struct {
	PublisherCode    string
	ResourceID       *string
	Variant          string
	MetricsStorageID string
}

func (op *metricsRoutingOp) Create(ctx context.Context, params MetricsRoutingCreateParams) (*v1.MetricsRouting, error) {
	rid, err := fromStringPtr[v1.OptNilInt64, int64](params.ResourceID)
	if err != nil {
		return nil, NewError("MetricsRouting.Create", err)
	}
	mid, err := strconv.ParseInt(params.MetricsStorageID, 10, 64)
	if err != nil {
		return nil, NewError("MetricsRouting.Create", err)
	}
	req := v1.MetricsRouting{
		PublisherCode:    intoOpt[v1.OptString](&params.PublisherCode),
		ResourceID:       rid,
		Variant:          params.Variant,
		MetricsStorageID: intoOptNil[v1.OptNilInt64](&mid),
		Publisher:        v1.Publisher{},
		MetricsStorage:   v1.MetricsStorage{},
	}

	// prevent ogen error (encoder is not accepting empty struct)
	req.Publisher.SetFake()
	req.MetricsStorage.SetFake()
	req.Publisher.SetVariants([]v1.PublisherVariant{})
	req.MetricsStorage.SetTags([]string{})

	resp, err := op.client.MetricsRoutingsCreate(ctx, &req)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("MetricsRouting.Create", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return nil, NewAPIError("MetricsRouting.Create", e.StatusCode, errors.Wrap(err, "invalid parameter, or the metrics storage cannot be routed"))
		default:
			return nil, NewAPIError("MetricsRouting.Create", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("MetricsRouting.Create", 0, err)
	} else {
		ret := new(v1.MetricsRouting)
		return Unwrap(ret, resp)
	}
}

func (op *metricsRoutingOp) Read(ctx context.Context, id uuid.UUID) (*v1.MetricsRouting, error) {
	resp, err := op.client.MetricsRoutingsRetrieve(ctx, v1.MetricsRoutingsRetrieveParams{UID: id})
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("MetricsRouting.Read", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return nil, NewAPIError("MetricsRouting.Read", e.StatusCode, errors.Wrap(err, "invalid parameter, or invalid metrics storage"))
		default:
			return nil, NewAPIError("MetricsRouting.Read", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("MetricsRouting.Read", 0, err)
	} else {
		ret := new(v1.MetricsRouting)
		return Unwrap(ret, resp)
	}
}

type MetricsRoutingUpdateParams struct {
	PublisherCode    *string
	ResourceID       *string
	Variant          *string
	MetricsStorageID *string
}

func (op *metricsRoutingOp) Update(ctx context.Context, id uuid.UUID, params MetricsRoutingUpdateParams) (*v1.MetricsRouting, error) {
	rid, err := fromStringPtr[v1.OptNilInt64, int64](params.ResourceID)
	if err != nil {
		return nil, NewError("MetricsRouting.Update", err)
	}
	mid, err := fromStringPtr[v1.OptNilInt64, int64](params.MetricsStorageID)
	if err != nil {
		return nil, NewError("MetricsRouting.Update", err)
	}
	patch := v1.PatchedMetricsRouting{
		PublisherCode:    intoOpt[v1.OptString](params.PublisherCode),
		ResourceID:       rid,
		Variant:          intoOpt[v1.OptString](params.Variant),
		MetricsStorageID: mid,
	}
	req := v1.NewOptPatchedMetricsRouting(patch)
	resp, err := op.client.MetricsRoutingsPartialUpdate(ctx, req, v1.MetricsRoutingsPartialUpdateParams{UID: id})
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("MetricsRouting.Update", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return nil, NewAPIError("MetricsRouting.Update", e.StatusCode, errors.Wrap(err, "invalid parameter, or invalid metrics storage"))
		default:
			return nil, NewAPIError("MetricsRouting.Update", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("MetricsRouting.Update", 0, err)
	} else {
		ret := new(v1.MetricsRouting)
		return Unwrap(ret, resp)
	}
}

func (op *metricsRoutingOp) Delete(ctx context.Context, id uuid.UUID) error {
	err := op.client.MetricsRoutingsDestroy(ctx, v1.MetricsRoutingsDestroyParams{UID: id})
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return NewAPIError("MetricsRouting.Delete", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return NewAPIError("MetricsRouting.Delete", e.StatusCode, errors.Wrap(err, "the request resource is not eligible for deletion"))
		default:
			return NewAPIError("MetricsRouting.Delete", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return NewAPIError("MetricsRouting.Delete", 0, err)
	}
	return nil
}

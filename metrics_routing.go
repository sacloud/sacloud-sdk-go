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

	"github.com/go-faster/errors"
	ogen "github.com/ogen-go/ogen/validate"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
)

type MetricsRoutingAPI interface {
	List(ctx context.Context, params v1.MetricsRoutingsListParams) ([]v1.MetricsRouting, error)
	Create(ctx context.Context, request v1.MetricsRouting) (*v1.MetricsRouting, error)
	Read(ctx context.Context, id int64) (*v1.MetricsRouting, error)
	Update(ctx context.Context, id int64, request *v1.MetricsRouting) (*v1.MetricsRouting, error)
	Delete(ctx context.Context, id int64) error
}

var _ MetricsRoutingAPI = (*metricsRoutingOp)(nil)

type metricsRoutingOp struct {
	client *v1.Client
}

func NewMetricsRoutingOp(client *v1.Client) MetricsRoutingAPI {
	return &metricsRoutingOp{client: client}
}

func (op *metricsRoutingOp) List(ctx context.Context, params v1.MetricsRoutingsListParams) ([]v1.MetricsRouting, error) {
	resp, err := op.client.MetricsRoutingsList(ctx, params)
	if err != nil {
		return nil, err
	}
	return resp.Results, nil
}

func (op *metricsRoutingOp) Create(ctx context.Context, request v1.MetricsRouting) (*v1.MetricsRouting, error) {
	resp, err := op.client.MetricsRoutingsCreate(ctx, &request)
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

func (op *metricsRoutingOp) Read(ctx context.Context, id int64) (*v1.MetricsRouting, error) {
	resp, err := op.client.MetricsRoutingsRetrieve(ctx, v1.MetricsRoutingsRetrieveParams{ID: id})
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

func (op *metricsRoutingOp) Update(ctx context.Context, id int64, request *v1.MetricsRouting) (*v1.MetricsRouting, error) {
	resp, err := op.client.MetricsRoutingsUpdate(ctx, request, v1.MetricsRoutingsUpdateParams{ID: id})
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

func (op *metricsRoutingOp) Delete(ctx context.Context, id int64) error {
	err := op.client.MetricsRoutingsDestroy(ctx, v1.MetricsRoutingsDestroyParams{ID: id})
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

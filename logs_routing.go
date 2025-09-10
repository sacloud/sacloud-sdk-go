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
	ogen "github.com/ogen-go/ogen/validate"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
)

type LogRoutingAPI interface {
	List(ctx context.Context, params v1.LogsRoutingsListParams) ([]v1.LogRouting, error)
	Create(ctx context.Context, request v1.LogRouting) (*v1.LogRouting, error)
	Read(ctx context.Context, id string) (*v1.LogRouting, error)
	Update(ctx context.Context, id string, request *v1.LogRouting) (*v1.LogRouting, error)
	Delete(ctx context.Context, id string) error
}

var _ LogRoutingAPI = (*logRoutingOp)(nil)

type logRoutingOp struct {
	client *v1.Client
}

func NewLogRoutingOp(client *v1.Client) LogRoutingAPI {
	return &logRoutingOp{client: client}
}

func (op *logRoutingOp) List(ctx context.Context, params v1.LogsRoutingsListParams) ([]v1.LogRouting, error) {
	resp, err := op.client.LogsRoutingsList(ctx, params)
	if err != nil {
		return nil, NewAPIError("LogRouting.List", 0, err)
	}
	return resp.Results, nil
}

func (op *logRoutingOp) Create(ctx context.Context, request v1.LogRouting) (*v1.LogRouting, error) {
	resp, err := op.client.LogsRoutingsCreate(ctx, &request)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("LogRouting.Create", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return nil, NewAPIError("LogRouting.Create", e.StatusCode, errors.Wrap(err, "invalid parameter, or the log storage cannot be routed"))
		default:
			return nil, NewAPIError("LogRouting.Create", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("LogRouting.Create", 0, err)
	} else {
		ret := new(v1.LogRouting)
		return Unwrap(ret, resp)
	}
}

func (op *logRoutingOp) Read(ctx context.Context, id string) (*v1.LogRouting, error) {
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, NewAPIError("LogRouting.Read", 0, err)
	}
	resp, err := op.client.LogsRoutingsRetrieve(ctx, v1.LogsRoutingsRetrieveParams{ID: intId})
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("LogRouting.Read", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return nil, NewAPIError("LogRouting.Read", e.StatusCode, errors.Wrap(err, "invalid parameter, or invalid log storage"))
		default:
			return nil, NewAPIError("LogRouting.Read", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("LogRouting.Read", 0, err)
	} else {
		ret := new(v1.LogRouting)
		return Unwrap(ret, resp)
	}
}

func (op *logRoutingOp) Update(ctx context.Context, id string, request *v1.LogRouting) (*v1.LogRouting, error) {
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, NewAPIError("LogRouting.Update", 0, err)
	}
	resp, err := op.client.LogsRoutingsUpdate(ctx, request, v1.LogsRoutingsUpdateParams{ID: intId})
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("LogRouting.Update", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return nil, NewAPIError("LogRouting.Update", e.StatusCode, errors.Wrap(err, "invalid parameter, or invalid log storage"))
		default:
			return nil, NewAPIError("LogRouting.Update", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("LogRouting.Update", 0, err)
	} else {
		ret := new(v1.LogRouting)
		return Unwrap(ret, resp)
	}
}

func (op *logRoutingOp) Delete(ctx context.Context, id string) error {
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return NewAPIError("LogRouting.Delete", 0, err)
	}
	err = op.client.LogsRoutingsDestroy(ctx, v1.LogsRoutingsDestroyParams{ID: intId})
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return NewAPIError("LogRouting.Delete", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return NewAPIError("LogRouting.Delete", e.StatusCode, errors.Wrap(err, "the request resource is not eligible for deletion"))
		default:
			return NewAPIError("LogRouting.Delete", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return NewAPIError("LogRouting.Delete", 0, err)
	}
	return nil
}

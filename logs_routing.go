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

func (op *logRoutingOp) List(ctx context.Context, p LogsRoutingsListParams) ([]v1.LogRouting, error) {
	params := v1.LogsRoutingsListParams{
		Count:         intoOpt[v1.OptInt](p.Count),
		From:          intoOpt[v1.OptInt](p.From),
		PublisherCode: intoOpt[v1.OptString](p.PublisherCode),
		ResourceID:    intoOpt[v1.OptInt64](p.ResourceID),
		Variant:       intoOpt[v1.OptString](p.Variant),
	}
	resp, err := op.client.LogsRoutingsList(ctx, params)
	if err != nil {
		return nil, NewAPIError("LogRouting.List", 0, err)
	}
	return resp.Results, nil
}

type LogsRoutingCreateParams struct {
	PublisherCode string
	ResourceID    *string
	Variant       string
	LogStorageID  string
}

func (op *logRoutingOp) Create(ctx context.Context, params LogsRoutingCreateParams) (*v1.LogRouting, error) {
	rid, err := fromStringPtr[v1.OptNilInt64, int64](params.ResourceID)
	if err != nil {
		return nil, errors.Wrap(err, "invalid ResourceID")
	}
	lid, err := strconv.ParseInt(params.LogStorageID, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "invalid LogStorageID")
	}
	request := v1.LogRouting{
		PublisherCode: intoOpt[v1.OptString](&params.PublisherCode),
		ResourceID:    rid,
		Variant:       params.Variant,
		LogStorageID:  intoOptNil[v1.OptNilInt64](&lid),
	}

	// prevent ogen error (encoder is not accepting empty struct)
	request.Publisher.SetFake()
	request.LogStorage.SetFake()
	request.Publisher.SetVariants([]v1.PublisherVariant{})
	request.LogStorage.SetTags([]string{})

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

func (op *logRoutingOp) Read(ctx context.Context, id uuid.UUID) (*v1.LogRouting, error) {
	resp, err := op.client.LogsRoutingsRetrieve(ctx, v1.LogsRoutingsRetrieveParams{UID: id})
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

type LogsRoutingUpdateParams struct {
	PublisherCode *string
	ResourceID    *string
	Variant       *string
	LogStorageID  *string
}

func (op *logRoutingOp) Update(ctx context.Context, id uuid.UUID, params LogsRoutingUpdateParams) (*v1.LogRouting, error) {
	rid, err := fromStringPtr[v1.OptNilInt64, int64](params.ResourceID)
	if err != nil {
		return nil, errors.Wrap(err, "invalid ResourceID")
	}
	lid, err := fromStringPtr[v1.OptNilInt64, int64](params.LogStorageID)
	if err != nil {
		return nil, errors.Wrap(err, "invalid LogStorageID")
	}
	patch := v1.PatchedLogRouting{
		PublisherCode: intoOpt[v1.OptString](params.PublisherCode),
		ResourceID:    rid,
		Variant:       intoOpt[v1.OptString](params.Variant),
		LogStorageID:  lid,
	}
	request := v1.NewOptPatchedLogRouting(patch)
	resp, err := op.client.LogsRoutingsPartialUpdate(ctx, request, v1.LogsRoutingsPartialUpdateParams{UID: id})
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

func (op *logRoutingOp) Delete(ctx context.Context, id uuid.UUID) error {
	err := op.client.LogsRoutingsDestroy(ctx, v1.LogsRoutingsDestroyParams{UID: id})
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

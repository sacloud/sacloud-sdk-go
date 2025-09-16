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

type MetricsStorageAPI interface {
	List(ctx context.Context, count *int, from *int) ([]v1.MetricsStorage, error)
	Create(ctx context.Context, request v1.MetricsStorageCreate) (*v1.MetricsStorage, error)
	Read(ctx context.Context, id string) (*v1.MetricsStorage, error)
	Update(ctx context.Context, id string, request *v1.MetricsStorage) (*v1.MetricsStorage, error)
	Delete(ctx context.Context, id string) error

	ListKeys(ctx context.Context, metricsResourceId string, count *int, from *int) ([]v1.MetricsStorageAccessKey, error)
	CreateKey(ctx context.Context, metricsResourceId string, request *v1.MetricsStorageAccessKey) (*v1.MetricsStorageAccessKey, error)
	ReadKey(ctx context.Context, metricsResourceId string, id string) (*v1.MetricsStorageAccessKey, error)
	UpdateKey(ctx context.Context, metricsResourceId string, id string, request *v1.MetricsStorageAccessKey) (*v1.MetricsStorageAccessKey, error)
	DeleteKey(ctx context.Context, metricsResourceId string, id string) error
}

var _ MetricsStorageAPI = (*metricsStorageOp)(nil)

type metricsStorageOp struct {
	client *v1.Client
}

func NewMetricsStorageOp(client *v1.Client) MetricsStorageAPI {
	return &metricsStorageOp{client: client}
}

func (op *metricsStorageOp) List(ctx context.Context, count *int, from *int) ([]v1.MetricsStorage, error) {
	params := v1.MetricsStoragesListParams{
		Count: intoOpt[v1.OptInt](count),
		From:  intoOpt[v1.OptInt](from),
	}
	result, err := op.client.MetricsStoragesList(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("MetricsStorage.List", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		default:
			return nil, NewAPIError("MetricsStorage.List", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("MetricsStorage.List", 0, err)
	} else {
		return result.GetResults(), nil
	}
}

func (op *metricsStorageOp) Read(ctx context.Context, resourceID string) (*v1.MetricsStorage, error) {
	id, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return nil, NewAPIError("MetricsStorage.Read", 0, err)
	}
	params := v1.MetricsStoragesRetrieveParams{ResourceID: id}
	result, err := op.client.MetricsStoragesRetrieve(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("MetricsStorage.Read", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("MetricsStorage.Read", e.StatusCode, errors.Wrap(err, "metrics tank not found"))
		default:
			return nil, NewAPIError("MetricsStorage.Read", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("MetricsStorage.Read", 0, err)
	} else {
		ret := new(v1.MetricsStorage)
		return Unwrap(ret, result)
	}
}

func (op *metricsStorageOp) Create(ctx context.Context, body v1.MetricsStorageCreate) (*v1.MetricsStorage, error) {
	result, err := op.client.MetricsStoragesCreate(ctx, &body)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("MetricsStorage.Create", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return nil, NewAPIError("MetricsStorage.Create", e.StatusCode, errors.Wrap(err, "invalid parameter, or no space left for a new storage"))
		default:
			return nil, NewAPIError("MetricsStorage.Create", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("MetricsStorage.Create", 0, err)
	} else {
		return result, nil
	}
}

func (op *metricsStorageOp) Update(ctx context.Context, id string, resource *v1.MetricsStorage) (*v1.MetricsStorage, error) {
	rid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, NewAPIError("MetricsStorage.Update", 0, err)
	}
	query := v1.MetricsStoragesUpdateParams{ResourceID: rid}
	body := v1.NewOptMetricsStorage(*resource)
	result, err := op.client.MetricsStoragesUpdate(ctx, body, query)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("MetricsStorage.Update", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return nil, NewAPIError("MetricsStorage.Update", e.StatusCode, errors.Wrap(err, "invalid parameter"))
		default:
			return nil, NewAPIError("MetricsStorage.Update", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("MetricsStorage.Update", 0, err)
	} else {
		ret := new(v1.MetricsStorage)
		return Unwrap(ret, result)
	}
}

func (op *metricsStorageOp) Delete(ctx context.Context, resourceID string) error {
	rid, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return NewAPIError("MetricsStorage.Delete", 0, err)
	}
	params := v1.MetricsStoragesDestroyParams{ResourceID: rid}
	err = op.client.MetricsStoragesDestroy(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return NewAPIError("MetricsStorage.Delete", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return NewAPIError("MetricsStorage.Delete", e.StatusCode, errors.Wrap(err, "the request resource is not eligible for deletion"))
		default:
			return NewAPIError("MetricsStorage.Delete", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return NewAPIError("MetricsStorage.Delete", 0, err)
	}
	return nil
}

func (op *metricsStorageOp) ListKeys(ctx context.Context, metricsResourceId string, count *int, from *int) ([]v1.MetricsStorageAccessKey, error) {
	rid, err := strconv.ParseInt(metricsResourceId, 10, 64)
	if err != nil {
		return nil, NewAPIError("MetricsStorage.ListKeys", 0, err)
	}
	params := v1.MetricsStoragesKeysListParams{
		MetricsResourceID: rid,
		Count:             intoOpt[v1.OptInt](count),
		From:              intoOpt[v1.OptInt](from),
	}
	result, err := op.client.MetricsStoragesKeysList(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("MetricsStorage.ListKeys", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		default:
			return nil, NewAPIError("MetricsStorage.ListKeys", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("MetricsStorage.ListKeys", 0, err)
	} else {
		return result.GetResults(), nil
	}
}

func (op *metricsStorageOp) CreateKey(ctx context.Context, metricsResourceId string, request *v1.MetricsStorageAccessKey) (*v1.MetricsStorageAccessKey, error) {
	rid, err := strconv.ParseInt(metricsResourceId, 10, 64)
	if err != nil {
		return nil, NewAPIError("MetricsStorage.CreateKey", 0, err)
	}
	params := v1.MetricsStoragesKeysCreateParams{MetricsResourceID: rid}
	opt := v1.NewOptMetricsStorageAccessKey(*request)
	result, err := op.client.MetricsStoragesKeysCreate(ctx, opt, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("MetricsStorage.CreateKey", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("MetricsStorage.CreateKey", e.StatusCode, errors.Wrap(err, "requested metrics storage not found"))
		case http.StatusBadRequest:
			return nil, NewAPIError("MetricsStorage.CreateKey", e.StatusCode, errors.Wrap(err, "this metrics storage cannot have a key"))
		default:
			return nil, NewAPIError("MetricsStorage.CreateKey", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("MetricsStorage.CreateKey", 0, err)
	} else {
		ret := new(v1.MetricsStorageAccessKey)
		return Unwrap(ret, result)
	}
}

func (op *metricsStorageOp) ReadKey(ctx context.Context, metricsResourceId string, id string) (*v1.MetricsStorageAccessKey, error) {
	rid, err := strconv.ParseInt(metricsResourceId, 10, 64)
	if err != nil {
		return nil, NewAPIError("MetricsStorage.ReadKey", 0, err)
	}
	kid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, NewAPIError("MetricsStorage.ReadKey", 0, err)
	}
	params := v1.MetricsStoragesKeysRetrieveParams{MetricsResourceID: rid, ID: kid}
	result, err := op.client.MetricsStoragesKeysRetrieve(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("MetricsStorage.ReadKey", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("MetricsStorage.ReadKey", e.StatusCode, errors.Wrap(err, "access key not found"))
		default:
			return nil, NewAPIError("MetricsStorage.ReadKey", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("MetricsStorage.ReadKey", 0, err)
	} else {
		ret := new(v1.MetricsStorageAccessKey)
		return Unwrap(ret, result)
	}
}

func (op *metricsStorageOp) UpdateKey(ctx context.Context, metricsResourceId string, id string, request *v1.MetricsStorageAccessKey) (*v1.MetricsStorageAccessKey, error) {
	rid, err := strconv.ParseInt(metricsResourceId, 10, 64)
	if err != nil {
		return nil, NewAPIError("MetricsStorage.UpdateKey", 0, err)
	}
	kid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, NewAPIError("MetricsStorage.UpdateKey", 0, err)
	}
	params := v1.MetricsStoragesKeysUpdateParams{MetricsResourceID: rid, ID: kid}
	opt := v1.NewOptMetricsStorageAccessKey(*request)
	result, err := op.client.MetricsStoragesKeysUpdate(ctx, opt, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("MetricsStorage.UpdateKey", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return nil, NewAPIError("MetricsStorage.UpdateKey", e.StatusCode, errors.Wrap(err, "invalid parameter"))
		default:
			return nil, NewAPIError("MetricsStorage.UpdateKey", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("MetricsStorage.UpdateKey", 0, err)
	} else {
		ret := new(v1.MetricsStorageAccessKey)
		return Unwrap(ret, result)
	}
}

// DeleteKey deletes an access key for a metrics storage resource.
func (op *metricsStorageOp) DeleteKey(ctx context.Context, metricsResourceId string, id string) error {
	rid, err := strconv.ParseInt(metricsResourceId, 10, 64)
	if err != nil {
		return NewAPIError("MetricsStorage.DeleteKey", 0, err)
	}
	kid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return NewAPIError("MetricsStorage.DeleteKey", 0, err)
	}
	params := v1.MetricsStoragesKeysDestroyParams{MetricsResourceID: rid, ID: kid}
	err = op.client.MetricsStoragesKeysDestroy(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return NewAPIError("MetricsStorage.DeleteKey", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return NewAPIError("MetricsStorage.DeleteKey", e.StatusCode, errors.Wrap(err, "the request resource is not eligible for deletion"))
		default:
			return NewAPIError("MetricsStorage.DeleteKey", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return NewAPIError("MetricsStorage.DeleteKey", 0, err)
	}
	return nil
}

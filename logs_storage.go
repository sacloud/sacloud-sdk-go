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

type LogsStorageAPI interface {
	List(ctx context.Context, params v1.LogsStoragesListParams) ([]v1.LogStorage, error)
	Create(ctx context.Context, params v1.LogStorageCreate) (*v1.LogStorage, error)
	Read(ctx context.Context, id string) (*v1.LogStorage, error)
	Update(ctx context.Context, id string, request *v1.LogStorage) (*v1.LogStorage, error)
	Delete(ctx context.Context, id string) error

	ListKeys(ctx context.Context, logResourceId string, count *int, from *int) ([]v1.LogStorageAccessKey, error)
	CreateKey(ctx context.Context, logResourceId string, request *v1.LogStorageAccessKey) (*v1.LogStorageAccessKey, error)
	ReadKey(ctx context.Context, logResourceId string, id string) (*v1.LogStorageAccessKey, error)
	UpdateKey(ctx context.Context, logResourceId string, id string, request *v1.LogStorageAccessKey) (*v1.LogStorageAccessKey, error)
	DeleteKey(ctx context.Context, logResourceId string, id string) error
}

var _ LogsStorageAPI = (*logsStorageOp)(nil)

type logsStorageOp struct {
	client *v1.Client
}

func NewLogsStorageOp(client *v1.Client) LogsStorageAPI {
	return &logsStorageOp{client: client}
}

func (op *logsStorageOp) List(ctx context.Context, params v1.LogsStoragesListParams) ([]v1.LogStorage, error) {
	result, err := op.client.LogsStoragesList(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("LogsStorage.List", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		default:
			return nil, NewAPIError("LogsStorage.List", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("LogsStorage.List", 0, err)
	} else {
		return result.GetResults(), nil
	}
}

func (op *logsStorageOp) Read(ctx context.Context, resourceID string) (*v1.LogStorage, error) {
	id, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return nil, NewAPIError("LogsStorage.Read", 0, err)
	}
	params := v1.LogsStoragesRetrieveParams{ResourceID: id}
	result, err := op.client.LogsStoragesRetrieve(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("LogsStorage.Read", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("LogsStorage.Read", e.StatusCode, errors.Wrap(err, "log table not found"))
		default:
			return nil, NewAPIError("LogsStorage.Read", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("LogsStorage.Read", 0, err)
	} else {
		ret := new(v1.LogStorage)
		return Unwrap(ret, result)
	}
}

func (op *logsStorageOp) Create(ctx context.Context, params v1.LogStorageCreate) (*v1.LogStorage, error) {
	result, err := op.client.LogsStoragesCreate(ctx, &params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("LogsStorage.Create", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return nil, NewAPIError("LogsStorage.Create", e.StatusCode, errors.Wrap(err, "invalid parameter, or no space left for a new storage"))
		default:
			return nil, NewAPIError("LogsStorage.Create", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("LogsStorage.Create", 0, err)
	} else {
		return result, nil
	}
}

func (op *logsStorageOp) Update(ctx context.Context, id string, resource *v1.LogStorage) (*v1.LogStorage, error) {
	rid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, NewAPIError("LogsStorage.Update", 0, err)
	}
	params := v1.LogsStoragesUpdateParams{ResourceID: rid}
	body := v1.NewOptLogStorage(*resource)
	result, err := op.client.LogsStoragesUpdate(ctx, body, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("LogsStorage.Update", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return nil, NewAPIError("LogsStorage.Update", e.StatusCode, errors.Wrap(err, "invalid parameter"))
		default:
			return nil, NewAPIError("LogsStorage.Update", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("LogsStorage.Update", 0, err)
	} else {
		ret := new(v1.LogStorage)
		return Unwrap(ret, result)
	}
}

func (op *logsStorageOp) Delete(ctx context.Context, id string) error {
	rid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return NewAPIError("LogsStorage.Delete", 0, err)
	}
	params := v1.LogsStoragesDestroyParams{ResourceID: rid}
	err = op.client.LogsStoragesDestroy(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return NewAPIError("LogsStorage.Delete", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return NewAPIError("LogsStorage.Delete", e.StatusCode, errors.Wrap(err, "the request resource is not eligible for deletion"))
		default:
			return NewAPIError("LogsStorage.Delete", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return NewAPIError("LogsStorage.Delete", 0, err)
	}
	return nil
}

func (op *logsStorageOp) ListKeys(ctx context.Context, logResourceId string, count *int, from *int) ([]v1.LogStorageAccessKey, error) {
	rid, err := strconv.ParseInt(logResourceId, 10, 64)
	if err != nil {
		return nil, NewAPIError("LogsStorage.ListKeys", 0, err)
	}
	params := v1.LogsStoragesKeysListParams{
		Count:         intoOpt[v1.OptInt](count),
		From:          intoOpt[v1.OptInt](from),
		LogResourceID: rid,
	}
	result, err := op.client.LogsStoragesKeysList(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("LogsStorage.ListKeys", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		default:
			return nil, NewAPIError("LogsStorage.ListKeys", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("LogsStorage.ListKeys", 0, err)
	} else {
		return result.GetResults(), nil
	}
}

func (op *logsStorageOp) CreateKey(ctx context.Context, logResourceId string, request *v1.LogStorageAccessKey) (*v1.LogStorageAccessKey, error) {
	rid, err := strconv.ParseInt(logResourceId, 10, 64)
	if err != nil {
		return nil, NewAPIError("LogsStorage.CreateKey", 0, err)
	}
	params := v1.LogsStoragesKeysCreateParams{LogResourceID: rid}
	opt := v1.NewOptLogStorageAccessKey(*request)
	result, err := op.client.LogsStoragesKeysCreate(ctx, opt, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("LogsStorage.CreateKey", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("LogsStorage.CreateKey", e.StatusCode, errors.Wrap(err, "requested log storage not found"))
		case http.StatusBadRequest:
			return nil, NewAPIError("LogsStorage.CreateKey", e.StatusCode, errors.Wrap(err, "this log storage cannot have a key"))
		default:
			return nil, NewAPIError("LogsStorage.CreateKey", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("LogsStorage.CreateKey", 0, err)
	} else {
		ret := new(v1.LogStorageAccessKey)
		return Unwrap(ret, result)
	}
}

func (op *logsStorageOp) ReadKey(ctx context.Context, logResourceId string, id string) (*v1.LogStorageAccessKey, error) {
	rid, err := strconv.ParseInt(logResourceId, 10, 64)
	if err != nil {
		return nil, NewAPIError("LogsStorage.ReadKey", 0, err)
	}
	kid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, NewAPIError("LogsStorage.ReadKey", 0, err)
	}
	params := v1.LogsStoragesKeysRetrieveParams{LogResourceID: rid, ID: kid}
	result, err := op.client.LogsStoragesKeysRetrieve(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("LogsStorage.ReadKey", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("LogsStorage.ReadKey", e.StatusCode, errors.Wrap(err, "access key not found"))
		default:
			return nil, NewAPIError("LogsStorage.ReadKey", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("LogsStorage.ReadKey", 0, err)
	} else {
		ret := new(v1.LogStorageAccessKey)
		return Unwrap(ret, result)
	}
}

func (op *logsStorageOp) UpdateKey(ctx context.Context, logResourceId string, id string, request *v1.LogStorageAccessKey) (*v1.LogStorageAccessKey, error) {
	rid, err := strconv.ParseInt(logResourceId, 10, 64)
	if err != nil {
		return nil, NewAPIError("LogsStorage.UpdateKey", 0, err)
	}
	kid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, NewAPIError("LogsStorage.UpdateKey", 0, err)
	}
	params := v1.LogsStoragesKeysUpdateParams{LogResourceID: rid, ID: kid}
	opt := v1.NewOptLogStorageAccessKey(*request)
	result, err := op.client.LogsStoragesKeysUpdate(ctx, opt, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("LogsStorage.UpdateKey", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return nil, NewAPIError("LogsStorage.UpdateKey", e.StatusCode, errors.Wrap(err, "invalid parameter"))
		default:
			return nil, NewAPIError("LogsStorage.UpdateKey", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("LogsStorage.UpdateKey", 0, err)
	} else {
		key := new(v1.LogStorageAccessKey)
		return Unwrap(key, result)
	}
}

func (op *logsStorageOp) DeleteKey(ctx context.Context, logResourceId string, id string) error {
	rid, err := strconv.ParseInt(logResourceId, 10, 64)
	if err != nil {
		return NewAPIError("LogsStorage.DeleteKey", 0, err)
	}
	kid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return NewAPIError("LogsStorage.DeleteKey", 0, err)
	}
	params := v1.LogsStoragesKeysDestroyParams{LogResourceID: rid, ID: kid}
	err = op.client.LogsStoragesKeysDestroy(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return NewAPIError("LogsStorage.DeleteKey", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return NewAPIError("LogsStorage.DeleteKey", e.StatusCode, errors.Wrap(err, "the request resource is not eligible for deletion"))
		default:
			return NewAPIError("LogsStorage.DeleteKey", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return NewAPIError("LogsStorage.DeleteKey", 0, err)
	}
	return nil
}

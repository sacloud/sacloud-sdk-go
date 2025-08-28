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

type LogsStorageAPI interface {
	List(ctx context.Context, params v1.LogsStoragesListParams) ([]v1.LogTable, error)
	Create(ctx context.Context, params v1.LogTableCreate) (*v1.LogTable, error)
	Read(ctx context.Context, id int64) (*v1.LogTable, error)
	Update(ctx context.Context, id int64, request *v1.LogTable) (*v1.LogTable, error)
	Delete(ctx context.Context, id int64) error

	ListKeys(ctx context.Context, logResourceId int64, count int, from int) ([]v1.LogTableAccessKey, error)
	CreateKey(ctx context.Context, logResourceId int64, request *v1.LogTableAccessKey) (*v1.LogTableAccessKey, error)
	ReadKey(ctx context.Context, logResourceId int64, id int64) (*v1.LogTableAccessKey, error)
	UpdateKey(ctx context.Context, logResourceId int64, id int64, request *v1.LogTableAccessKey) (*v1.LogTableAccessKey, error)
	DeleteKey(ctx context.Context, logResourceId int64, id int64) error
}

var _ LogsStorageAPI = (*logsStorageOp)(nil)

type logsStorageOp struct {
	client *v1.Client
}

func NewLogsStorageOp(client *v1.Client) LogsStorageAPI {
	return &logsStorageOp{client: client}
}

func convertLogIcon(wrapped v1.WrappedLogTableIcon) v1.NilLogTableIcon {
	var icon v1.LogTableIcon

	icon.SetID(wrapped.GetID())
	return v1.NewNilLogTableIcon(icon)
}

func convertLogTableEndpointsIngester(wrapped v1.WrappedLogTableEndpointsIngester) v1.LogTableEndpointsIngester {
	var ret v1.LogTableEndpointsIngester

	ret.SetAddress(wrapped.GetAddress())
	ret.SetInsecure(wrapped.GetInsecure())
	return ret
}

func convertLogEndpoints(wrapped v1.WrappedLogTableEndpoints) v1.LogTableEndpoints {
	var ret v1.LogTableEndpoints

	wi := wrapped.GetIngester()
	raw := convertLogTableEndpointsIngester(wi)
	ret.SetIngester(raw)
	return ret
}

func convertLogUsage(wrapped v1.WrappedLogTableUsage) v1.LogTableUsage {
	var ret v1.LogTableUsage

	ret.SetLogRoutings(wrapped.GetLogRoutings())
	ret.SetLogRecordingRules(wrapped.GetLogRecordingRules())
	return ret
}

func convertLogTable(wrapped v1.WrappedLogTable) v1.LogTable {
	var ret v1.LogTable
	ret.SetID(v1.NewNilInt64(wrapped.GetID()))
	ret.SetName(wrapped.GetName())
	ret.SetDescription(wrapped.GetDescription())
	ret.SetTags(wrapped.GetTags())
	if wt, ok := wrapped.GetIcon().Get(); ok {
		ret.SetIcon(convertLogIcon(wt))
	} else {
		var icon v1.NilLogTableIcon
		icon.SetToNull()
		ret.SetIcon(icon)
	}
	ret.SetExpireDay(wrapped.GetExpireDay())
	ret.SetCreatedAt(wrapped.GetCreatedAt())
	ret.SetEndpoints(convertLogEndpoints(wrapped.GetEndpoints()))
	ret.SetAccountID(wrapped.GetAccountID())
	ret.SetResourceID(wrapped.GetResourceID())
	ret.SetIsSystem(wrapped.GetIsSystem())
	ret.SetUsage(convertLogUsage(wrapped.GetUsage()))
	return ret
}

func convertLogKey(wrapped v1.WrappedLogTableAccessKey) v1.LogTableAccessKey {
	var ret v1.LogTableAccessKey
	ret.SetID(wrapped.GetID())
	ret.SetSecret(wrapped.GetSecret())
	ret.SetDescription(wrapped.GetDescription())
	return ret
}

func (op *logsStorageOp) List(ctx context.Context, params v1.LogsStoragesListParams) ([]v1.LogTable, error) {
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

func (op *logsStorageOp) Read(ctx context.Context, resourceID int64) (*v1.LogTable, error) {
	params := v1.LogsStoragesRetrieveParams{ResourceID: resourceID}
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
		ret := convertLogTable(*result)
		return &ret, nil
	}
}

func (op *logsStorageOp) Create(ctx context.Context, params v1.LogTableCreate) (*v1.LogTable, error) {
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

func (op *logsStorageOp) Update(ctx context.Context, id int64, resource *v1.LogTable) (*v1.LogTable, error) {
	params := v1.LogsStoragesUpdateParams{ResourceID: id}
	body := v1.NewOptLogTable(*resource)
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
		ret := convertLogTable(*result)
		return &ret, nil
	}
}

func (op *logsStorageOp) Delete(ctx context.Context, id int64) error {
	params := v1.LogsStoragesDestroyParams{ResourceID: id}
	err := op.client.LogsStoragesDestroy(ctx, params)
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

func (op *logsStorageOp) ListKeys(ctx context.Context, logResourceId int64, count int, from int) ([]v1.LogTableAccessKey, error) {
	params := v1.LogsStoragesKeysListParams{
		Count:         v1.NewOptInt(count),
		From:          v1.NewOptInt(from),
		LogResourceID: logResourceId,
	}
	params.Count.SetTo(count)
	params.From.SetTo(from)
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

func (op *logsStorageOp) CreateKey(ctx context.Context, logResourceId int64, request *v1.LogTableAccessKey) (*v1.LogTableAccessKey, error) {
	params := v1.LogsStoragesKeysCreateParams{LogResourceID: logResourceId}
	opt := v1.NewOptLogTableAccessKey(*request)
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
		key := convertLogKey(*result)
		return &key, nil
	}
}

func (op *logsStorageOp) ReadKey(ctx context.Context, logResourceId int64, id int64) (*v1.LogTableAccessKey, error) {
	params := v1.LogsStoragesKeysRetrieveParams{LogResourceID: logResourceId, ID: id}
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
		key := convertLogKey(*result)
		return &key, nil
	}
}

func (op *logsStorageOp) UpdateKey(ctx context.Context, logResourceId int64, id int64, request *v1.LogTableAccessKey) (*v1.LogTableAccessKey, error) {
	params := v1.LogsStoragesKeysUpdateParams{LogResourceID: logResourceId, ID: id}
	opt := v1.NewOptLogTableAccessKey(*request)
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
		key := convertLogKey(*result)
		return &key, nil
	}
}

func (op *logsStorageOp) DeleteKey(ctx context.Context, logResourceId int64, id int64) error {
	params := v1.LogsStoragesKeysDestroyParams{LogResourceID: logResourceId, ID: id}
	err := op.client.LogsStoragesKeysDestroy(ctx, params)
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

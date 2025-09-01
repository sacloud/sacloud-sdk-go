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

type MetricsStorageAPI interface {
	List(ctx context.Context, count int, from int) ([]v1.MetricsTank, error)
	Create(ctx context.Context, request v1.MetricsTankCreate) (*v1.MetricsTank, error)
	Read(ctx context.Context, id int64) (*v1.MetricsTank, error)
	Update(ctx context.Context, id int64, request *v1.MetricsTank) (*v1.MetricsTank, error)
	Delete(ctx context.Context, id int64) error

	ListKeys(ctx context.Context, metricsResourceId int64, count int, from int) ([]v1.MetricsTankAccessKey, error)
	CreateKey(ctx context.Context, metricsResourceId int64, request *v1.MetricsTankAccessKey) (*v1.MetricsTankAccessKey, error)
	ReadKey(ctx context.Context, metricsResourceId int64, id int64) (*v1.MetricsTankAccessKey, error)
	UpdateKey(ctx context.Context, metricsResourceId int64, id int64, request *v1.MetricsTankAccessKey) (*v1.MetricsTankAccessKey, error)
	DeleteKey(ctx context.Context, metricsResourceId int64, id int64) error
}

var _ MetricsStorageAPI = (*metricsStorageOp)(nil)

type metricsStorageOp struct {
	client *v1.Client
}

func NewMetricsStorageOp(client *v1.Client) MetricsStorageAPI {
	return &metricsStorageOp{client: client}
}

func convertIcon(wrapped v1.WrappedMetricsTankIcon) v1.NilMetricsTankIcon {
	var icon v1.MetricsTankIcon

	icon.SetID(wrapped.GetID())
	return v1.NewNilMetricsTankIcon(icon)
}

func convertEndpoints(wrapped v1.WrappedMetricsTankEndpoints) v1.MetricsTankEndpoints {
	var endpoints v1.MetricsTankEndpoints

	endpoints.SetAddress(wrapped.GetAddress())
	return endpoints
}

func convertUsage(wrapped v1.WrappedMetricsTankUsage) v1.MetricsTankUsage {
	var usage v1.MetricsTankUsage

	usage.SetMetricsRoutings(wrapped.GetMetricsRoutings())
	usage.SetAlertRules(wrapped.GetAlertRules())
	usage.SetLogRecordingRules(wrapped.GetLogRecordingRules())
	return usage
}

func convertTank(result v1.WrappedMetricsTank) v1.MetricsTank {
	var tank v1.MetricsTank

	tank.SetID(v1.NewNilInt64(result.GetID()))
	tank.SetName(result.GetName())
	tank.SetDescription(result.GetDescription())
	tank.SetTags(result.GetTags())
	if wt, ok := result.GetIcon().Get(); ok {
		tank.SetIcon(convertIcon(wt))
	}
	tank.SetIsSystem(result.GetIsSystem())
	tank.SetAccountID(result.GetAccountID())
	tank.SetResourceID(result.GetResourceID())
	tank.SetEndpoints(convertEndpoints(result.GetEndpoints()))
	tank.SetCreatedAt(result.GetCreatedAt())
	tank.SetUpdatedAt(result.GetUpdatedAt())
	tank.SetUsage(convertUsage(result.GetUsage()))
	return tank
}

func convertKey(result v1.WrappedMetricsTankAccessKey) v1.MetricsTankAccessKey {
	var key v1.MetricsTankAccessKey

	key.SetID(result.GetID())
	key.SetSecret(result.GetSecret())
	key.SetDescription(result.GetDescription())
	return key
}

func (op *metricsStorageOp) List(ctx context.Context, count int, from int) ([]v1.MetricsTank, error) {
	params := v1.MetricsStoragesListParams{}
	params.Count.SetTo(count)
	params.From.SetTo(from)
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

func (op *metricsStorageOp) Read(ctx context.Context, resourceID int64) (*v1.MetricsTank, error) {
	params := v1.MetricsStoragesRetrieveParams{ResourceID: resourceID}
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
		ret := convertTank(*result)
		return &ret, nil
	}
}

func (op *metricsStorageOp) Create(ctx context.Context, body v1.MetricsTankCreate) (*v1.MetricsTank, error) {
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

func (op *metricsStorageOp) Update(ctx context.Context, id int64, resource *v1.MetricsTank) (*v1.MetricsTank, error) {
	query := v1.MetricsStoragesUpdateParams{ResourceID: id}
	body := v1.NewOptMetricsTank(*resource)
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
		ret := convertTank(*result)
		return &ret, nil
	}
}

func (op *metricsStorageOp) Delete(ctx context.Context, resourceID int64) error {
	params := v1.MetricsStoragesDestroyParams{ResourceID: resourceID}
	err := op.client.MetricsStoragesDestroy(ctx, params)
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

func (op *metricsStorageOp) ListKeys(ctx context.Context, metricsResourceId int64, count int, from int) ([]v1.MetricsTankAccessKey, error) {
	params := v1.MetricsStoragesKeysListParams{MetricsResourceID: metricsResourceId}
	params.Count.SetTo(count)
	params.From.SetTo(from)
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

func (op *metricsStorageOp) CreateKey(ctx context.Context, metricsResourceId int64, request *v1.MetricsTankAccessKey) (*v1.MetricsTankAccessKey, error) {
	params := v1.MetricsStoragesKeysCreateParams{MetricsResourceID: metricsResourceId}
	opt := v1.NewOptMetricsTankAccessKey(*request)
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
		key := convertKey(*result)
		return &key, nil
	}
}

func (op *metricsStorageOp) ReadKey(ctx context.Context, metricsResourceId int64, id int64) (*v1.MetricsTankAccessKey, error) {
	params := v1.MetricsStoragesKeysRetrieveParams{MetricsResourceID: metricsResourceId, ID: id}
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
		key := convertKey(*result)
		return &key, nil
	}
}

func (op *metricsStorageOp) UpdateKey(ctx context.Context, metricsResourceId int64, id int64, request *v1.MetricsTankAccessKey) (*v1.MetricsTankAccessKey, error) {
	params := v1.MetricsStoragesKeysUpdateParams{MetricsResourceID: metricsResourceId, ID: id}
	opt := v1.NewOptMetricsTankAccessKey(*request)
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
		key := convertKey(*result)
		return &key, nil
	}
}

// DeleteKey deletes an access key for a metrics storage resource.
func (op *metricsStorageOp) DeleteKey(ctx context.Context, metricsResourceId int64, id int64) error {
	params := v1.MetricsStoragesKeysDestroyParams{MetricsResourceID: metricsResourceId, ID: id}
	err := op.client.MetricsStoragesKeysDestroy(ctx, params)
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

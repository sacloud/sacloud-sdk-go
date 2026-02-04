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
	"time"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	ogen "github.com/ogen-go/ogen/validate"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
)

type MetricsStorageAPI interface {
	List(ctx context.Context, params MetricsStorageListParams) ([]v1.MetricsStorage, error)
	Create(ctx context.Context, request MetricsStorageCreateParams) (*v1.MetricsStorage, error)
	Read(ctx context.Context, id string) (*v1.MetricsStorage, error)
	Update(ctx context.Context, id string, request MetricsStorageUpdateParams) (*v1.MetricsStorage, error)
	Delete(ctx context.Context, id string) error

	ReadDailyStats(ctx context.Context, resourceID string, startDate, endDate *time.Time) ([]v1.MetricsStorageDailyUsage, error)
	ReadMonthlyStats(ctx context.Context, resourceID string, year int) ([]v1.MetricsStorageMonthlyUsage, error)

	ListKeys(ctx context.Context, metricsResourceId string, count *int, from *int) ([]v1.MetricsStorageAccessKey, error)
	CreateKey(ctx context.Context, metricsResourceId string, description *string) (*v1.MetricsStorageAccessKey, error)
	ReadKey(ctx context.Context, metricsResourceId string, id uuid.UUID) (*v1.MetricsStorageAccessKey, error)
	UpdateKey(ctx context.Context, metricsResourceId string, id uuid.UUID, description *string) (*v1.MetricsStorageAccessKey, error)
	DeleteKey(ctx context.Context, metricsResourceId string, id uuid.UUID) error
}

var _ MetricsStorageAPI = (*metricsStorageOp)(nil)

type metricsStorageOp struct {
	client *v1.Client
}

func NewMetricsStorageOp(client *v1.Client) MetricsStorageAPI {
	return &metricsStorageOp{client: client}
}

type MetricsStorageListParams struct {
	Count      *int
	From       *int
	AccountID  *string
	ResourceID *string
	IsSystem   *bool
}

func (op *metricsStorageOp) List(ctx context.Context, params MetricsStorageListParams) ([]v1.MetricsStorage, error) {
	resourceId, err := fromStringPtr[v1.OptInt64, int64](params.ResourceID)
	if err != nil {
		return nil, NewError("MetricsStorage.List", err)
	}
	result, err := op.client.MetricsStoragesList(ctx, v1.MetricsStoragesListParams{
		Count:      intoOpt[v1.OptInt](params.Count),
		From:       intoOpt[v1.OptInt](params.From),
		AccountID:  intoOpt[v1.OptString](params.AccountID),
		ResourceID: resourceId,
		IsSystem:   intoOpt[v1.OptBool](params.IsSystem),
	})
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
		return nil, NewError("MetricsStorage.Read", err)
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

type MetricsStorageCreateParams struct {
	Name        string
	Description *string
	IsSystem    bool
}

func (op *metricsStorageOp) Create(ctx context.Context, params MetricsStorageCreateParams) (*v1.MetricsStorage, error) {
	body := v1.MetricsStorageCreate{
		Name:        params.Name,
		Description: intoOpt[v1.OptString](params.Description),
		IsSystem:    params.IsSystem,
	}
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

type MetricsStorageUpdateParams struct {
	Name        *string
	Description *string
}

func (op *metricsStorageOp) Update(ctx context.Context, id string, params MetricsStorageUpdateParams) (*v1.MetricsStorage, error) {
	rid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, NewError("MetricsStorage.Update", err)
	}
	query := v1.MetricsStoragesPartialUpdateParams{ResourceID: rid}
	body := v1.NewOptPatchedMetricsStorage(v1.PatchedMetricsStorage{
		Name:        intoOpt[v1.OptString](params.Name),
		Description: intoOpt[v1.OptString](params.Description),
	})
	result, err := op.client.MetricsStoragesPartialUpdate(ctx, body, query)
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
		return NewError("MetricsStorage.Delete", err)
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

func (op *metricsStorageOp) ReadDailyStats(ctx context.Context, resourceID string, startDate, endDate *time.Time) ([]v1.MetricsStorageDailyUsage, error) {
	rid, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return nil, NewError("MetricsStorage.ReadDailyStats", err)
	}
	query := v1.MetricsStoragesStatsDailyRetrieveParams{
		ResourceID: rid,
		StartDate:  intoOpt[v1.OptDate](startDate),
		EndDate:    intoOpt[v1.OptDate](endDate),
	}
	result, err := op.client.MetricsStoragesStatsDailyRetrieve(ctx, query)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("MetricsStorage.ReadDailyStats", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("MetricsStorage.ReadDailyStats", e.StatusCode, errors.Wrap(err, "metrics storage not found"))
		case http.StatusBadRequest:
			return nil, NewAPIError("MetricsStorage.ReadDailyStats", e.StatusCode, errors.Wrap(err, "invalid parameter"))
		default:
			return nil, NewAPIError("MetricsStorage.ReadDailyStats", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("MetricsStorage.ReadDailyStats", 0, err)
	} else {
		return result.GetUsages(), nil
	}
}

func (op *metricsStorageOp) ReadMonthlyStats(ctx context.Context, resourceID string, year int) ([]v1.MetricsStorageMonthlyUsage, error) {
	rid, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return nil, NewError("MetricsStorage.ReadMonthlyStats", err)
	}
	query := v1.MetricsStoragesStatsMonthlyRetrieveParams{
		ResourceID: rid,
		Year:       year,
	}
	result, err := op.client.MetricsStoragesStatsMonthlyRetrieve(ctx, query)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("MetricsStorage.ReadMonthlyStats", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("MetricsStorage.ReadMonthlyStats", e.StatusCode, errors.Wrap(err, "metrics storage not found"))
		case http.StatusBadRequest:
			return nil, NewAPIError("MetricsStorage.ReadMonthlyStats", e.StatusCode, errors.Wrap(err, "invalid parameter, year must be between 1970 and 2100"))
		default:
			return nil, NewAPIError("MetricsStorage.ReadMonthlyStats", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("MetricsStorage.ReadMonthlyStats", 0, err)
	} else {
		return result.GetUsages(), nil
	}
}

func (op *metricsStorageOp) ListKeys(ctx context.Context, metricsResourceId string, count *int, from *int) ([]v1.MetricsStorageAccessKey, error) {
	rid, err := strconv.ParseInt(metricsResourceId, 10, 64)
	if err != nil {
		return nil, NewError("MetricsStorage.ListKeys", err)
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

func (op *metricsStorageOp) CreateKey(ctx context.Context, metricsResourceId string, description *string) (*v1.MetricsStorageAccessKey, error) {
	rid, err := strconv.ParseInt(metricsResourceId, 10, 64)
	if err != nil {
		return nil, NewError("MetricsStorage.CreateKey", err)
	}
	params := v1.MetricsStoragesKeysCreateParams{MetricsResourceID: rid}
	opt := v1.NewOptMetricsStorageAccessKey(v1.MetricsStorageAccessKey{
		Description: intoOpt[v1.OptString](description),
	})
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

func (op *metricsStorageOp) ReadKey(ctx context.Context, metricsResourceId string, id uuid.UUID) (*v1.MetricsStorageAccessKey, error) {
	rid, err := strconv.ParseInt(metricsResourceId, 10, 64)
	if err != nil {
		return nil, NewError("MetricsStorage.ReadKey", err)
	}
	params := v1.MetricsStoragesKeysRetrieveParams{MetricsResourceID: rid, UID: id}
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

func (op *metricsStorageOp) UpdateKey(ctx context.Context, metricsResourceId string, id uuid.UUID, description *string) (*v1.MetricsStorageAccessKey, error) {
	rid, err := strconv.ParseInt(metricsResourceId, 10, 64)
	if err != nil {
		return nil, NewError("MetricsStorage.UpdateKey", err)
	}
	params := v1.MetricsStoragesKeysUpdateParams{MetricsResourceID: rid, UID: id}
	opt := v1.NewOptMetricsStorageAccessKey(v1.MetricsStorageAccessKey{
		Description: intoOpt[v1.OptString](description),
	})
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
func (op *metricsStorageOp) DeleteKey(ctx context.Context, metricsResourceId string, id uuid.UUID) error {
	rid, err := strconv.ParseInt(metricsResourceId, 10, 64)
	if err != nil {
		return NewError("MetricsStorage.DeleteKey", err)
	}
	params := v1.MetricsStoragesKeysDestroyParams{MetricsResourceID: rid, UID: id}
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

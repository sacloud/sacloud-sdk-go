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

type LogsStorageAPI interface {
	List(ctx context.Context, params LogsStoragesListParams) ([]v1.LogStorage, error)
	Create(ctx context.Context, params LogStorageCreateParams) (*v1.LogStorage, error)
	Read(ctx context.Context, id string) (*v1.LogStorage, error)
	Update(ctx context.Context, id string, params LogStorageUpdateParams) (*v1.LogStorage, error)
	Delete(ctx context.Context, id string) error

	SetExpire(ctx context.Context, resourceID string, days int) (*v1.LogStorage, error)
	StatsDaily(ctx context.Context, resourceID string, startDate, endDate *time.Time) ([]v1.LogStorageDailyUsage, error)
	StatsMonthly(ctx context.Context, resourceID string, year int) ([]v1.LogStorageMonthlyUsage, error)

	ListKeys(ctx context.Context, logResourceId string, count *int, from *int) ([]v1.LogStorageAccessKey, error)
	CreateKey(ctx context.Context, logResourceId string, description *string) (*v1.LogStorageAccessKey, error)
	ReadKey(ctx context.Context, logResourceId string, id uuid.UUID) (*v1.LogStorageAccessKey, error)
	UpdateKey(ctx context.Context, logResourceId string, id uuid.UUID, description *string) (*v1.LogStorageAccessKey, error)
	DeleteKey(ctx context.Context, logResourceId string, id uuid.UUID) error
}

var _ LogsStorageAPI = (*logsStorageOp)(nil)

type logsStorageOp struct {
	client *v1.Client
}

func NewLogsStorageOp(client *v1.Client) LogsStorageAPI {
	return &logsStorageOp{client: client}
}

type LogsStoragesListParams struct {
	AccountID            *string
	BucketClassification *v1.LogsStoragesListBucketClassification
	Count                *int
	From                 *int
	IsSystem             *bool
	ResourceID           *string
	Status               *v1.LogsStoragesListStatus
}

func (op *logsStorageOp) List(ctx context.Context, p LogsStoragesListParams) ([]v1.LogStorage, error) {
	id, err := fromStringPtr[v1.OptInt64, int64](p.ResourceID)
	if err != nil {
		return nil, NewError("LogsStorage.List", err)
	}
	params := v1.LogsStoragesListParams{
		AccountID:            intoOpt[v1.OptString](p.AccountID),
		BucketClassification: intoOpt[v1.OptLogsStoragesListBucketClassification](p.BucketClassification),
		Count:                intoOpt[v1.OptInt](p.Count),
		From:                 intoOpt[v1.OptInt](p.From),
		IsSystem:             intoOpt[v1.OptBool](p.IsSystem),
		ResourceID:           id,
		Status:               intoOpt[v1.OptLogsStoragesListStatus](p.Status),
	}
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
		return nil, NewError("LogsStorage.Read", err)
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

type LogStorageCreateParams struct {
	Name           string
	Description    *string
	IsSystem       bool
	Classification *v1.LogStorageCreateClassification
}

func (op *logsStorageOp) Create(ctx context.Context, params LogStorageCreateParams) (*v1.LogStorage, error) {
	req := v1.LogStorageCreate{
		Classification: intoOpt[v1.OptLogStorageCreateClassification](params.Classification),
		IsSystem:       params.IsSystem,
		Name:           params.Name,
		Description:    intoOpt[v1.OptString](params.Description),
	}
	result, err := op.client.LogsStoragesCreate(ctx, &req)
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

type LogStorageUpdateParams struct {
	Name        *string
	Description *string
}

func (op *logsStorageOp) Update(ctx context.Context, id string, p LogStorageUpdateParams) (*v1.LogStorage, error) {
	resource := v1.PatchedLogStorage{
		Name:        intoOpt[v1.OptString](p.Name),
		Description: intoOpt[v1.OptString](p.Description),
	}
	rid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, NewError("LogsStorage.Update", err)
	}
	params := v1.LogsStoragesPartialUpdateParams{ResourceID: rid}
	body := v1.NewOptPatchedLogStorage(resource)
	result, err := op.client.LogsStoragesPartialUpdate(ctx, body, params)
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
		return NewError("LogsStorage.Delete", err)
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

func (op *logsStorageOp) SetExpire(ctx context.Context, resourceID string, days int) (*v1.LogStorage, error) {
	rid, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return nil, NewError("LogsStorage.SetExpire", err)
	}
	params := v1.LogsStoragesSetExpireCreateParams{ResourceID: rid}
	body := v1.SetLogStorageExpireDay{Days: days}
	result, err := op.client.LogsStoragesSetExpireCreate(ctx, &body, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("LogsStorage.SetExpire", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("LogsStorage.SetExpire", e.StatusCode, errors.Wrap(err, "log storage not found"))
		case http.StatusBadRequest:
			return nil, NewAPIError("LogsStorage.SetExpire", e.StatusCode, errors.Wrap(err, "invalid parameter, days must be between 1 and 730"))
		default:
			return nil, NewAPIError("LogsStorage.SetExpire", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("LogsStorage.SetExpire", 0, err)
	} else {
		return result, nil
	}
}

func (op *logsStorageOp) StatsDaily(ctx context.Context, resourceID string, startDate, endDate *time.Time) ([]v1.LogStorageDailyUsage, error) {
	rid, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return nil, NewError("LogsStorage.StatsDaily", err)
	}
	query := v1.LogsStoragesStatsDailyRetrieveParams{
		ResourceID: rid,
		StartDate:  intoOpt[v1.OptDate](startDate),
		EndDate:    intoOpt[v1.OptDate](endDate),
	}
	result, err := op.client.LogsStoragesStatsDailyRetrieve(ctx, query)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("LogsStorage.StatsDaily", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("LogsStorage.StatsDaily", e.StatusCode, errors.Wrap(err, "log storage not found"))
		case http.StatusBadRequest:
			return nil, NewAPIError("LogsStorage.StatsDaily", e.StatusCode, errors.Wrap(err, "invalid parameter"))
		default:
			return nil, NewAPIError("LogsStorage.StatsDaily", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("LogsStorage.StatsDaily", 0, err)
	} else {
		return result.GetUsages(), nil
	}
}

func (op *logsStorageOp) StatsMonthly(ctx context.Context, resourceID string, year int) ([]v1.LogStorageMonthlyUsage, error) {
	rid, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return nil, NewError("LogsStorage.StatsMonthly", err)
	}
	query := v1.LogsStoragesStatsMonthlyRetrieveParams{
		ResourceID: rid,
		Year:       year,
	}
	result, err := op.client.LogsStoragesStatsMonthlyRetrieve(ctx, query)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("LogsStorage.StatsMonthly", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("LogsStorage.StatsMonthly", e.StatusCode, errors.Wrap(err, "log storage not found"))
		case http.StatusBadRequest:
			return nil, NewAPIError("LogsStorage.StatsMonthly", e.StatusCode, errors.Wrap(err, "invalid parameter, year must be between 1970 and 2100"))
		default:
			return nil, NewAPIError("LogsStorage.StatsMonthly", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("LogsStorage.StatsMonthly", 0, err)
	} else {
		return result.GetUsages(), nil
	}
}

func (op *logsStorageOp) ListKeys(ctx context.Context, logResourceId string, count *int, from *int) ([]v1.LogStorageAccessKey, error) {
	rid, err := strconv.ParseInt(logResourceId, 10, 64)
	if err != nil {
		return nil, NewError("LogsStorage.ListKeys", err)
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

func (op *logsStorageOp) CreateKey(ctx context.Context, logResourceId string, description *string) (*v1.LogStorageAccessKey, error) {
	request := v1.LogStorageAccessKey{
		Description: intoOpt[v1.OptString](description),
	}
	rid, err := strconv.ParseInt(logResourceId, 10, 64)
	if err != nil {
		return nil, NewError("LogsStorage.CreateKey", err)
	}
	params := v1.LogsStoragesKeysCreateParams{LogResourceID: rid}
	opt := v1.NewOptLogStorageAccessKey(request)
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

func (op *logsStorageOp) ReadKey(ctx context.Context, logResourceId string, id uuid.UUID) (*v1.LogStorageAccessKey, error) {
	rid, err := strconv.ParseInt(logResourceId, 10, 64)
	if err != nil {
		return nil, NewError("LogsStorage.ReadKey", err)
	}
	params := v1.LogsStoragesKeysRetrieveParams{LogResourceID: rid, UID: id}
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

func (op *logsStorageOp) UpdateKey(ctx context.Context, logResourceId string, id uuid.UUID, description *string) (*v1.LogStorageAccessKey, error) {
	// :FIXME: does it make sense to allow a nil description?
	request := v1.LogStorageAccessKey{
		Description: intoOpt[v1.OptString](description),
	}
	rid, err := strconv.ParseInt(logResourceId, 10, 64)
	if err != nil {
		return nil, NewError("LogsStorage.UpdateKey", err)
	}
	params := v1.LogsStoragesKeysUpdateParams{LogResourceID: rid, UID: id}
	opt := v1.NewOptLogStorageAccessKey(request)
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

func (op *logsStorageOp) DeleteKey(ctx context.Context, logResourceId string, id uuid.UUID) error {
	rid, err := strconv.ParseInt(logResourceId, 10, 64)
	if err != nil {
		return NewError("LogsStorage.DeleteKey", err)
	}
	params := v1.LogsStoragesKeysDestroyParams{LogResourceID: rid, UID: id}
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

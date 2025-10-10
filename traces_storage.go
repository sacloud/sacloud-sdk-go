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

type TracesStorageAPI interface {
	List(ctx context.Context, params TracesStorageListParams) ([]v1.TraceStorage, error)
	Create(ctx context.Context, request TracesStorageCreateParams) (*v1.TraceStorage, error)
	Read(ctx context.Context, id string) (*v1.TraceStorage, error)
	Update(ctx context.Context, id string, request TracesStorageUpdateParams) (*v1.TraceStorage, error)
	Delete(ctx context.Context, id string) error

	ListKeys(ctx context.Context, tracesResourceId string, count *int, from *int) ([]v1.TraceStorageAccessKey, error)
	CreateKey(ctx context.Context, tracesResourceId string, description *string) (*v1.TraceStorageAccessKey, error)
	ReadKey(ctx context.Context, tracesResourceId string, id uuid.UUID) (*v1.TraceStorageAccessKey, error)
	UpdateKey(ctx context.Context, tracesResourceId string, id uuid.UUID, description *string) (*v1.TraceStorageAccessKey, error)
	DeleteKey(ctx context.Context, tracesResourceId string, id uuid.UUID) error
}

var _ TracesStorageAPI = (*tracesStorageOp)(nil)

type tracesStorageOp struct {
	client *v1.Client
}

func NewTracesStorageOp(client *v1.Client) TracesStorageAPI {
	return &tracesStorageOp{client: client}
}

type TracesStorageListParams struct {
	Count                *int
	From                 *int
	AccountID            *string
	ResourceID           *string
	BucketClassification *v1.TracesStoragesListLogStorageBucketClassification
}

func (op *tracesStorageOp) List(ctx context.Context, params TracesStorageListParams) ([]v1.TraceStorage, error) {
	resourceId, err := fromStringPtr[v1.OptInt64, int64](params.ResourceID)
	if err != nil {
		return nil, NewError("TracesStorage.List", err)
	}
	result, err := op.client.TracesStoragesList(ctx, v1.TracesStoragesListParams{
		Count:                          intoOpt[v1.OptInt](params.Count),
		From:                           intoOpt[v1.OptInt](params.From),
		AccountID:                      intoOpt[v1.OptString](params.AccountID),
		ResourceID:                     resourceId,
		LogStorageBucketClassification: intoOpt[v1.OptTracesStoragesListLogStorageBucketClassification](params.BucketClassification),
	})
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("TracesStorage.List", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		default:
			return nil, NewAPIError("TracesStorage.List", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("TracesStorage.List", 0, err)
	} else {
		return result.Results, nil
	}
}

func (op *tracesStorageOp) Read(ctx context.Context, resourceID string) (*v1.TraceStorage, error) {
	id, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return nil, NewError("TracesStorage.Read", err)
	}
	params := v1.TracesStoragesRetrieveParams{ResourceID: id}
	result, err := op.client.TracesStoragesRetrieve(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("TracesStorage.Read", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("TracesStorage.Read", e.StatusCode, errors.Wrap(err, "resource not found"))
		default:
			return nil, NewAPIError("TracesStorage.Read", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("TracesStorage.Read", 0, err)
	} else {
		ret := new(v1.TraceStorage)
		return Unwrap(ret, result)
	}
}

type TracesStorageCreateParams struct {
	Name           string
	Description    *string
	Classification *v1.TraceStorageCreateClassification
}

func (op *tracesStorageOp) Create(ctx context.Context, params TracesStorageCreateParams) (*v1.TraceStorage, error) {
	body := v1.TraceStorageCreate{
		Name:           params.Name,
		Description:    intoOpt[v1.OptString](params.Description),
		Classification: intoOpt[v1.OptTraceStorageCreateClassification](params.Classification),
	}
	result, err := op.client.TracesStoragesCreate(ctx, &body)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("TracesStorage.Create", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return nil, NewAPIError("TracesStorage.Create", e.StatusCode, errors.Wrap(err, "invalid request"))
		default:
			return nil, NewAPIError("TracesStorage.Create", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("TracesStorage.Create", 0, err)
	}
	return result, nil
}

type TracesStorageUpdateParams struct {
	Name        *string
	Description *string
}

func (op *tracesStorageOp) Update(ctx context.Context, id string, params TracesStorageUpdateParams) (*v1.TraceStorage, error) {
	rid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, NewError("TracesStorage.Update", err)
	}
	query := v1.TracesStoragesPartialUpdateParams{ResourceID: rid}
	body := v1.NewOptPatchedTraceStorage(v1.PatchedTraceStorage{
		Name:        intoOpt[v1.OptString](params.Name),
		Description: intoOpt[v1.OptString](params.Description),
	})
	result, err := op.client.TracesStoragesPartialUpdate(ctx, body, query)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("TracesStorage.Update", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("TracesStorage.Update", e.StatusCode, errors.Wrap(err, "resource not found"))
		case http.StatusBadRequest:
			return nil, NewAPIError("TracesStorage.Update", e.StatusCode, errors.Wrap(err, "invalid parameter"))
		default:
			return nil, NewAPIError("TracesStorage.Update", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("TracesStorage.Update", 0, err)
	} else {
		ret := new(v1.TraceStorage)
		return Unwrap(ret, result)
	}
}

func (op *tracesStorageOp) Delete(ctx context.Context, resourceID string) error {
	rid, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return NewError("TracesStorage.Delete", err)
	}
	params := v1.TracesStoragesDestroyParams{ResourceID: rid}
	err = op.client.TracesStoragesDestroy(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return NewAPIError("TracesStorage.Delete", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return NewAPIError("TracesStorage.Delete", e.StatusCode, errors.Wrap(err, "the request resource is not eligible for deletion"))
		default:
			return NewAPIError("TracesStorage.Delete", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return NewAPIError("TracesStorage.Delete", 0, err)
	}
	return nil
}

func (op *tracesStorageOp) ListKeys(ctx context.Context, tracesResourceId string, count *int, from *int) ([]v1.TraceStorageAccessKey, error) {
	rid, err := strconv.ParseInt(tracesResourceId, 10, 64)
	if err != nil {
		return nil, NewError("TracesStorage.ListKeys", err)
	}
	params := v1.TracesStoragesKeysListParams{
		TraceResourceID: rid,
		Count:           intoOpt[v1.OptInt](count),
		From:            intoOpt[v1.OptInt](from),
	}
	result, err := op.client.TracesStoragesKeysList(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("TracesStorage.ListKeys", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("TracesStorage.ListKeys", e.StatusCode, errors.Wrap(err, "resource not found"))
		default:
			return nil, NewAPIError("TracesStorage.ListKeys", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("TracesStorage.ListKeys", 0, err)
	} else {
		return result.Results, nil
	}
}

func (op *tracesStorageOp) CreateKey(ctx context.Context, tracesResourceId string, description *string) (*v1.TraceStorageAccessKey, error) {
	rid, err := strconv.ParseInt(tracesResourceId, 10, 64)
	if err != nil {
		return nil, NewError("TracesStorage.CreateKey", err)
	}
	params := v1.TracesStoragesKeysCreateParams{TraceResourceID: rid}
	opt := v1.NewOptTraceStorageAccessKey(v1.TraceStorageAccessKey{
		Description: intoOpt[v1.OptString](description),
	})
	result, err := op.client.TracesStoragesKeysCreate(ctx, opt, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("TracesStorage.CreateKey", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("TracesStorage.CreateKey", e.StatusCode, errors.Wrap(err, "requested traces storage not found"))
		case http.StatusBadRequest:
			return nil, NewAPIError("TracesStorage.CreateKey", e.StatusCode, errors.Wrap(err, "this traces storage cannot have a key"))
		default:
			return nil, NewAPIError("TracesStorage.CreateKey", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("TracesStorage.CreateKey", 0, err)
	} else {
		ret := new(v1.TraceStorageAccessKey)
		return Unwrap(ret, result)
	}
}

func (op *tracesStorageOp) ReadKey(ctx context.Context, tracesResourceId string, id uuid.UUID) (*v1.TraceStorageAccessKey, error) {
	rid, err := strconv.ParseInt(tracesResourceId, 10, 64)
	if err != nil {
		return nil, NewError("TracesStorage.ReadKey", err)
	}
	params := v1.TracesStoragesKeysRetrieveParams{TraceResourceID: rid, UID: id}
	result, err := op.client.TracesStoragesKeysRetrieve(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("TracesStorage.ReadKey", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("TracesStorage.ReadKey", e.StatusCode, errors.Wrap(err, "access key not found"))
		default:
			return nil, NewAPIError("TracesStorage.ReadKey", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("TracesStorage.ReadKey", 0, err)
	} else {
		ret := new(v1.TraceStorageAccessKey)
		return Unwrap(ret, result)
	}
}

func (op *tracesStorageOp) UpdateKey(ctx context.Context, tracesResourceId string, id uuid.UUID, description *string) (*v1.TraceStorageAccessKey, error) {
	rid, err := strconv.ParseInt(tracesResourceId, 10, 64)
	if err != nil {
		return nil, NewError("TracesStorage.UpdateKey", err)
	}
	params := v1.TracesStoragesKeysUpdateParams{TraceResourceID: rid, UID: id}
	opt := v1.NewOptTraceStorageAccessKey(v1.TraceStorageAccessKey{
		Description: intoOpt[v1.OptString](description),
	})
	result, err := op.client.TracesStoragesKeysUpdate(ctx, opt, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("TracesStorage.UpdateKey", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return nil, NewAPIError("TracesStorage.UpdateKey", e.StatusCode, errors.Wrap(err, "invalid parameter"))
		default:
			return nil, NewAPIError("TracesStorage.UpdateKey", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("TracesStorage.UpdateKey", 0, err)
	} else {
		ret := new(v1.TraceStorageAccessKey)
		return Unwrap(ret, result)
	}
}

// DeleteKey deletes an access key for a traces storage resource.
func (op *tracesStorageOp) DeleteKey(ctx context.Context, tracesResourceId string, id uuid.UUID) error {
	rid, err := strconv.ParseInt(tracesResourceId, 10, 64)
	if err != nil {
		return NewError("TracesStorage.DeleteKey", err)
	}
	params := v1.TracesStoragesKeysDestroyParams{TraceResourceID: rid, UID: id}
	err = op.client.TracesStoragesKeysDestroy(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return NewAPIError("TracesStorage.DeleteKey", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusBadRequest:
			return NewAPIError("TracesStorage.DeleteKey", e.StatusCode, errors.Wrap(err, "the request resource is not eligible for deletion"))
		default:
			return NewAPIError("TracesStorage.DeleteKey", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return NewAPIError("TracesStorage.DeleteKey", 0, err)
	}
	return nil
}

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
	"strconv"
	"time"

	"github.com/google/uuid"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
)

type TracesStorageAPI interface {
	List(ctx context.Context, params TracesStorageListParams) ([]v1.TraceStorage, error)
	Create(ctx context.Context, request TracesStorageCreateParams) (*v1.TraceStorage, error)
	Read(ctx context.Context, id string) (*v1.TraceStorage, error)
	Update(ctx context.Context, id string, request TracesStorageUpdateParams) (*v1.TraceStorage, error)
	Delete(ctx context.Context, id string) error

	SetExpire(ctx context.Context, resourceID string, days int) (*v1.TraceStorage, error)
	ReadDailyStats(ctx context.Context, resourceID string, startDate, endDate *time.Time) ([]v1.LogStorageDailyUsage, error)
	ReadMonthlyStats(ctx context.Context, resourceID string, year int) ([]v1.LogStorageMonthlyUsage, error)

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

func (op *tracesStorageOp) List(ctx context.Context, params TracesStorageListParams) (ret []v1.TraceStorage, err error) {
	res, err := errorFromDecodedResponse("TracesStorage.List", func() (*v1.PaginatedTraceStorageList, error) {
		resourceId, err := fromStringPtr[v1.OptInt64, int64](params.ResourceID)
		if err != nil {
			return nil, err
		}
		return op.client.TracesStoragesList(ctx, v1.TracesStoragesListParams{
			Count:                          intoOpt[v1.OptInt](params.Count),
			From:                           intoOpt[v1.OptInt](params.From),
			AccountID:                      intoOpt[v1.OptString](params.AccountID),
			ResourceID:                     resourceId,
			LogStorageBucketClassification: intoOpt[v1.OptTracesStoragesListLogStorageBucketClassification](params.BucketClassification),
		})
	})
	if err == nil {
		ret = res.GetResults()
	}
	return
}

func (op *tracesStorageOp) Read(ctx context.Context, resourceID string) (*v1.TraceStorage, error) {
	res, err := errorFromDecodedResponse("TracesStorage.Read", func() (*v1.WrappedTraceStorage, error) {
		id, err := strconv.ParseInt(resourceID, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.TracesStoragesRetrieve(ctx, v1.TracesStoragesRetrieveParams{ResourceID: id})
	})
	return unwrapE[*v1.TraceStorage](res, err)
}

type TracesStorageCreateParams struct {
	Name           string
	Description    *string
	Classification *v1.TraceStorageCreateClassification
}

func (op *tracesStorageOp) Create(ctx context.Context, params TracesStorageCreateParams) (*v1.TraceStorage, error) {
	res, err := errorFromDecodedResponse("TracesStorage.Create", func() (*v1.TraceStorage, error) {
		req := v1.TraceStorageCreate{
			Name:           params.Name,
			Description:    intoOpt[v1.OptString](params.Description),
			Classification: intoOpt[v1.OptTraceStorageCreateClassification](params.Classification),
		}
		return op.client.TracesStoragesCreate(ctx, &req)
	})
	return unwrapE[*v1.TraceStorage](res, err)
}

type TracesStorageUpdateParams struct {
	Name        *string
	Description *string
}

func (op *tracesStorageOp) Update(ctx context.Context, id string, p TracesStorageUpdateParams) (*v1.TraceStorage, error) {
	res, err := errorFromDecodedResponse("TracesStorage.Update", func() (*v1.WrappedTraceStorage, error) {
		rid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.TracesStoragesPartialUpdate(ctx, v1.NewOptPatchedTraceStorage(v1.PatchedTraceStorage{
			Name:        intoOpt[v1.OptString](p.Name),
			Description: intoOpt[v1.OptString](p.Description),
		}), v1.TracesStoragesPartialUpdateParams{ResourceID: rid})
	})
	return unwrapE[*v1.TraceStorage](res, err)
}

func (op *tracesStorageOp) Delete(ctx context.Context, id string) error {
	return errorFromDecodedResponse1("TracesStorage.Delete", func() error {
		rid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return err
		}
		return op.client.TracesStoragesDestroy(ctx, v1.TracesStoragesDestroyParams{ResourceID: rid})
	})
}

func (op *tracesStorageOp) SetExpire(ctx context.Context, resourceID string, days int) (*v1.TraceStorage, error) {
	res, err := errorFromDecodedResponse("TracesStorage.SetExpire", func() (*v1.TraceStorage, error) {
		rid, err := strconv.ParseInt(resourceID, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.TracesStoragesSetExpireCreate(
			ctx,
			&v1.SetTraceStorageExpireDay{Days: days},
			v1.TracesStoragesSetExpireCreateParams{ResourceID: rid},
		)
	})
	return unwrapE[*v1.TraceStorage](res, err)
}

func (op *tracesStorageOp) ReadDailyStats(ctx context.Context, resourceID string, startDate, endDate *time.Time) (ret []v1.LogStorageDailyUsage, err error) {
	res, err := errorFromDecodedResponse("TracesStorage.ReadDailyStats", func() (*v1.TraceStorageDailyUsageBody, error) {
		rid, err := strconv.ParseInt(resourceID, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.TracesStoragesStatsDailyRetrieve(ctx, v1.TracesStoragesStatsDailyRetrieveParams{
			ResourceID: rid,
			StartDate:  intoOpt[v1.OptDate](startDate),
			EndDate:    intoOpt[v1.OptDate](endDate),
		})
	})
	if err == nil {
		ret = res.GetUsages()
	}
	return
}

func (op *tracesStorageOp) ReadMonthlyStats(ctx context.Context, resourceID string, year int) (ret []v1.LogStorageMonthlyUsage, err error) {
	res, err := errorFromDecodedResponse("TracesStorage.ReadMonthlyStats", func() (*v1.TraceStorageMonthlyUsageBody, error) {
		rid, err := strconv.ParseInt(resourceID, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.TracesStoragesStatsMonthlyRetrieve(ctx, v1.TracesStoragesStatsMonthlyRetrieveParams{
			ResourceID: rid,
			Year:       year,
		})
	})
	if err == nil {
		ret = res.GetUsages()
	}
	return
}

func (op *tracesStorageOp) ListKeys(ctx context.Context, tracesResourceId string, count *int, from *int) (ret []v1.TraceStorageAccessKey, err error) {
	res, err := errorFromDecodedResponse("TracesStorage.ListKeys", func() (*v1.PaginatedTraceStorageAccessKeyList, error) {
		rid, err := strconv.ParseInt(tracesResourceId, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.TracesStoragesKeysList(ctx, v1.TracesStoragesKeysListParams{
			TraceResourceID: rid,
			Count:           intoOpt[v1.OptInt](count),
			From:            intoOpt[v1.OptInt](from),
		})
	})
	if err == nil {
		ret = res.GetResults()
	}
	return
}

func (op *tracesStorageOp) CreateKey(ctx context.Context, tracesResourceId string, description *string) (*v1.TraceStorageAccessKey, error) {
	res, err := errorFromDecodedResponse("TracesStorage.CreateKey", func() (*v1.WrappedTraceStorageAccessKey, error) {
		rid, err := strconv.ParseInt(tracesResourceId, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.TracesStoragesKeysCreate(ctx, v1.NewOptTraceStorageAccessKey(v1.TraceStorageAccessKey{
			Description: intoOpt[v1.OptString](description),
		}), v1.TracesStoragesKeysCreateParams{TraceResourceID: rid})
	})
	return unwrapE[*v1.TraceStorageAccessKey](res, err)
}

func (op *tracesStorageOp) ReadKey(ctx context.Context, tracesResourceId string, id uuid.UUID) (*v1.TraceStorageAccessKey, error) {
	res, err := errorFromDecodedResponse("TracesStorage.ReadKey", func() (*v1.WrappedTraceStorageAccessKey, error) {
		rid, err := strconv.ParseInt(tracesResourceId, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.TracesStoragesKeysRetrieve(ctx, v1.TracesStoragesKeysRetrieveParams{
			TraceResourceID: rid,
			UID:             id,
		})
	})
	return unwrapE[*v1.TraceStorageAccessKey](res, err)
}

func (op *tracesStorageOp) UpdateKey(ctx context.Context, tracesResourceId string, id uuid.UUID, description *string) (*v1.TraceStorageAccessKey, error) {
	res, err := errorFromDecodedResponse("TracesStorage.UpdateKey", func() (*v1.WrappedTraceStorageAccessKey, error) {
		rid, err := strconv.ParseInt(tracesResourceId, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.TracesStoragesKeysUpdate(ctx, v1.NewOptTraceStorageAccessKey(v1.TraceStorageAccessKey{
			Description: intoOpt[v1.OptString](description),
		}), v1.TracesStoragesKeysUpdateParams{
			TraceResourceID: rid,
			UID:             id,
		})
	})
	return unwrapE[*v1.TraceStorageAccessKey](res, err)
}

// DeleteKey deletes an access key for a traces storage resource.
func (op *tracesStorageOp) DeleteKey(ctx context.Context, tracesResourceId string, id uuid.UUID) error {
	return errorFromDecodedResponse1("TracesStorage.DeleteKey", func() error {
		rid, err := strconv.ParseInt(tracesResourceId, 10, 64)
		if err != nil {
			return err
		}
		return op.client.TracesStoragesKeysDestroy(ctx, v1.TracesStoragesKeysDestroyParams{
			TraceResourceID: rid,
			UID:             id,
		})
	})
}

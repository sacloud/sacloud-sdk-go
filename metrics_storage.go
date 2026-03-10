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

func (op *metricsStorageOp) List(ctx context.Context, params MetricsStorageListParams) (ret []v1.MetricsStorage, err error) {
	res, err := errorFromDecodedResponse("MetricsStorage.List", func() (*v1.PaginatedMetricsStorageList, error) {
		resourceId, err := fromStringPtr[v1.OptInt64, int64](params.ResourceID)
		if err != nil {
			return nil, err
		}
		return op.client.MetricsStoragesList(ctx, v1.MetricsStoragesListParams{
			Count:      intoOpt[v1.OptInt](params.Count),
			From:       intoOpt[v1.OptInt](params.From),
			AccountID:  intoOpt[v1.OptString](params.AccountID),
			ResourceID: resourceId,
			IsSystem:   intoOpt[v1.OptBool](params.IsSystem),
		})
	})
	if err == nil {
		ret = res.GetResults()
	}
	return
}

func (op *metricsStorageOp) Read(ctx context.Context, resourceID string) (*v1.MetricsStorage, error) {
	res, err := errorFromDecodedResponse("MetricsStorage.Read", func() (*v1.WrappedMetricsStorage, error) {
		id, err := strconv.ParseInt(resourceID, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.MetricsStoragesRetrieve(ctx, v1.MetricsStoragesRetrieveParams{ResourceID: id})
	})
	return unwrapE[*v1.MetricsStorage](res, err)
}

type MetricsStorageCreateParams struct {
	Name        string
	Description *string
	IsSystem    bool
}

func (op *metricsStorageOp) Create(ctx context.Context, params MetricsStorageCreateParams) (*v1.MetricsStorage, error) {
	res, err := errorFromDecodedResponse("MetricsStorage.Create", func() (*v1.MetricsStorage, error) {
		return op.client.MetricsStoragesCreate(ctx, &v1.MetricsStorageCreate{
			Name:        params.Name,
			Description: intoOpt[v1.OptString](params.Description),
			IsSystem:    params.IsSystem,
		})
	})
	return unwrapE[*v1.MetricsStorage](res, err)
}

type MetricsStorageUpdateParams struct {
	Name        *string
	Description *string
}

func (op *metricsStorageOp) Update(ctx context.Context, id string, params MetricsStorageUpdateParams) (*v1.MetricsStorage, error) {
	res, err := errorFromDecodedResponse("MetricsStorage.Update", func() (*v1.WrappedMetricsStorage, error) {
		rid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.MetricsStoragesPartialUpdate(ctx, v1.NewOptPatchedMetricsStorage(v1.PatchedMetricsStorage{
			Name:        intoOpt[v1.OptString](params.Name),
			Description: intoOpt[v1.OptString](params.Description),
		}), v1.MetricsStoragesPartialUpdateParams{ResourceID: rid})
	})
	return unwrapE[*v1.MetricsStorage](res, err)
}

func (op *metricsStorageOp) Delete(ctx context.Context, resourceID string) error {
	return errorFromDecodedResponse1("MetricsStorage.Delete", func() error {
		rid, err := strconv.ParseInt(resourceID, 10, 64)
		if err != nil {
			return err
		}
		return op.client.MetricsStoragesDestroy(ctx, v1.MetricsStoragesDestroyParams{ResourceID: rid})
	})
}

func (op *metricsStorageOp) ReadDailyStats(ctx context.Context, resourceID string, startDate, endDate *time.Time) (ret []v1.MetricsStorageDailyUsage, err error) {
	res, err := errorFromDecodedResponse("MetricsStorage.ReadDailyStats", func() (*v1.MetricsStorageDailyUsageBody, error) {
		rid, err := strconv.ParseInt(resourceID, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.MetricsStoragesStatsDailyRetrieve(ctx, v1.MetricsStoragesStatsDailyRetrieveParams{
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

func (op *metricsStorageOp) ReadMonthlyStats(ctx context.Context, resourceID string, year int) (ret []v1.MetricsStorageMonthlyUsage, err error) {
	res, err := errorFromDecodedResponse("MetricsStorage.ReadMonthlyStats", func() (*v1.MetricsStorageMonthlyUsageBody, error) {
		rid, err := strconv.ParseInt(resourceID, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.MetricsStoragesStatsMonthlyRetrieve(ctx, v1.MetricsStoragesStatsMonthlyRetrieveParams{
			ResourceID: rid,
			Year:       year,
		})
	})
	if err == nil {
		ret = res.GetUsages()
	}
	return
}

func (op *metricsStorageOp) ListKeys(ctx context.Context, metricsResourceId string, count *int, from *int) (ret []v1.MetricsStorageAccessKey, err error) {
	res, err := errorFromDecodedResponse("MetricsStorage.ListKeys", func() (*v1.PaginatedMetricsStorageAccessKeyList, error) {
		rid, err := strconv.ParseInt(metricsResourceId, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.MetricsStoragesKeysList(ctx, v1.MetricsStoragesKeysListParams{
			MetricsResourceID: rid,
			Count:             intoOpt[v1.OptInt](count),
			From:              intoOpt[v1.OptInt](from),
		})
	})
	if err == nil {
		ret = res.GetResults()
	}
	return
}

func (op *metricsStorageOp) CreateKey(ctx context.Context, metricsResourceId string, description *string) (*v1.MetricsStorageAccessKey, error) {
	res, err := errorFromDecodedResponse("MetricsStorage.CreateKey", func() (*v1.WrappedMetricsStorageAccessKey, error) {
		rid, err := strconv.ParseInt(metricsResourceId, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.MetricsStoragesKeysCreate(ctx, v1.NewOptMetricsStorageAccessKey(v1.MetricsStorageAccessKey{
			Description: intoOpt[v1.OptString](description),
		}), v1.MetricsStoragesKeysCreateParams{MetricsResourceID: rid})
	})
	return unwrapE[*v1.MetricsStorageAccessKey](res, err)
}

func (op *metricsStorageOp) ReadKey(ctx context.Context, metricsResourceId string, id uuid.UUID) (*v1.MetricsStorageAccessKey, error) {
	res, err := errorFromDecodedResponse("MetricsStorage.ReadKey", func() (*v1.WrappedMetricsStorageAccessKey, error) {
		rid, err := strconv.ParseInt(metricsResourceId, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.MetricsStoragesKeysRetrieve(ctx, v1.MetricsStoragesKeysRetrieveParams{
			MetricsResourceID: rid,
			UID:               id,
		})
	})
	return unwrapE[*v1.MetricsStorageAccessKey](res, err)
}

func (op *metricsStorageOp) UpdateKey(ctx context.Context, metricsResourceId string, id uuid.UUID, description *string) (*v1.MetricsStorageAccessKey, error) {
	res, err := errorFromDecodedResponse("MetricsStorage.UpdateKey", func() (*v1.WrappedMetricsStorageAccessKey, error) {
		rid, err := strconv.ParseInt(metricsResourceId, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.MetricsStoragesKeysUpdate(ctx, v1.NewOptMetricsStorageAccessKey(v1.MetricsStorageAccessKey{
			Description: intoOpt[v1.OptString](description),
		}), v1.MetricsStoragesKeysUpdateParams{
			MetricsResourceID: rid,
			UID:               id,
		})
	})
	return unwrapE[*v1.MetricsStorageAccessKey](res, err)
}

// DeleteKey deletes an access key for a metrics storage resource.
func (op *metricsStorageOp) DeleteKey(ctx context.Context, metricsResourceId string, id uuid.UUID) error {
	return errorFromDecodedResponse1("MetricsStorage.DeleteKey", func() error {
		rid, err := strconv.ParseInt(metricsResourceId, 10, 64)
		if err != nil {
			return err
		}
		return op.client.MetricsStoragesKeysDestroy(ctx, v1.MetricsStoragesKeysDestroyParams{
			MetricsResourceID: rid,
			UID:               id,
		})
	})
}

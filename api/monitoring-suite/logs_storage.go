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
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
)

type LogsStorageAPI interface {
	List(ctx context.Context, params LogsStoragesListParams) ([]v1.LogStorage, error)
	Create(ctx context.Context, params LogStorageCreateParams) (*v1.LogStorage, error)
	Read(ctx context.Context, id string) (*v1.LogStorage, error)
	Update(ctx context.Context, id string, params LogStorageUpdateParams) (*v1.LogStorage, error)
	Delete(ctx context.Context, id string) error

	SetExpire(ctx context.Context, resourceID string, days int) (*v1.LogStorage, error)
	ReadDailyStats(ctx context.Context, resourceID string, startDate, endDate *time.Time) ([]v1.LogStorageDailyUsage, error)
	ReadMonthlyStats(ctx context.Context, resourceID string, year int) ([]v1.LogStorageMonthlyUsage, error)

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

func (op *logsStorageOp) List(ctx context.Context, p LogsStoragesListParams) (ret []v1.LogStorage, err error) {
	res, err := errorFromDecodedResponse("LogsStorage.List", func() (*v1.PaginatedLogStorageList, error) {
		id, err := fromStringPtr[v1.OptInt64, int64](p.ResourceID)
		if err != nil {
			return nil, fmt.Errorf("LogsStoragesListParams.ResourceID: %w", err)
		}
		return op.client.LogsStoragesList(ctx, v1.LogsStoragesListParams{
			AccountID:            intoOpt[v1.OptString](p.AccountID),
			BucketClassification: intoOpt[v1.OptLogsStoragesListBucketClassification](p.BucketClassification),
			Count:                intoOpt[v1.OptInt](p.Count),
			From:                 intoOpt[v1.OptInt](p.From),
			IsSystem:             intoOpt[v1.OptBool](p.IsSystem),
			ResourceID:           id,
			Status:               intoOpt[v1.OptLogsStoragesListStatus](p.Status),
		})
	})
	if err == nil {
		ret = res.GetResults()
	}
	return
}

func (op *logsStorageOp) Read(ctx context.Context, resourceID string) (*v1.LogStorage, error) {
	res, err := errorFromDecodedResponse("LogsStorage.Read", func() (*v1.WrappedLogStorage, error) {
		id, err := strconv.ParseInt(resourceID, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.LogsStoragesRetrieve(ctx, v1.LogsStoragesRetrieveParams{ResourceID: id})
	})
	return unwrapE[*v1.LogStorage](res, err)
}

type LogStorageCreateParams struct {
	Name           string
	Description    *string
	IsSystem       bool
	Classification *v1.LogStorageCreateRequestClassification
}

func (op *logsStorageOp) Create(ctx context.Context, params LogStorageCreateParams) (*v1.LogStorage, error) {
	res, err := errorFromDecodedResponse("LogsStorage.Create", func() (*v1.LogStorage, error) {
		req := v1.LogStorageCreateRequest{
			Classification: intoOpt[v1.OptLogStorageCreateRequestClassification](params.Classification),
			IsSystem:       params.IsSystem,
			Name:           params.Name,
			Description:    intoOpt[v1.OptString](params.Description),
		}
		return op.client.LogsStoragesCreate(ctx, &req)
	})
	return unwrapE[*v1.LogStorage](res, err)
}

type LogStorageUpdateParams struct {
	Name        *string
	Description *string
}

func (op *logsStorageOp) Update(ctx context.Context, id string, p LogStorageUpdateParams) (*v1.LogStorage, error) {
	res, err := errorFromDecodedResponse("LogsStorage.Update", func() (*v1.WrappedLogStorage, error) {
		rid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.LogsStoragesPartialUpdate(ctx, v1.NewOptPatchedLogStorageRequest(v1.PatchedLogStorageRequest{
			Name:        intoOpt[v1.OptString](p.Name),
			Description: intoOpt[v1.OptString](p.Description),
		}), v1.LogsStoragesPartialUpdateParams{ResourceID: rid})
	})
	return unwrapE[*v1.LogStorage](res, err)
}

func (op *logsStorageOp) Delete(ctx context.Context, id string) error {
	return errorFromDecodedResponse1("LogsStorage.Delete", func() error {
		rid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return err
		}
		return op.client.LogsStoragesDestroy(ctx, v1.LogsStoragesDestroyParams{ResourceID: rid})
	})
}

func (op *logsStorageOp) SetExpire(ctx context.Context, resourceID string, days int) (*v1.LogStorage, error) {
	res, err := errorFromDecodedResponse("LogsStorage.SetExpire", func() (*v1.LogStorage, error) {
		rid, err := strconv.ParseInt(resourceID, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.LogsStoragesSetExpireCreate(
			ctx,
			&v1.SetLogStorageExpireDayRequest{Days: days},
			v1.LogsStoragesSetExpireCreateParams{ResourceID: rid},
		)
	})
	return unwrapE[*v1.LogStorage](res, err)
}

func (op *logsStorageOp) ReadDailyStats(ctx context.Context, resourceID string, startDate, endDate *time.Time) (ret []v1.LogStorageDailyUsage, err error) {
	res, err := errorFromDecodedResponse("LogsStorage.ReadDailyStats", func() (*v1.LogStorageDailyUsageBody, error) {
		rid, err := strconv.ParseInt(resourceID, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.LogsStoragesStatsDailyRetrieve(ctx, v1.LogsStoragesStatsDailyRetrieveParams{
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

func (op *logsStorageOp) ReadMonthlyStats(ctx context.Context, resourceID string, year int) (ret []v1.LogStorageMonthlyUsage, err error) {
	res, err := errorFromDecodedResponse("LogsStorage.ReadMonthlyStats", func() (*v1.LogStorageMonthlyUsageBody, error) {
		rid, err := strconv.ParseInt(resourceID, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.LogsStoragesStatsMonthlyRetrieve(ctx, v1.LogsStoragesStatsMonthlyRetrieveParams{
			ResourceID: rid,
			Year:       year,
		})
	})
	if err == nil {
		ret = res.GetUsages()
	}
	return
}

func (op *logsStorageOp) ListKeys(ctx context.Context, logResourceId string, count *int, from *int) (ret []v1.LogStorageAccessKey, err error) {
	res, err := errorFromDecodedResponse("LogsStorage.ListKeys", func() (*v1.PaginatedLogStorageAccessKeyList, error) {
		rid, err := strconv.ParseInt(logResourceId, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.LogsStoragesKeysList(ctx, v1.LogsStoragesKeysListParams{
			Count:         intoOpt[v1.OptInt](count),
			From:          intoOpt[v1.OptInt](from),
			LogResourceID: rid,
		})
	})
	if err == nil {
		ret = res.GetResults()
	}
	return
}

func (op *logsStorageOp) CreateKey(ctx context.Context, logResourceId string, description *string) (*v1.LogStorageAccessKey, error) {
	res, err := errorFromDecodedResponse("LogsStorage.CreateKey", func() (*v1.WrappedLogStorageAccessKey, error) {
		rid, err := strconv.ParseInt(logResourceId, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.LogsStoragesKeysCreate(ctx, v1.NewOptLogStorageAccessKeyRequest(v1.LogStorageAccessKeyRequest{
			Description: intoOpt[v1.OptString](description),
		}), v1.LogsStoragesKeysCreateParams{LogResourceID: rid})
	})
	return unwrapE[*v1.LogStorageAccessKey](res, err)
}

func (op *logsStorageOp) ReadKey(ctx context.Context, logResourceId string, id uuid.UUID) (*v1.LogStorageAccessKey, error) {
	res, err := errorFromDecodedResponse("LogsStorage.ReadKey", func() (*v1.WrappedLogStorageAccessKey, error) {
		rid, err := strconv.ParseInt(logResourceId, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.LogsStoragesKeysRetrieve(ctx, v1.LogsStoragesKeysRetrieveParams{
			LogResourceID: rid,
			UID:           id,
		})
	})
	return unwrapE[*v1.LogStorageAccessKey](res, err)
}

func (op *logsStorageOp) UpdateKey(ctx context.Context, logResourceId string, id uuid.UUID, description *string) (*v1.LogStorageAccessKey, error) {
	res, err := errorFromDecodedResponse("LogsStorage.UpdateKey", func() (*v1.WrappedLogStorageAccessKey, error) {
		rid, err := strconv.ParseInt(logResourceId, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.LogsStoragesKeysUpdate(ctx, v1.NewOptLogStorageAccessKeyRequest(v1.LogStorageAccessKeyRequest{
			Description: intoOpt[v1.OptString](description),
		}), v1.LogsStoragesKeysUpdateParams{
			LogResourceID: rid,
			UID:           id,
		})
	})
	return unwrapE[*v1.LogStorageAccessKey](res, err)
}

func (op *logsStorageOp) DeleteKey(ctx context.Context, logResourceId string, id uuid.UUID) error {
	return errorFromDecodedResponse1("LogsStorage.DeleteKey", func() error {
		rid, err := strconv.ParseInt(logResourceId, 10, 64)
		if err != nil {
			return err
		}
		return op.client.LogsStoragesKeysDestroy(ctx, v1.LogsStoragesKeysDestroyParams{
			LogResourceID: rid,
			UID:           id,
		})
	})
}

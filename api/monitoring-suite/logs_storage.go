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
	v1 "github.com/sacloud/sacloud-sdk-go/api/monitoring-suite/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/packages/into"
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
	BucketClassification *v1.ListLogStoragesBucketClassification
	Count                *int
	From                 *int
	IsSystem             *bool
	ResourceID           *string
	Status               *v1.ListLogStoragesStatus
}

func (op *logsStorageOp) List(ctx context.Context, p LogsStoragesListParams) (ret []v1.LogStorage, err error) {
	res, err := errorFromDecodedResponse("LogsStorage.List", func() (*v1.PaginatedLogStorageList, error) {
		id, err := into.FromStringPtr[v1.OptInt64, int64](p.ResourceID)
		if err != nil {
			return nil, fmt.Errorf("LogsStoragesListParams.ResourceID: %w", err)
		}
		return op.client.ListLogStorages(ctx, v1.ListLogStoragesParams{
			AccountID:            into.Opt[v1.OptString](p.AccountID),
			BucketClassification: into.Opt[v1.OptListLogStoragesBucketClassification](p.BucketClassification),
			Count:                into.Opt[v1.OptInt](p.Count),
			From:                 into.Opt[v1.OptInt](p.From),
			IsSystem:             into.Opt[v1.OptBool](p.IsSystem),
			ResourceID:           id,
			Status:               into.Opt[v1.OptListLogStoragesStatus](p.Status),
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
		return op.client.ReadLogStorage(ctx, v1.ReadLogStorageParams{ResourceId: id})
	})
	return unwrapE[*v1.LogStorage](res, err)
}

type LogStorageCreateParams struct {
	Name               string
	Description        *string
	IsSystem           bool
	Classification     *v1.LogStorageCreateRequestClassification
	KmsKeyID           v1.OptNilInt64
	ServicePrincipalID v1.OptNilInt64
}

func (op *logsStorageOp) Create(ctx context.Context, params LogStorageCreateParams) (*v1.LogStorage, error) {
	res, err := errorFromDecodedResponse("LogsStorage.Create", func() (*v1.LogStorage, error) {
		req := v1.LogStorageCreateRequest{
			Classification:     into.Opt[v1.OptLogStorageCreateRequestClassification](params.Classification),
			IsSystem:           params.IsSystem,
			Name:               params.Name,
			Description:        into.Opt[v1.OptString](params.Description),
			KmsKeyID:           params.KmsKeyID,
			ServicePrincipalID: params.ServicePrincipalID,
		}
		return op.client.CreateLogStorage(ctx, &req)
	})
	return unwrapE[*v1.LogStorage](res, err)
}

type LogStorageUpdateParams struct {
	Name               *string
	Description        *string
	ServicePrincipalID v1.OptNilInt64
}

func (op *logsStorageOp) Update(ctx context.Context, id string, p LogStorageUpdateParams) (*v1.LogStorage, error) {
	res, err := errorFromDecodedResponse("LogsStorage.Update", func() (*v1.WrappedLogStorage, error) {
		rid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.PatchLogStorage(ctx, v1.NewOptPatchedLogStorageRequest(v1.PatchedLogStorageRequest{
			Name:               into.Opt[v1.OptString](p.Name),
			Description:        into.Opt[v1.OptString](p.Description),
			ServicePrincipalID: p.ServicePrincipalID,
		}), v1.PatchLogStorageParams{ResourceId: rid})
	})
	return unwrapE[*v1.LogStorage](res, err)
}

func (op *logsStorageOp) Delete(ctx context.Context, id string) error {
	return errorFromDecodedResponse1("LogsStorage.Delete", func() error {
		rid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return err
		}
		return op.client.DeleteLogStorage(ctx, v1.DeleteLogStorageParams{ResourceId: rid})
	})
}

func (op *logsStorageOp) SetExpire(ctx context.Context, resourceID string, days int) (*v1.LogStorage, error) {
	res, err := errorFromDecodedResponse("LogsStorage.SetExpire", func() (*v1.LogStorage, error) {
		rid, err := strconv.ParseInt(resourceID, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.SetLogStorageExpire(
			ctx,
			&v1.SetLogStorageExpireDayRequest{Days: days},
			v1.SetLogStorageExpireParams{ResourceId: rid},
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
		return op.client.ReadLogStorageStatsDaily(ctx, v1.ReadLogStorageStatsDailyParams{
			ResourceId: rid,
			StartDate:  into.Opt[v1.OptDate](startDate),
			EndDate:    into.Opt[v1.OptDate](endDate),
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
		return op.client.ReadLogStorageStatsMonthly(ctx, v1.ReadLogStorageStatsMonthlyParams{
			ResourceId: rid,
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
		return op.client.ListLogStorageKeys(ctx, v1.ListLogStorageKeysParams{
			Count:         into.Opt[v1.OptInt](count),
			From:          into.Opt[v1.OptInt](from),
			LogResourceId: rid,
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
		return op.client.CreateLogStorageKey(ctx, v1.NewOptLogStorageAccessKeyRequest(v1.LogStorageAccessKeyRequest{
			Description: into.Opt[v1.OptString](description),
		}), v1.CreateLogStorageKeyParams{LogResourceId: rid})
	})
	return unwrapE[*v1.LogStorageAccessKey](res, err)
}

func (op *logsStorageOp) ReadKey(ctx context.Context, logResourceId string, id uuid.UUID) (*v1.LogStorageAccessKey, error) {
	res, err := errorFromDecodedResponse("LogsStorage.ReadKey", func() (*v1.WrappedLogStorageAccessKey, error) {
		rid, err := strconv.ParseInt(logResourceId, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.ReadLogStorageKey(ctx, v1.ReadLogStorageKeyParams{
			LogResourceId: rid,
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
		return op.client.UpdateLogStorageKey(ctx, v1.NewOptLogStorageAccessKeyRequest(v1.LogStorageAccessKeyRequest{
			Description: into.Opt[v1.OptString](description),
		}), v1.UpdateLogStorageKeyParams{
			LogResourceId: rid,
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
		return op.client.DeleteLogStorageKey(ctx, v1.DeleteLogStorageKeyParams{
			LogResourceId: rid,
			UID:           id,
		})
	})
}

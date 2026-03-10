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

type LogMeasureRuleAPI interface {
	List(ctx context.Context, projectId string, count *int, from *int) ([]v1.LogMeasureRule, error)
	Create(ctx context.Context, projectId string, params LogMeasureRuleCreateParams) (*v1.LogMeasureRule, error)
	Read(ctx context.Context, projectId string, ruleId uuid.UUID) (*v1.LogMeasureRule, error)
	Update(ctx context.Context, projectId string, ruleId uuid.UUID, params LogMeasureRuleUpdateParams) (*v1.LogMeasureRule, error)
	Delete(ctx context.Context, projectId string, ruleId uuid.UUID) error

	ListHistories(ctx context.Context, projectId string, params LogMeasureRuleListHistoriesParams) ([]v1.History, error)
	ReadHistory(ctx context.Context, projectId string, historyId uuid.UUID) (*v1.History, error)
}

var _ LogMeasureRuleAPI = (*logMeasureRuleOp)(nil)

type logMeasureRuleOp struct {
	client *v1.Client
}

func NewLogMeasureRuleOp(client *v1.Client) LogMeasureRuleAPI {
	return &logMeasureRuleOp{client: client}
}

func (op *logMeasureRuleOp) List(ctx context.Context, projectId string, count *int, from *int) (ret []v1.LogMeasureRule, err error) {
	res, err := errorFromDecodedResponse("LogMeasureRule.List", func() (*v1.PaginatedLogMeasureRuleList, error) {
		if id, err := strconv.ParseInt(projectId, 10, 64); err != nil {
			return nil, err
		} else {
			return op.client.AlertsProjectsLogMeasureRulesList(ctx, v1.AlertsProjectsLogMeasureRulesListParams{
				ProjectResourceID: id,
				Count:             intoOpt[v1.OptInt](count),
				From:              intoOpt[v1.OptInt](from),
			})
		}
	})
	if err == nil {
		ret = res.GetResults()
	}
	return
}

type LogMeasureRuleCreateParams struct {
	LogStorageID     string
	MetricsStorageID string
	Name             *string
	Description      *string
	Rule             v1.LogMeasureRuleModel
}

func (op *logMeasureRuleOp) Create(ctx context.Context, projectId string, p LogMeasureRuleCreateParams) (*v1.LogMeasureRule, error) {
	return errorFromDecodedResponse("LogMeasureRule.Create", func() (*v1.LogMeasureRule, error) {
		if pid, err := strconv.ParseInt(projectId, 10, 64); err != nil {
			return nil, fmt.Errorf("projectId: %w", err)
		} else if lid, err := strconv.ParseInt(p.LogStorageID, 10, 64); err != nil {
			return nil, fmt.Errorf("LogMeasureRuleCreateParams.LogStorageID: %w", err)
		} else if mid, err := strconv.ParseInt(p.MetricsStorageID, 10, 64); err != nil {
			return nil, fmt.Errorf("LogMeasureRuleCreateParams.MetricsStorageID: %w", err)
		} else {
			params := &v1.LogMeasureRule{
				LogStorageID:     v1.NewOptNilInt64(lid),
				MetricsStorageID: v1.NewOptNilInt64(mid),
				Name:             intoOpt[v1.OptString](p.Name),
				Description:      intoOpt[v1.OptString](p.Description),
				Rule:             p.Rule,
			}
			// prevent ogen error (encoder is not accepting empty struct)
			params.LogStorage.SetFake()
			params.MetricsStorage.SetFake()
			params.LogStorage.SetTags(make([]string, 0))
			params.MetricsStorage.SetTags(make([]string, 0))
			return op.client.AlertsProjectsLogMeasureRulesCreate(ctx, params, v1.AlertsProjectsLogMeasureRulesCreateParams{
				ProjectResourceID: pid,
			})
		}
	})
}

func (op *logMeasureRuleOp) Read(ctx context.Context, projectId string, ruleId uuid.UUID) (*v1.LogMeasureRule, error) {
	return errorFromDecodedResponse("LogMeasureRule.Read", func() (*v1.LogMeasureRule, error) {
		if pid, err := strconv.ParseInt(projectId, 10, 64); err != nil {
			return nil, err
		} else {
			return op.client.AlertsProjectsLogMeasureRulesRetrieve(ctx, v1.AlertsProjectsLogMeasureRulesRetrieveParams{
				ProjectResourceID: pid,
				UID:               ruleId,
			})
		}
	})
}

type LogMeasureRuleUpdateParams struct {
	LogStorageID     *string
	MetricsStorageID *string
	Name             *string
	Description      *string
	Rule             *v1.LogMeasureRuleModel
}

func (op *logMeasureRuleOp) Update(ctx context.Context, projectId string, ruleId uuid.UUID, p LogMeasureRuleUpdateParams) (*v1.LogMeasureRule, error) {
	return errorFromDecodedResponse("LogMeasureRule.Update", func() (*v1.LogMeasureRule, error) {
		if pid, err := strconv.ParseInt(projectId, 10, 64); err != nil {
			return nil, fmt.Errorf("projectId: %w", err)
		} else if lid, err := fromStringPtr[v1.OptNilInt64, int64](p.LogStorageID); err != nil {
			return nil, fmt.Errorf("LogMeasureRuleUpdateParams.LogStorageID: %w", err)
		} else if mid, err := fromStringPtr[v1.OptNilInt64, int64](p.MetricsStorageID); err != nil {
			return nil, fmt.Errorf("LogMeasureRuleUpdateParams.MetricsStorageID: %w", err)
		} else {
			return op.client.AlertsProjectsLogMeasureRulesPartialUpdate(ctx, v1.NewOptPatchedLogMeasureRule(v1.PatchedLogMeasureRule{
				LogStorageID:     lid,
				MetricsStorageID: mid,
				Name:             intoOpt[v1.OptString](p.Name),
				Description:      intoOpt[v1.OptString](p.Description),
				Rule:             intoOpt[v1.OptLogMeasureRuleModel](p.Rule),
			}), v1.AlertsProjectsLogMeasureRulesPartialUpdateParams{
				ProjectResourceID: pid,
				UID:               ruleId,
			})
		}
	})
}

func (op *logMeasureRuleOp) Delete(ctx context.Context, projectId string, ruleId uuid.UUID) error {
	return errorFromDecodedResponse1("LogMeasureRule.Delete", func() error {
		if pid, err := strconv.ParseInt(projectId, 10, 64); err != nil {
			return err
		} else {
			return op.client.AlertsProjectsLogMeasureRulesDestroy(ctx, v1.AlertsProjectsLogMeasureRulesDestroyParams{
				ProjectResourceID: pid,
				UID:               ruleId,
			})
		}
	})
}

type LogMeasureRuleListHistoriesParams struct {
	Count    *int
	From     *int
	Open     *bool
	Severity *v1.AlertsProjectsHistoriesListSeverity
	StartsAt *time.Time
}

func (op *logMeasureRuleOp) ListHistories(ctx context.Context, projectId string, params LogMeasureRuleListHistoriesParams) (ret []v1.History, err error) {
	res, err := errorFromDecodedResponse("LogMeasureRule.ListHistories", func() (*v1.PaginatedHistoryList, error) {
		if pid, err := strconv.ParseInt(projectId, 10, 64); err != nil {
			return nil, err
		} else {
			return op.client.AlertsProjectsHistoriesList(ctx, v1.AlertsProjectsHistoriesListParams{
				ProjectResourceID: pid,
				Count:             intoOpt[v1.OptInt](params.Count),
				From:              intoOpt[v1.OptInt](params.From),
				Open:              intoOpt[v1.OptBool](params.Open),
				Severity:          intoOpt[v1.OptAlertsProjectsHistoriesListSeverity](params.Severity),
				StartsAt:          intoOpt[v1.OptDateTime](params.StartsAt),
			})
		}
	})
	if err == nil {
		ret = res.GetResults()
	}
	return
}

func (op *logMeasureRuleOp) ReadHistory(ctx context.Context, projectId string, historyId uuid.UUID) (*v1.History, error) {
	return errorFromDecodedResponse("LogMeasureRule.ReadHistory", func() (*v1.History, error) {
		if pid, err := strconv.ParseInt(projectId, 10, 64); err != nil {
			return nil, err
		} else {
			return op.client.AlertsProjectsHistoriesRetrieve(ctx, v1.AlertsProjectsHistoriesRetrieveParams{
				ProjectResourceID: pid,
				UID:               historyId,
			})
		}
	})
}

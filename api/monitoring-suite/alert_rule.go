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

type AlertRuleAPI interface {
	List(ctx context.Context, projectId string, count *int, from *int) ([]v1.AlertRule, error)
	Create(ctx context.Context, projectId string, params AlertRuleCreateParams) (*v1.AlertRule, error)
	Read(ctx context.Context, projectId string, ruleId uuid.UUID) (*v1.AlertRule, error)
	Update(ctx context.Context, projectId string, ruleId uuid.UUID, params AlertRuleUpdateParams) (*v1.AlertRule, error)
	Delete(ctx context.Context, projectId string, ruleId uuid.UUID) error

	ListHistories(ctx context.Context, projectId string, ruleId uuid.UUID, params AlertRuleListHistoriesParams) ([]v1.History, error)
	ReadHistory(ctx context.Context, projectId string, ruleId uuid.UUID, historyId uuid.UUID) (*v1.History, error)
}

var _ AlertRuleAPI = (*alertRuleOp)(nil)

type alertRuleOp struct {
	client *v1.Client
}

func NewAlertRuleOp(client *v1.Client) AlertRuleAPI {
	return &alertRuleOp{client: client}
}

func (op *alertRuleOp) List(ctx context.Context, projectId string, count *int, from *int) (ret []v1.AlertRule, err error) {
	res, err := errorFromDecodedResponse("AlertRule.List", func() (*v1.PaginatedAlertRuleList, error) {
		id, err := strconv.ParseInt(projectId, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.AlertsProjectsRulesList(ctx, v1.AlertsProjectsRulesListParams{
			ProjectResourceID: id,
			Count:             into.Opt[v1.OptInt](count),
			From:              into.Opt[v1.OptInt](from),
		})
	})
	if err == nil {
		ret = res.GetResults()
	}
	return
}

type AlertRuleCreateParams struct {
	MetricsStorageID          string // mandatory
	Name                      *string
	Query                     string // mandatory
	Format                    *string
	Template                  *string
	EnabledWarning            *bool
	EnabledCritical           *bool
	ThresholdWarning          *string
	ThresholdCritical         *string
	ThresholdDurationWarning  *int64
	ThresholdDurationCritical *int64
}

func (op *alertRuleOp) Create(ctx context.Context, projectId string, p AlertRuleCreateParams) (*v1.AlertRule, error) {
	return errorFromDecodedResponse("AlertRule.Create", func() (*v1.AlertRule, error) {
		intProjectId, err := strconv.ParseInt(projectId, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("projectId: %w", err)
		}
		intStorageId, err := strconv.ParseInt(p.MetricsStorageID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("AlertRuleCreateParams.MetricsStorageID: %w", err)
		}
		return op.client.AlertsProjectsRulesCreate(ctx, &v1.AlertRuleRequest{
			MetricsStorageID:          into.Nil[v1.NilInt64](&intStorageId),
			Name:                      into.Opt[v1.OptString](p.Name),
			Query:                     p.Query,
			Format:                    into.Opt[v1.OptString](p.Format),
			Template:                  into.Opt[v1.OptString](p.Template),
			EnabledWarning:            into.Opt[v1.OptBool](p.EnabledWarning),
			EnabledCritical:           into.Opt[v1.OptBool](p.EnabledCritical),
			ThresholdWarning:          into.OptNil[v1.OptNilString](p.ThresholdWarning),
			ThresholdCritical:         into.OptNil[v1.OptNilString](p.ThresholdCritical),
			ThresholdDurationWarning:  into.Opt[v1.OptInt64](p.ThresholdDurationWarning),
			ThresholdDurationCritical: into.Opt[v1.OptInt64](p.ThresholdDurationCritical),
		}, v1.AlertsProjectsRulesCreateParams{
			ProjectResourceID: intProjectId,
		})
	})
}

func (op *alertRuleOp) Read(ctx context.Context, projectId string, ruleId uuid.UUID) (*v1.AlertRule, error) {
	return errorFromDecodedResponse("AlertRule.Read", func() (*v1.AlertRule, error) {
		intProjectId, err := strconv.ParseInt(projectId, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.AlertsProjectsRulesRetrieve(ctx, v1.AlertsProjectsRulesRetrieveParams{
			ProjectResourceID: intProjectId,
			UID:               ruleId,
		})
	})
}

type AlertRuleUpdateParams struct {
	MetricsStorageID          *string
	Name                      *string
	Query                     *string
	Format                    *string
	Template                  *string
	EnabledWarning            *bool
	EnabledCritical           *bool
	ThresholdWarning          *string
	ThresholdCritical         *string
	ThresholdDurationWarning  *int64
	ThresholdDurationCritical *int64
}

func (op *alertRuleOp) Update(ctx context.Context, projectId string, ruleId uuid.UUID, p AlertRuleUpdateParams) (*v1.AlertRule, error) {
	return errorFromDecodedResponse("AlertRule.Update", func() (*v1.AlertRule, error) {
		intProjectId, err := strconv.ParseInt(projectId, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("projectId: %w", err)
		}
		storageId, err := into.FromStringPtr[v1.OptNilInt64, int64](p.MetricsStorageID)
		if err != nil {
			return nil, fmt.Errorf("AlertRuleUpdateParams.MetricsStorageID: %w", err)
		}
		return op.client.AlertsProjectsRulesPartialUpdate(ctx, v1.NewOptPatchedAlertRuleRequest(v1.PatchedAlertRuleRequest{
			MetricsStorageID:          storageId,
			Name:                      into.Opt[v1.OptString](p.Name),
			Query:                     into.Opt[v1.OptString](p.Query),
			Format:                    into.Opt[v1.OptString](p.Format),
			Template:                  into.Opt[v1.OptString](p.Template),
			EnabledWarning:            into.Opt[v1.OptBool](p.EnabledWarning),
			EnabledCritical:           into.Opt[v1.OptBool](p.EnabledCritical),
			ThresholdWarning:          into.OptNil[v1.OptNilString](p.ThresholdWarning),
			ThresholdCritical:         into.OptNil[v1.OptNilString](p.ThresholdCritical),
			ThresholdDurationWarning:  into.Opt[v1.OptInt64](p.ThresholdDurationWarning),
			ThresholdDurationCritical: into.Opt[v1.OptInt64](p.ThresholdDurationCritical),
		}), v1.AlertsProjectsRulesPartialUpdateParams{
			ProjectResourceID: intProjectId,
			UID:               ruleId,
		})
	})
}

func (op *alertRuleOp) Delete(ctx context.Context, projectId string, ruleId uuid.UUID) error {
	return errorFromDecodedResponse1("AlertRule.Delete", func() error {
		intProjectId, err := strconv.ParseInt(projectId, 10, 64)
		if err != nil {
			return err
		}
		return op.client.AlertsProjectsRulesDestroy(ctx, v1.AlertsProjectsRulesDestroyParams{
			ProjectResourceID: intProjectId,
			UID:               ruleId,
		})
	})
}

type AlertRuleListHistoriesParams struct {
	Count    *int
	From     *int
	Open     *bool
	Severity *v1.AlertsProjectsRulesHistoriesListSeverity
	StartsAt *time.Time
}

func (op *alertRuleOp) ListHistories(ctx context.Context, projectId string, ruleId uuid.UUID, p AlertRuleListHistoriesParams) (ret []v1.History, err error) {
	res, err := errorFromDecodedResponse("AlertRule.ListHistories", func() (*v1.PaginatedHistoryList, error) {
		intProjectId, err := strconv.ParseInt(projectId, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.AlertsProjectsRulesHistoriesList(ctx, v1.AlertsProjectsRulesHistoriesListParams{
			ProjectResourceID: intProjectId,
			RuleUID:           ruleId,
			Count:             into.Opt[v1.OptInt](p.Count),
			From:              into.Opt[v1.OptInt](p.From),
			Open:              into.Opt[v1.OptBool](p.Open),
			Severity:          into.Opt[v1.OptAlertsProjectsRulesHistoriesListSeverity](p.Severity),
			StartsAt:          into.Opt[v1.OptDateTime](p.StartsAt),
		})
	})
	if err == nil {
		ret = res.GetResults()
	}
	return
}

func (op *alertRuleOp) ReadHistory(ctx context.Context, projectId string, ruleId uuid.UUID, historyId uuid.UUID) (*v1.History, error) {
	return errorFromDecodedResponse("AlertRule.ReadHistory", func() (*v1.History, error) {
		intProjectId, err := strconv.ParseInt(projectId, 10, 64)
		if err != nil {
			return nil, err
		}
		return op.client.AlertsProjectsRulesHistoriesRetrieve(ctx, v1.AlertsProjectsRulesHistoriesRetrieveParams{
			ProjectResourceID: intProjectId,
			RuleUID:           ruleId,
			UID:               historyId,
		})
	})
}

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

func (op *alertRuleOp) List(ctx context.Context, projectId string, count *int, from *int) ([]v1.AlertRule, error) {
	id, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewAPIError("AlertRule.List", 0, err)
	}
	params := v1.AlertsProjectsRulesListParams{
		ProjectResourceID: id,
		Count:             intoOpt[v1.OptInt](count),
		From:              intoOpt[v1.OptInt](from),
	}
	result, err := op.client.AlertsProjectsRulesList(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("AlertRule.List", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("AlertRule.List", e.StatusCode, errors.Wrap(err, "project not found"))
		default:
			return nil, NewAPIError("AlertRule.List", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("AlertRule.List", 0, err)
	} else {
		return result.GetResults(), nil
	}
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
	intProjectId, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewAPIError("AlertRule.Create", 0, errors.Wrap(err, "invalid ProjectID"))
	}
	intStorageId, err := strconv.ParseInt(p.MetricsStorageID, 10, 64)
	if err != nil {
		return nil, NewAPIError("AlertRule.Create", 0, errors.Wrap(err, "invalid MetricsStorageID"))
	}
	params := &v1.AlertRule{
		MetricsStorageID:          intoNil[v1.NilInt64](&intStorageId),
		Name:                      intoOpt[v1.OptString](p.Name),
		Query:                     p.Query,
		Format:                    intoOpt[v1.OptString](p.Format),
		Template:                  intoOpt[v1.OptString](p.Template),
		EnabledWarning:            intoOpt[v1.OptBool](p.EnabledWarning),
		EnabledCritical:           intoOpt[v1.OptBool](p.EnabledCritical),
		ThresholdWarning:          intoOptNil[v1.OptNilString](p.ThresholdWarning),
		ThresholdCritical:         intoOptNil[v1.OptNilString](p.ThresholdCritical),
		ThresholdDurationWarning:  intoOpt[v1.OptInt64](p.ThresholdDurationWarning),
		ThresholdDurationCritical: intoOpt[v1.OptInt64](p.ThresholdDurationCritical),
	}
	query := v1.AlertsProjectsRulesCreateParams{
		ProjectResourceID: intProjectId,
	}
	result, err := op.client.AlertsProjectsRulesCreate(ctx, params, query)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("AlertRule.Create", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("AlertRule.Create", e.StatusCode, errors.Wrap(err, "project not found"))
		case http.StatusBadRequest:
			return nil, NewAPIError("AlertRule.Create", e.StatusCode, errors.Wrap(err, "invalid parameter, or no space left for a new alert rule"))
		default:
			return nil, NewAPIError("AlertRule.Create", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("AlertRule.Create", 0, err)
	} else {
		return result, nil
	}
}

func (op *alertRuleOp) Read(ctx context.Context, projectId string, ruleId uuid.UUID) (*v1.AlertRule, error) {
	intProjectId, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewAPIError("AlertRule.Read", 0, err)
	}
	query := v1.AlertsProjectsRulesRetrieveParams{
		ProjectResourceID: intProjectId,
		UID:               ruleId,
	}
	result, err := op.client.AlertsProjectsRulesRetrieve(ctx, query)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("AlertRule.Read", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("AlertRule.Read", e.StatusCode, errors.Wrap(err, "alert rule not found"))
		default:
			return nil, NewAPIError("AlertRule.Read", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("AlertRule.Read", 0, err)
	}
	return result, nil
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
	intProjectId, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewAPIError("AlertRule.Update", 0, err)
	}
	query := v1.AlertsProjectsRulesPartialUpdateParams{
		ProjectResourceID: intProjectId,
		UID:               ruleId,
	}
	storageId, err := fromStringPtr[v1.OptNilInt64, int64](p.MetricsStorageID)
	if err != nil {
		return nil, NewAPIError("AlertRule.Update", 0, errors.Wrap(err, "invalid MetricsStorageID"))
	}
	params := v1.NewOptPatchedAlertRule(v1.PatchedAlertRule{
		MetricsStorageID:          storageId,
		Name:                      intoOpt[v1.OptString](p.Name),
		Query:                     intoOpt[v1.OptString](p.Query),
		Format:                    intoOpt[v1.OptString](p.Format),
		Template:                  intoOpt[v1.OptString](p.Template),
		EnabledWarning:            intoOpt[v1.OptBool](p.EnabledWarning),
		EnabledCritical:           intoOpt[v1.OptBool](p.EnabledCritical),
		ThresholdWarning:          intoOptNil[v1.OptNilString](p.ThresholdWarning),
		ThresholdCritical:         intoOptNil[v1.OptNilString](p.ThresholdCritical),
		ThresholdDurationWarning:  intoOpt[v1.OptInt64](p.ThresholdDurationWarning),
		ThresholdDurationCritical: intoOpt[v1.OptInt64](p.ThresholdDurationCritical),
	})
	result, err := op.client.AlertsProjectsRulesPartialUpdate(ctx, params, query)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("AlertRule.Update", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("AlertRule.Update", e.StatusCode, errors.Wrap(err, "alert rule not found"))
		case http.StatusBadRequest:
			return nil, NewAPIError("AlertRule.Update", e.StatusCode, errors.Wrap(err, "invalid parameter"))
		default:
			return nil, NewAPIError("AlertRule.Update", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("AlertRule.Update", 0, err)
	}
	return result, nil
}

func (op *alertRuleOp) Delete(ctx context.Context, projectId string, ruleId uuid.UUID) error {
	intProjectId, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return NewAPIError("AlertRule.Delete", 0, err)
	}
	query := v1.AlertsProjectsRulesDestroyParams{
		ProjectResourceID: intProjectId,
		UID:               ruleId,
	}
	err = op.client.AlertsProjectsRulesDestroy(ctx, query)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return NewAPIError("AlertRule.Delete", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return NewAPIError("AlertRule.Delete", e.StatusCode, errors.Wrap(err, "alert rule not found"))
		case http.StatusBadRequest:
			return NewAPIError("AlertRule.Delete", e.StatusCode, errors.Wrap(err, "the request resource is not eligible for deletion"))
		default:
			return NewAPIError("AlertRule.Delete", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return NewAPIError("AlertRule.Delete", 0, err)
	}
	return nil
}

type AlertRuleListHistoriesParams struct {
	Count    *int
	From     *int
	Open     *bool
	Severity *v1.AlertsProjectsRulesHistoriesListSeverity
	StartsAt *time.Time
}

func (op *alertRuleOp) ListHistories(ctx context.Context, projectId string, ruleId uuid.UUID, p AlertRuleListHistoriesParams) ([]v1.History, error) {
	intProjectId, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewAPIError("AlertRule.ListHistories", 0, err)
	}
	params := v1.AlertsProjectsRulesHistoriesListParams{
		ProjectResourceID: intProjectId,
		RuleUID:           ruleId,
		Count:             intoOpt[v1.OptInt](p.Count),
		From:              intoOpt[v1.OptInt](p.From),
		Open:              intoOpt[v1.OptBool](p.Open),
		Severity:          intoOpt[v1.OptAlertsProjectsRulesHistoriesListSeverity](p.Severity),
		StartsAt:          intoOpt[v1.OptDateTime](p.StartsAt),
	}
	result, err := op.client.AlertsProjectsRulesHistoriesList(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("AlertRule.ListHistories", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("AlertRule.ListHistories", e.StatusCode, errors.Wrap(err, "project or rule not found"))
		default:
			return nil, NewAPIError("AlertRule.ListHistories", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("AlertRule.ListHistories", 0, err)
	} else {
		return result.GetResults(), nil
	}
}

func (op *alertRuleOp) ReadHistory(ctx context.Context, projectId string, ruleId uuid.UUID, historyId uuid.UUID) (*v1.History, error) {
	intProjectId, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewAPIError("AlertRule.ReadHistory", 0, err)
	}
	query := v1.AlertsProjectsRulesHistoriesRetrieveParams{
		ProjectResourceID: intProjectId,
		RuleUID:           ruleId,
		UID:               historyId,
	}
	result, err := op.client.AlertsProjectsRulesHistoriesRetrieve(ctx, query)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("AlertRule.ReadHistory", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("AlertRule.ReadHistory", e.StatusCode, errors.Wrap(err, "history not found"))
		default:
			return nil, NewAPIError("AlertRule.ReadHistory", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("AlertRule.ReadHistory", 0, err)
	}
	return result, nil
}

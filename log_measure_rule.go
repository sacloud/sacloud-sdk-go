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

func (op *logMeasureRuleOp) List(ctx context.Context, projectId string, count *int, from *int) ([]v1.LogMeasureRule, error) {
	pid, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewAPIError("LogMeasureRule.List", 0, errors.Wrap(err, "invalid projectID"))
	}
	params := v1.AlertsProjectsLogMeasureRulesListParams{
		ProjectResourceID: pid,
		Count:             intoOpt[v1.OptInt](count),
		From:              intoOpt[v1.OptInt](from),
	}
	result, err := op.client.AlertsProjectsLogMeasureRulesList(ctx, params)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("LogMeasureRule.List", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("LogMeasureRule.List", e.StatusCode, errors.Wrap(err, "project not found"))
		default:
			return nil, NewAPIError("LogMeasureRule.List", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("LogMeasureRule.List", 0, err)
	} else {
		return result.GetResults(), nil
	}
}

type LogMeasureRuleCreateParams struct {
	LogStorageID     string
	MetricsStorageID string
	Name             *string
	Description      *string
	Rule             v1.LogMeasureRuleModel
}

func (op *logMeasureRuleOp) Create(ctx context.Context, projectId string, p LogMeasureRuleCreateParams) (*v1.LogMeasureRule, error) {
	pid, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewAPIError("LogMeasureRule.Create", 0, errors.Wrap(err, "invalid projectID"))
	}
	lid, err := strconv.ParseInt(p.LogStorageID, 10, 64)
	if err != nil {
		return nil, NewAPIError("LogMeasureRule.Create", 0, errors.Wrap(err, "invalid LogStorageID"))
	}
	mid, err := strconv.ParseInt(p.MetricsStorageID, 10, 64)
	if err != nil {
		return nil, NewAPIError("LogMeasureRule.Create", 0, errors.Wrap(err, "invalid MetricsStorageID"))
	}
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
	params.LogStorage.SetTags([]string{})
	params.MetricsStorage.SetTags([]string{})

	query := v1.AlertsProjectsLogMeasureRulesCreateParams{
		ProjectResourceID: pid,
	}
	result, err := op.client.AlertsProjectsLogMeasureRulesCreate(ctx, params, query)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("LogMeasureRule.Create", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("LogMeasureRule.Create", e.StatusCode, errors.Wrap(err, "project not found"))
		case http.StatusBadRequest:
			return nil, NewAPIError("LogMeasureRule.Create", e.StatusCode, errors.Wrap(err, "invalid parameter"))
		default:
			return nil, NewAPIError("LogMeasureRule.Create", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("LogMeasureRule.Create", 0, err)
	} else {
		return result, nil
	}
}

func (op *logMeasureRuleOp) Read(ctx context.Context, projectId string, ruleId uuid.UUID) (*v1.LogMeasureRule, error) {
	pid, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewAPIError("LogMeasureRule.Read", 0, errors.Wrap(err, "invalid projectID"))
	}
	query := v1.AlertsProjectsLogMeasureRulesRetrieveParams{
		ProjectResourceID: pid,
		UID:               ruleId,
	}
	result, err := op.client.AlertsProjectsLogMeasureRulesRetrieve(ctx, query)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("LogMeasureRule.Read", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("LogMeasureRule.Read", e.StatusCode, errors.Wrap(err, "log measure rule not found"))
		default:
			return nil, NewAPIError("LogMeasureRule.Read", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("LogMeasureRule.Read", 0, err)
	}
	return result, nil
}

type LogMeasureRuleUpdateParams struct {
	LogStorageID     *string
	MetricsStorageID *string
	Name             *string
	Description      *string
	Rule             *v1.LogMeasureRuleModel
}

func (op *logMeasureRuleOp) Update(ctx context.Context, projectId string, ruleId uuid.UUID, p LogMeasureRuleUpdateParams) (*v1.LogMeasureRule, error) {
	pid, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewAPIError("LogMeasureRule.Update", 0, errors.Wrap(err, "invalid projectID"))
	}
	query := v1.AlertsProjectsLogMeasureRulesPartialUpdateParams{
		ProjectResourceID: pid,
		UID:               ruleId,
	}
	lid, err := fromStringPtr[v1.OptNilInt64, int64](p.LogStorageID)
	if err != nil {
		return nil, NewAPIError("LogMeasureRule.Update", 0, errors.Wrap(err, "invalid LogStorageID"))
	}
	mid, err := fromStringPtr[v1.OptNilInt64, int64](p.MetricsStorageID)
	if err != nil {
		return nil, NewAPIError("LogMeasureRule.Update", 0, errors.Wrap(err, "invalid MetricsStorageID"))
	}
	params := v1.NewOptPatchedLogMeasureRule(v1.PatchedLogMeasureRule{
		LogStorageID:     lid,
		MetricsStorageID: mid,
		Name:             intoOpt[v1.OptString](p.Name),
		Description:      intoOpt[v1.OptString](p.Description),
		Rule:             intoOpt[v1.OptLogMeasureRuleModel](p.Rule),
	})
	result, err := op.client.AlertsProjectsLogMeasureRulesPartialUpdate(ctx, params, query)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("LogMeasureRule.Update", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("LogMeasureRule.Update", e.StatusCode, errors.Wrap(err, "log measure rule not found"))
		case http.StatusBadRequest:
			return nil, NewAPIError("LogMeasureRule.Update", e.StatusCode, errors.Wrap(err, "invalid parameter"))
		default:
			return nil, NewAPIError("LogMeasureRule.Update", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("LogMeasureRule.Update", 0, err)
	}
	return result, nil
}

func (op *logMeasureRuleOp) Delete(ctx context.Context, projectId string, ruleId uuid.UUID) error {
	pid, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return NewAPIError("LogMeasureRule.Delete", 0, errors.Wrap(err, "invalid projectID"))
	}
	query := v1.AlertsProjectsLogMeasureRulesDestroyParams{
		ProjectResourceID: pid,
		UID:               ruleId,
	}
	err = op.client.AlertsProjectsLogMeasureRulesDestroy(ctx, query)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return NewAPIError("LogMeasureRule.Delete", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return NewAPIError("LogMeasureRule.Delete", e.StatusCode, errors.Wrap(err, "log measure rule not found"))
		case http.StatusBadRequest:
			return NewAPIError("LogMeasureRule.Delete", e.StatusCode, errors.Wrap(err, "the request resource is not eligible for deletion"))
		default:
			return NewAPIError("LogMeasureRule.Delete", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return NewAPIError("LogMeasureRule.Delete", 0, err)
	} else {
		return nil
	}
}

type LogMeasureRuleListHistoriesParams struct {
	Count    *int
	From     *int
	Open     *bool
	Severity *v1.AlertsProjectsHistoriesListSeverity
	StartsAt *time.Time
}

func (op *logMeasureRuleOp) ListHistories(ctx context.Context, projectId string, params LogMeasureRuleListHistoriesParams) ([]v1.History, error) {
	pid, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewAPIError("LogMeasureRule.ListHistories", 0, errors.Wrap(err, "invalid projectID"))
	}
	apiParams := v1.AlertsProjectsHistoriesListParams{
		ProjectResourceID: pid,
		Count:             intoOpt[v1.OptInt](params.Count),
		From:              intoOpt[v1.OptInt](params.From),
		Open:              intoOpt[v1.OptBool](params.Open),
		Severity:          intoOpt[v1.OptAlertsProjectsHistoriesListSeverity](params.Severity),
		StartsAt:          intoOpt[v1.OptDateTime](params.StartsAt),
	}
	result, err := op.client.AlertsProjectsHistoriesList(ctx, apiParams)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("LogMeasureRule.ListHistories", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("LogMeasureRule.ListHistories", e.StatusCode, errors.Wrap(err, "project not found"))
		default:
			return nil, NewAPIError("LogMeasureRule.ListHistories", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("LogMeasureRule.ListHistories", 0, err)
	} else {
		return result.GetResults(), nil
	}
}

func (op *logMeasureRuleOp) ReadHistory(ctx context.Context, projectId string, historyId uuid.UUID) (*v1.History, error) {
	pid, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		return nil, NewAPIError("LogMeasureRule.ReadHistory", 0, errors.Wrap(err, "invalid projectID"))
	}
	query := v1.AlertsProjectsHistoriesRetrieveParams{
		ProjectResourceID: pid,
		UID:               historyId,
	}
	result, err := op.client.AlertsProjectsHistoriesRetrieve(ctx, query)
	if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		switch e.StatusCode {
		case http.StatusForbidden:
			return nil, NewAPIError("LogMeasureRule.ReadHistory", e.StatusCode, errors.Wrap(err, "insufficient permissions"))
		case http.StatusNotFound:
			return nil, NewAPIError("LogMeasureRule.ReadHistory", e.StatusCode, errors.Wrap(err, "history not found"))
		default:
			return nil, NewAPIError("LogMeasureRule.ReadHistory", e.StatusCode, errors.Wrap(err, "internal server error"))
		}
	} else if err != nil {
		return nil, NewAPIError("LogMeasureRule.ReadHistory", 0, err)
	} else {
		return result, nil
	}
}

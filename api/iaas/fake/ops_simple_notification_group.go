// Copyright 2022-2025 The sacloud/iaas-api-go Authors
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

package fake

import (
	"context"
	"time"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
)

// Find is fake implementation
func (o *SimpleNotificationGroupOp) Find(ctx context.Context, conditions *iaas.FindCondition) (*iaas.SimpleNotificationGroupFindResult, error) {
	results, _ := find(o.key, iaas.APIDefaultZone, conditions)
	var values []*iaas.SimpleNotificationGroup
	for _, res := range results {
		dest := &iaas.SimpleNotificationGroup{}
		copySameNameField(res, dest)
		values = append(values, dest)
	}
	return &iaas.SimpleNotificationGroupFindResult{
		Total:                    len(results),
		Count:                    len(results),
		From:                     0,
		SimpleNotificationGroups: values,
	}, nil
}

// Create is fake implementation
func (o *SimpleNotificationGroupOp) Create(ctx context.Context, param *iaas.SimpleNotificationGroupCreateRequest) (*iaas.SimpleNotificationGroup, error) {
	result := &iaas.SimpleNotificationGroup{}
	copySameNameField(param, result)
	fill(result, fillID, fillCreatedAt)

	result.Availability = types.Availabilities.Available
	putSimpleNotificationGroup(iaas.APIDefaultZone, result)
	return result, nil
}

// Read is fake implementation
func (o *SimpleNotificationGroupOp) Read(ctx context.Context, id types.ID) (*iaas.SimpleNotificationGroup, error) {
	value := getSimpleNotificationGroupByID(iaas.APIDefaultZone, id)
	if value == nil {
		return nil, newErrorNotFound(o.key, id)
	}
	dest := &iaas.SimpleNotificationGroup{}
	copySameNameField(value, dest)
	return dest, nil
}

// Update is fake implementation
func (o *SimpleNotificationGroupOp) Update(ctx context.Context, id types.ID, param *iaas.SimpleNotificationGroupUpdateRequest) (*iaas.SimpleNotificationGroup, error) {
	value, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}
	copySameNameField(param, value)
	fill(value, fillModifiedAt)

	return value, nil
}

// UpdateSettings is fake implementation
func (o *SimpleNotificationGroupOp) UpdateSettings(ctx context.Context, id types.ID, param *iaas.SimpleNotificationGroupUpdateSettingsRequest) (*iaas.SimpleNotificationGroup, error) {
	value, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}
	copySameNameField(param, value)
	fill(value, fillModifiedAt)

	return value, nil
}

// Delete is fake implementation
func (o *SimpleNotificationGroupOp) Delete(ctx context.Context, id types.ID) error {
	_, err := o.Read(ctx, id)
	if err != nil {
		return err
	}

	ds().Delete(o.key, iaas.APIDefaultZone, id)
	return nil
}

// PostMessage is fake implementation
func (o *SimpleNotificationGroupOp) PostMessage(ctx context.Context, id types.ID, message string) error {
	return nil
}

func (o *SimpleNotificationGroupOp) History(ctx context.Context) (*iaas.SimpleNotificationHistories, error) {
	return &iaas.SimpleNotificationHistories{
		NotificationHistories: []*iaas.SimpleNotificationHistory{
			{
				RequestID:  "11111",
				SourceID:   "1",
				ReceivedAt: time.Now(),
				Message: &iaas.SimpleNotificationHistoryMessage{
					Body:      "body",
					Color:     "color",
					ColorCode: "#000000",
					IconURL:   "",
					ImageURL:  "",
					Title:     "title",
				},
				Statuses: []*iaas.SimpleNotificationHistoryStatus{
					{
						ID:                    "1",
						Status:                1,
						ErrorInfo:             "error",
						NotificationRequestID: "11111",
						GroupID:               "123456789012",
						CreatedAt:             time.Now(),
						UpdatedAt:             time.Now(),
					},
				},
			},
		},
	}, nil
}

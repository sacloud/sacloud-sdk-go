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
func (o *SimpleNotificationDestinationOp) Find(ctx context.Context, conditions *iaas.FindCondition) (*iaas.SimpleNotificationDestinationFindResult, error) {
	results, _ := find(o.key, iaas.APIDefaultZone, conditions)
	var values []*iaas.SimpleNotificationDestination
	for _, res := range results {
		dest := &iaas.SimpleNotificationDestination{}
		copySameNameField(res, dest)
		values = append(values, dest)
	}
	return &iaas.SimpleNotificationDestinationFindResult{
		Total:                          len(results),
		Count:                          len(results),
		From:                           0,
		SimpleNotificationDestinations: values,
	}, nil
}

// Create is fake implementation
func (o *SimpleNotificationDestinationOp) Create(ctx context.Context, param *iaas.SimpleNotificationDestinationCreateRequest) (*iaas.SimpleNotificationDestination, error) {
	result := &iaas.SimpleNotificationDestination{}
	copySameNameField(param, result)
	fill(result, fillID, fillCreatedAt)

	result.Availability = types.Availabilities.Available
	putSimpleNotificationDestination(iaas.APIDefaultZone, result)
	return result, nil
}

// Read is fake implementation
func (o *SimpleNotificationDestinationOp) Read(ctx context.Context, id types.ID) (*iaas.SimpleNotificationDestination, error) {
	value := getSimpleNotificationDestinationByID(iaas.APIDefaultZone, id)
	if value == nil {
		return nil, newErrorNotFound(o.key, id)
	}
	dest := &iaas.SimpleNotificationDestination{}
	copySameNameField(value, dest)
	return dest, nil
}

// Update is fake implementation
func (o *SimpleNotificationDestinationOp) Update(ctx context.Context, id types.ID, param *iaas.SimpleNotificationDestinationUpdateRequest) (*iaas.SimpleNotificationDestination, error) {
	value, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}
	copySameNameField(param, value)
	fill(value, fillModifiedAt)

	return value, nil
}

// UpdateSettings is fake implementation
func (o *SimpleNotificationDestinationOp) UpdateSettings(ctx context.Context, id types.ID, param *iaas.SimpleNotificationDestinationUpdateSettingsRequest) (*iaas.SimpleNotificationDestination, error) {
	value, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}
	copySameNameField(param, value)
	fill(value, fillModifiedAt)
	return value, err
}

// Delete is fake implementation
func (o *SimpleNotificationDestinationOp) Delete(ctx context.Context, id types.ID) error {
	_, err := o.Read(ctx, id)
	if err != nil {
		return err
	}

	ds().Delete(o.key, iaas.APIDefaultZone, id)
	return nil
}

// Status is fake implementation
func (o *SimpleNotificationDestinationOp) Status(ctx context.Context, id types.ID) (*iaas.SimpleNotificationDestinationStatus, error) {
	_, err := o.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	return &iaas.SimpleNotificationDestinationStatus{
		Disabled:   false,
		ModifiedAt: time.Now(),
	}, nil
}

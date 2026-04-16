// Copyright 2022-2025 The sacloud/iaas-service-go Authors
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

package autoscale

import (
	"testing"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/packages-go/pointer"
	"github.com/stretchr/testify/require"
)

func TestUpdateRequest_ToRequestParameter(t *testing.T) {
	type fields struct {
		ID                     types.ID
		Name                   *string
		Description            *string
		Tags                   *types.Tags
		IconID                 *types.ID
		Zones                  *[]string
		Config                 *string
		CPUThresholdScaling    *UpdateCPUThresholdScaling
		RouterThresholdScaling *UpdateRouterThresholdScaling
		ScheduleScaling        *[]*UpdateScheduleScaling
		Disabled               *bool
		SettingsHash           string
	}
	type args struct {
		current *iaas.AutoScale
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *iaas.AutoScaleUpdateRequest
		wantErr bool
	}{
		{
			name: "minimum",
			fields: fields{
				ID: 1,
				CPUThresholdScaling: &UpdateCPUThresholdScaling{
					Up:   pointer.NewInt(81),
					Down: pointer.NewInt(21),
				},
			},
			args: args{
				current: &iaas.AutoScale{
					ID:   1,
					Name: "name",
					CPUThresholdScaling: &iaas.AutoScaleCPUThresholdScaling{
						ServerPrefix: "foobar",
						Up:           80,
						Down:         20,
					},
				},
			},
			want: &iaas.AutoScaleUpdateRequest{
				Name: "name",
				CPUThresholdScaling: &iaas.AutoScaleCPUThresholdScaling{
					ServerPrefix: "foobar",
					Up:           81,
					Down:         21,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &UpdateRequest{
				ID:                     tt.fields.ID,
				Name:                   tt.fields.Name,
				Description:            tt.fields.Description,
				Tags:                   tt.fields.Tags,
				IconID:                 tt.fields.IconID,
				Zones:                  tt.fields.Zones,
				Config:                 tt.fields.Config,
				CPUThresholdScaling:    tt.fields.CPUThresholdScaling,
				RouterThresholdScaling: tt.fields.RouterThresholdScaling,
				ScheduleScaling:        tt.fields.ScheduleScaling,
				Disabled:               tt.fields.Disabled,
				SettingsHash:           tt.fields.SettingsHash,
			}
			got, err := req.ToRequestParameter(tt.args.current)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToRequestParameter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

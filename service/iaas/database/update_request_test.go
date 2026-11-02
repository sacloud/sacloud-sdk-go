// Copyright 2026- The sacloud/iaas-service-go Authors
// SPDX-License-Identifier: Apache-2.0

package database

import (
	"testing"

	"github.com/sacloud/sacloud-sdk-go/api/iaas/types"
)

func TestUpdateRequestValidateBackupv2Weekdays(t *testing.T) {
	tests := []struct {
		name    string
		request *UpdateRequest
		wantErr string
	}{
		{
			name:    "weekdays omitted",
			request: &UpdateRequest{Zone: "is1b", ID: 123},
		},
		{
			name: "required when backupv2 is enabled",
			request: &UpdateRequest{
				Zone:            "is1b",
				ID:              123,
				EnableBackupv2:  new(true),
				Backupv2Connect: new("nfs://192.0.2.1/export"),
			},
			wantErr: "Backupv2Weekdays: required when EnableBackupv2 is set",
		},
		{
			name: "seven weekdays",
			request: &UpdateRequest{
				Zone: "is1b",
				ID:   123,
				Backupv2Weekdays: &[]types.EDayOfTheWeek{
					types.DaysOfTheWeek.Sunday,
					types.DaysOfTheWeek.Monday,
					types.DaysOfTheWeek.Tuesday,
					types.DaysOfTheWeek.Wednesday,
					types.DaysOfTheWeek.Thursday,
					types.DaysOfTheWeek.Friday,
					types.DaysOfTheWeek.Saturday,
				},
			},
		},
		{
			name: "more than seven weekdays",
			request: &UpdateRequest{
				Zone: "is1b",
				ID:   123,
				Backupv2Weekdays: &[]types.EDayOfTheWeek{
					types.DaysOfTheWeek.Sunday,
					types.DaysOfTheWeek.Monday,
					types.DaysOfTheWeek.Tuesday,
					types.DaysOfTheWeek.Wednesday,
					types.DaysOfTheWeek.Thursday,
					types.DaysOfTheWeek.Friday,
					types.DaysOfTheWeek.Saturday,
					types.DaysOfTheWeek.Sunday,
				},
			},
			wantErr: "Backupv2Weekdays: must contain at most 7 weekdays",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("Validate() error = %v, want nil", err)
				}
				return
			}
			if err == nil || err.Error() != tt.wantErr {
				t.Fatalf("Validate() error = %v, want %q", err, tt.wantErr)
			}
		})
	}
}

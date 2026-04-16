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

package sim

import (
	"context"
	"testing"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/cleanup"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/packages-go/pointer"
	"github.com/stretchr/testify/require"
)

func TestSIMService_convertUpdateRequest(t *testing.T) {
	ctx := context.Background()
	caller := testutil.SingletonAPICaller()
	name := testutil.ResourceName("sim-service")

	// setup
	simOp := iaas.NewSIMOp(caller)
	sim, err := New(caller).CreateWithContext(ctx, &CreateRequest{
		Name:        name,
		Description: "desc",
		Tags:        types.Tags{"tag1", "tag2"},
		ICCID:       "aaaaaaaa",
		PassCode:    "bbbbbbbb",
		Activate:    true,
		IMEI:        "cccccccc",
		Carriers: []*iaas.SIMNetworkOperatorConfig{
			{
				Allow: true,
				Name:  types.SIMOperators.KDDI.String(),
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		cleanup.DeleteSIM(ctx, simOp, sim.ID) //nolint
	}()

	// test
	cases := []struct {
		in     *UpdateRequest
		expect *ApplyRequest
	}{
		{
			in: &UpdateRequest{
				ID:       sim.ID,
				Name:     pointer.NewString(name + "-upd"),
				Activate: pointer.NewBool(false),
				IMEI:     pointer.NewString(""),
				Carriers: &[]*iaas.SIMNetworkOperatorConfig{
					{Allow: true, Name: types.SIMOperators.SoftBank.String()},
				},
			},
			expect: &ApplyRequest{
				ID:          sim.ID,
				Name:        name + "-upd",
				Description: sim.Description,
				Tags:        sim.Tags,
				IconID:      sim.IconID,
				ICCID:       sim.ICCID,
				PassCode:    "",
				Activate:    false,
				IMEI:        "",
				Carriers: []*iaas.SIMNetworkOperatorConfig{
					{
						Allow: true,
						Name:  types.SIMOperators.SoftBank.String(),
					},
				},
			},
		},
	}

	for _, tc := range cases {
		req, err := tc.in.ApplyRequest(ctx, caller)
		require.NoError(t, err)
		require.EqualValues(t, tc.expect, req)
	}
}

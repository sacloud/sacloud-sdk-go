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

package setup

import (
	"context"
	"testing"
	"time"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/accessor"
	"github.com/sacloud/iaas-api-go/helper/power"
	"github.com/sacloud/iaas-api-go/helper/query"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
)

func TestRetryableSetup(t *testing.T) {
	var switchID types.ID
	testZone := testutil.TestZone()

	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		Parallel:           true,
		SetupAPICallerFunc: testutil.SingletonAPICaller,

		Setup: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			switchOp := iaas.NewSwitchOp(caller)
			sw, err := switchOp.Create(ctx, testZone, &iaas.SwitchCreateRequest{Name: "libsacloud-switch-for-util-setup"})
			if err != nil {
				return err
			}
			switchID = sw.ID
			return nil
		},

		Create: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
				nfsOp := iaas.NewNFSOp(caller)
				nfsSetup := &RetryableSetup{
					Create: func(ctx context.Context, zone string) (accessor.ID, error) {
						nfsPlanID, err := query.FindNFSPlanID(ctx, iaas.NewNoteOp(caller), types.NFSPlans.HDD, types.NFSHDDSizes.Size100GB)
						if err != nil {
							return nil, err
						}
						return nfsOp.Create(ctx, zone, &iaas.NFSCreateRequest{
							Name:           "libsacloud-nfs-for-util-setup",
							SwitchID:       switchID,
							PlanID:         nfsPlanID,
							IPAddresses:    []string{"192.168.0.11"},
							NetworkMaskLen: 24,
							DefaultRoute:   "192.168.0.1",
						})
					},
					Read: func(ctx context.Context, zone string, id types.ID) (interface{}, error) {
						return nfsOp.Read(ctx, zone, id)
					},
					Delete: func(ctx context.Context, zone string, id types.ID) error {
						return nfsOp.Delete(ctx, zone, id)
					},
					IsWaitForCopy: true,
					IsWaitForUp:   true,
					Options: &Options{
						RetryCount: 3,
					},
				}
				if !testutil.IsAccTest() {
					nfsSetup.Options.ProvisioningRetryInterval = time.Millisecond
					nfsSetup.Options.DeleteRetryInterval = time.Millisecond
					nfsSetup.Options.PollingInterval = time.Millisecond
				}

				return nfsSetup.Setup(ctx, testZone)
			},
		},
		Read: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
				nfsOp := iaas.NewNFSOp(caller)
				return nfsOp.Read(ctx, testZone, ctx.ID)
			},
			CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, i interface{}) error {
				nfs := i.(*iaas.NFS)
				return testutil.DoAsserts(
					testutil.AssertEqualFunc(t, types.Availabilities.Available, nfs.Availability, "NFS.Availability"),
				)
			},
		},

		Shutdown: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			return power.ShutdownNFS(ctx, iaas.NewNFSOp(caller), testZone, ctx.ID, true)
		},

		Delete: &testutil.CRUDTestDeleteFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
				nfsOp := iaas.NewNFSOp(caller)
				if err := nfsOp.Delete(ctx, testZone, ctx.ID); err != nil {
					return err
				}

				switchOp := iaas.NewSwitchOp(caller)
				if err := switchOp.Delete(ctx, testZone, switchID); err != nil {
					return err
				}
				return nil
			},
		},
	})
}

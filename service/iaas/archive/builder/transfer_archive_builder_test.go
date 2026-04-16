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

package builder

import (
	"testing"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/query"
	"github.com/sacloud/iaas-api-go/ostype"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
)

func TestTransferArchiveBuilder_Build(t *testing.T) {
	zoneFrom := "is1a"
	zoneTo := "is1b"
	var sourceArchive *iaas.Archive

	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		SetupAPICallerFunc: testutil.SingletonAPICaller,
		Parallel:           true,
		IgnoreStartupWait:  true,
		Setup: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			archiveOp := iaas.NewArchiveOp(caller)
			source, err := query.FindArchiveByOSType(ctx, archiveOp, zoneFrom, ostype.Ubuntu)
			if err != nil {
				return err
			}

			created, err := archiveOp.Create(ctx, zoneFrom, &iaas.ArchiveCreateRequest{
				SourceArchiveID: source.ID,
				Name:            testutil.ResourceName("source-archive-for-transfer"),
			})
			if err != nil {
				return err
			}
			sourceArchive = created
			_, err = iaas.WaiterForReady(func() (interface{}, error) {
				return archiveOp.Read(ctx, zoneFrom, sourceArchive.ID)
			}).WaitForState(ctx)

			return err
		},
		Create: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
				builder := &TransferArchiveBuilder{
					Name:              testutil.ResourceName("archive-from-other-zone"),
					Description:       "description",
					Tags:              types.Tags{"tag1", "tag2"},
					SourceArchiveID:   sourceArchive.ID,
					SourceArchiveZone: zoneFrom,
					Client:            NewAPIClient(caller),
				}
				return builder.Build(ctx, zoneTo)
			},
		},
		Read: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
				return iaas.NewArchiveOp(caller).Read(ctx, zoneTo, ctx.ID)
			},
			CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, value interface{}) error {
				archive := value.(*iaas.Archive)
				return testutil.DoAsserts(
					testutil.AssertNotNilFunc(t, archive, "Archive"),
				)
			},
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
				archiveOp := iaas.NewArchiveOp(caller)

				_, err := iaas.WaiterForReady(func() (interface{}, error) {
					return archiveOp.Read(ctx, zoneTo, ctx.ID)
				}).WaitForState(ctx)
				if err != nil {
					return err
				}
				if err := archiveOp.Delete(ctx, zoneTo, ctx.ID); err != nil {
					return err
				}

				if sourceArchive != nil {
					if err := archiveOp.Delete(ctx, zoneFrom, sourceArchive.ID); err != nil {
						return err
					}
				}
				return nil
			},
		},
	})
}

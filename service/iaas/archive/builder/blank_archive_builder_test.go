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
	"bytes"
	"context"
	"testing"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/stretchr/testify/require"
)

func TestBlankArchiveBuilder_Validate(t *testing.T) {
	builder := &BlankArchiveBuilder{
		Name:         "test",
		SizeGB:       20,
		SourceReader: bytes.NewBufferString(""),
		NoWait:       true,
		Client:       NewAPIClient(testutil.SingletonAPICaller()),
	}

	err := builder.Validate(context.Background(), testutil.TestZone())
	require.EqualError(t, err, "NoWait=true is not supported when uploading files and creating archives")
}

func TestBlankArchiveBuilder_Build(t *testing.T) {
	if !testutil.IsAccTest() {
		t.Skip("TestBlankArchiveBuilder_Build only exec when running an Acceptance Test")
	}

	testZone := testutil.TestZone()
	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		SetupAPICallerFunc: testutil.SingletonAPICaller,
		Parallel:           true,
		IgnoreStartupWait:  true,
		Create: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
				source := bytes.NewBufferString("dummy")

				builder := &BlankArchiveBuilder{
					Name:         testutil.ResourceName("archive-from-shared-builder"),
					Description:  "description",
					Tags:         types.Tags{"tag1", "tag2"},
					SizeGB:       20,
					SourceReader: source,
					Client:       NewAPIClient(caller),
				}
				return builder.Build(ctx, testZone)
			},
		},
		Read: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
				return iaas.NewArchiveOp(caller).Read(ctx, testZone, ctx.ID)
			},
			CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, value interface{}) error {
				archive := value.(*iaas.Archive)
				return testutil.DoAsserts(
					testutil.AssertNotNilFunc(t, archive, "Archive"),
					testutil.AssertTrueFunc(t, archive.Availability.IsAvailable(), "Archive.Availability.IsAvailable"),
				)
			},
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
				return iaas.NewArchiveOp(caller).Delete(ctx, testZone, ctx.ID)
			},
		},
	})
}

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

package test

import (
	"testing"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
)

func TestNoteOp_CRUD(t *testing.T) {
	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		Parallel:           true,
		IgnoreStartupWait:  true,
		SetupAPICallerFunc: singletonAPICaller,
		Create: &testutil.CRUDTestFunc{
			Func: testNoteCreate,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createNoteExpected,
				IgnoreFields: ignoreNoteFields,
			}),
		},
		Read: &testutil.CRUDTestFunc{
			Func: testNoteRead,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createNoteExpected,
				IgnoreFields: ignoreNoteFields,
			}),
		},
		Updates: []*testutil.CRUDTestFunc{
			{
				Func: testNoteUpdate,
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  updateNoteExpected,
					IgnoreFields: ignoreNoteFields,
				}),
			},
			{
				Func: testNoteUpdateToMin,
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  updateNoteToMinExpected,
					IgnoreFields: ignoreNoteFields,
				}),
			},
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: testNoteDelete,
		},
	})
}

var (
	ignoreNoteFields = []string{"ID", "CreatedAt", "ModifiedAt"}
	createNoteParam  = &iaas.NoteCreateRequest{
		Name:    testutil.ResourceName("note"),
		Tags:    []string{"tag1", "tag2"},
		Class:   "shell",
		Content: "test-content",
	}
	createNoteExpected = &iaas.Note{
		Name:         createNoteParam.Name,
		Tags:         createNoteParam.Tags,
		Class:        createNoteParam.Class,
		Content:      createNoteParam.Content,
		Scope:        types.Scopes.User,
		Availability: types.Availabilities.Available,
	}
	updateNoteParam = &iaas.NoteUpdateRequest{
		Name:    testutil.ResourceName("note-upd"),
		Tags:    []string{"tag1-upd", "tag2-upd"},
		Class:   "shell",
		Content: "test-content-upd",
		IconID:  testIconID,
	}
	updateNoteExpected = &iaas.Note{
		Name:         updateNoteParam.Name,
		Tags:         updateNoteParam.Tags,
		Class:        updateNoteParam.Class,
		Content:      updateNoteParam.Content,
		Scope:        types.Scopes.User,
		Availability: types.Availabilities.Available,
		IconID:       updateNoteParam.IconID,
	}
	updateNoteToMinParam = &iaas.NoteUpdateRequest{
		Name:    testutil.ResourceName("note-to-min"),
		Class:   "shell",
		Content: "test-content-upd",
	}
	updateNoteToMinExpected = &iaas.Note{
		Name:         updateNoteToMinParam.Name,
		Class:        updateNoteToMinParam.Class,
		Content:      updateNoteToMinParam.Content,
		Scope:        types.Scopes.User,
		Availability: types.Availabilities.Available,
	}
)

func testNoteCreate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewNoteOp(caller)
	return client.Create(ctx, createNoteParam)
}

func testNoteRead(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewNoteOp(caller)
	return client.Read(ctx, ctx.ID)
}

func testNoteUpdate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewNoteOp(caller)
	return client.Update(ctx, ctx.ID, updateNoteParam)
}

func testNoteUpdateToMin(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewNoteOp(caller)
	return client.Update(ctx, ctx.ID, updateNoteToMinParam)
}

func testNoteDelete(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
	client := iaas.NewNoteOp(caller)
	return client.Delete(ctx, ctx.ID)
}

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
	"errors"
	"os"
	"testing"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/packages-go/size"
	"github.com/stretchr/testify/assert"
)

func TestDiskOp_BlankDiskCRUD(t *testing.T) {
	if testutil.IsAccTest() && os.Getenv("SAKURACLOUD_KMS_KEY_ID") == "" {
		t.Skip("SAKURACLOUD_KMS_KEY_ID is required when running the acceptance test")
	}

	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		Parallel: true,

		SetupAPICallerFunc: singletonAPICaller,

		Create: &testutil.CRUDTestFunc{
			Func: testDiskCreate,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createDiskExpected,
				IgnoreFields: ignoreDiskFields,
			}),
		},

		Read: &testutil.CRUDTestFunc{
			Func: testDiskRead,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createDiskExpected,
				IgnoreFields: ignoreDiskFields,
			}),
		},

		Updates: []*testutil.CRUDTestFunc{
			{
				Func: testDiskUpdate,
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  updateDiskExpected,
					IgnoreFields: ignoreDiskFields,
				}),
			},
			{
				Func: testDiskUpdateToMin,
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  updateDiskToMinExpected,
					IgnoreFields: ignoreDiskFields,
				}),
			},
		},

		Delete: &testutil.CRUDTestDeleteFunc{
			Func: testDiskDelete,
		},
	})
}

var (
	ignoreDiskFields = []string{
		"ID",
		"DisplayOrder",
		"Availability",
		"DiskPlanName",
		"DiskPlanStorageClass",
		"SizeMB",
		"MigratedMB",
		"SourceDiskID",
		"SourceDiskAvailability",
		"SourceArchiveID",
		"SourceArchiveAvailability",
		"BundleInfo",
		"Server",
		"Storage",
		"CreatedAt",
		"ModifiedAt",
	}

	createDiskParam = &iaas.DiskCreateRequest{
		DiskPlanID:          types.DiskPlans.SSD,
		Connection:          types.DiskConnections.VirtIO,
		EncryptionAlgorithm: types.DiskEncryptionAlgorithms.AES256XTS,
		Name:                testutil.ResourceName("disk"),
		Description:         "desc",
		Tags:                []string{"tag1", "tag2"},
		SizeMB:              20 * size.GiB,
	}
	createDiskExpected = &iaas.Disk{
		Name:                createDiskParam.Name,
		Description:         createDiskParam.Description,
		Tags:                createDiskParam.Tags,
		DiskPlanID:          createDiskParam.DiskPlanID,
		Connection:          createDiskParam.Connection,
		EncryptionAlgorithm: createDiskParam.EncryptionAlgorithm,
		KMSKeyID:            types.StringID(os.Getenv("SAKURACLOUD_KMS_KEY_ID")),
	}
	updateDiskParam = &iaas.DiskUpdateRequest{
		Name:        testutil.ResourceName("disk-upd"),
		Description: "desc-upd",
		Tags:        []string{"tag1-upd", "tag2-upd"},
		IconID:      testIconID,
	}
	updateDiskExpected = &iaas.Disk{
		Name:                updateDiskParam.Name,
		Description:         updateDiskParam.Description,
		Tags:                updateDiskParam.Tags,
		DiskPlanID:          createDiskParam.DiskPlanID,
		Connection:          createDiskParam.Connection,
		EncryptionAlgorithm: createDiskParam.EncryptionAlgorithm,
		KMSKeyID:            types.StringID(os.Getenv("SAKURACLOUD_KMS_KEY_ID")),
		IconID:              updateDiskParam.IconID,
	}
	updateDiskToMinParam = &iaas.DiskUpdateRequest{
		Name: testutil.ResourceName("disk-to-min"),
	}
	updateDiskToMinExpected = &iaas.Disk{
		Name:                updateDiskToMinParam.Name,
		DiskPlanID:          createDiskParam.DiskPlanID,
		Connection:          createDiskParam.Connection,
		EncryptionAlgorithm: createDiskParam.EncryptionAlgorithm,
		KMSKeyID:            types.StringID(os.Getenv("SAKURACLOUD_KMS_KEY_ID")),
	}
)

func testDiskCreate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewDiskOp(caller)
	return client.Create(ctx, testZone, createDiskParam, nil, types.StringID(os.Getenv("SAKURACLOUD_KMS_KEY_ID")))
}

func testDiskRead(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewDiskOp(caller)
	return client.Read(ctx, testZone, ctx.ID)
}

func testDiskUpdate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewDiskOp(caller)
	return client.Update(ctx, testZone, ctx.ID, updateDiskParam)
}

func testDiskUpdateToMin(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewDiskOp(caller)
	return client.Update(ctx, testZone, ctx.ID, updateDiskToMinParam)
}

func testDiskDelete(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
	client := iaas.NewDiskOp(caller)
	return client.Delete(ctx, testZone, ctx.ID)
}

func TestDiskOp_Config(t *testing.T) {
	// source archive
	var archiveID types.ID

	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		Parallel:           true,
		SetupAPICallerFunc: singletonAPICaller,
		Setup: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			client := iaas.NewArchiveOp(singletonAPICaller())
			searched, err := client.Find(ctx, testZone, &iaas.FindCondition{
				Filter: search.Filter{
					search.Key("Tags.Name"): search.TagsAndEqual("current-stable", "distro-ubuntu"),
				},
			})
			if !assert.NoError(t, err) {
				return err
			}
			if searched.Count == 0 {
				return errors.New("archive is not found")
			}
			archiveID = searched.Archives[0].ID
			return nil
		},
		Create: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
				client := iaas.NewDiskOp(singletonAPICaller())
				disk, err := client.Create(ctx, testZone, &iaas.DiskCreateRequest{
					Name:            testutil.ResourceName("disk-edit"),
					DiskPlanID:      types.DiskPlans.SSD,
					SizeMB:          20 * size.GiB,
					SourceArchiveID: archiveID,
				}, nil, 0)
				if err != nil {
					return nil, err
				}
				if _, err = iaas.WaiterForReady(func() (interface{}, error) {
					return client.Read(ctx, testZone, disk.ID)
				}).WaitForState(ctx); err != nil {
					return disk, err
				}

				return disk, nil
			},
		},
		Read: &testutil.CRUDTestFunc{
			Func: testDiskRead,
		},
		Updates: []*testutil.CRUDTestFunc{
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					// edit disk
					client := iaas.NewDiskOp(singletonAPICaller())
					err := client.Config(ctx, testZone, ctx.ID, &iaas.DiskEditRequest{
						Background: true,
						Password:   "password",
						SSHKeys: []*iaas.DiskEditSSHKey{
							{
								PublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC4LDQuDiKecOJDPY9InS7EswZ2fPnoRZXc48T1EqyRLyJhgEYGSDWaBiMDs2R/lWgA81Hp37qhrNqZPjFHUkBr93FOXxt9W0m1TNlkNepK0Uyi+14B2n0pdoeqsKEkb3sTevWF0ztxxWrwUd7Mems2hf+wFODITHYye9RlDAKLKPCFRvlQ9xQj4bBWOogQwoaXMSK1znMPjudcm1tRry4KIifLdXmwVKU4qDPGxoXfqs44Dgsikk43UVBStQ7IFoqPgAqcJFSGHLoMS7tPKdTvY9+GME5QidWK84gl69piAkgIdwd+JTMUOc/J+9DXAt220HqZ6l3yhWG5nIgi0x8n",
							},
						},
						DisablePWAuth: true,
						EnableDHCP:    true,
						HostName:      "hostname",
						UserIPAddress: "192.2.0.11",
						UserSubnet: &iaas.DiskEditUserSubnet{
							DefaultRoute:   "192.2.0.1",
							NetworkMaskLen: 24,
						},
					})
					if err != nil {
						return nil, err
					}
					// wait
					_, err = iaas.WaiterForReady(func() (interface{}, error) {
						return client.Read(ctx, testZone, ctx.ID)
					}).WaitForState(ctx)
					return nil, err
				},
			},
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: testDiskDelete,
		},
	})
}

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
)

func TestSSHKeyOpCRUD(t *testing.T) {
	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		Parallel:           true,
		IgnoreStartupWait:  true,
		SetupAPICallerFunc: singletonAPICaller,
		Create: &testutil.CRUDTestFunc{
			Func: testSSHKeyCreate,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createSSHKeyExpected,
				IgnoreFields: ignoreSSHKeyFields,
			}),
		},
		Read: &testutil.CRUDTestFunc{
			Func: testSSHKeyRead,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createSSHKeyExpected,
				IgnoreFields: ignoreSSHKeyFields,
			}),
		},
		Updates: []*testutil.CRUDTestFunc{
			{
				Func: testSSHKeyUpdate,
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  updateSSHKeyExpected,
					IgnoreFields: ignoreSSHKeyFields,
				}),
			},
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: testSSHKeyDelete,
		},
	})
}

var (
	fakePublicKey   = `ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAs7YFtxjGrI49MCBnSFbUPxqz0e5HSGQPnLlPJ0u/9w4WLpoOZYmoQDTMfuFA61qv+0dp5mpMZPj3f5YEGlwUFKPy3Cmrp0ub1nYDb7n62s+Xf68TNvbVgQMLF0xdOaWxdRsQwmH8lOWan1Ubc8iwfOa3TNGwOzGLMjdW3PiJ7hcE7nFqnmbQUabHWow8G6JYDHKyjAdpz+edK8u+LY0iEP8M8VAjRJKJVg4p1/oDjHFKI0qjfjitKzoLm5FGaFv8afH2WQSpu/2To7d/RaLhfoMZsUReLSxeDnQkKGERXrAywTHnFu60cOaT3EvaAhP1H3BPj2LESm8M4ja9FaARnQ== `
	fakeFingerprint = "79:d7:ac:b8:cf:cf:01:44:b2:19:ba:d4:82:fd:c4:2d"

	ignoreSSHKeyFields = []string{
		"ID",
		"CreatedAt",
	}
	createSSHKeyParam = &iaas.SSHKeyCreateRequest{
		Name:        testutil.ResourceName("sshkey"),
		Description: "libsacloud-sshKey",
		PublicKey:   fakePublicKey,
	}
	createSSHKeyExpected = &iaas.SSHKey{
		Name:        createSSHKeyParam.Name,
		Description: createSSHKeyParam.Description,
		PublicKey:   fakePublicKey,
		Fingerprint: fakeFingerprint,
	}
	updateSSHKeyParam = &iaas.SSHKeyUpdateRequest{
		Name:        testutil.ResourceName("sshkey-upd"),
		Description: "libsacloud-sshKey-upd",
	}
	updateSSHKeyExpected = &iaas.SSHKey{
		Name:        updateSSHKeyParam.Name,
		Description: updateSSHKeyParam.Description,
		PublicKey:   fakePublicKey,
		Fingerprint: fakeFingerprint,
	}
)

func testSSHKeyCreate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewSSHKeyOp(caller)
	return client.Create(ctx, createSSHKeyParam)
}

func testSSHKeyRead(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewSSHKeyOp(caller)
	return client.Read(ctx, ctx.ID)
}

func testSSHKeyUpdate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewSSHKeyOp(caller)
	return client.Update(ctx, ctx.ID, updateSSHKeyParam)
}

func testSSHKeyDelete(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
	client := iaas.NewSSHKeyOp(caller)
	return client.Delete(ctx, ctx.ID)
}

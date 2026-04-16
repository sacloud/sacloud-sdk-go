// Copyright 2025- The sacloud/kms-api-go authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apigw_test

import (
	"context"
	"os"
	"testing"

	apigw "github.com/sacloud/apigw-api-go"
	v1 "github.com/sacloud/apigw-api-go/apis/v1"
	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomainAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN",
		"SAKURA_ACCESS_TOKEN_SECRET", "SAKURA_TEST_DOMAIN")(t)

	var theClient saclient.Client
	client, err := apigw.NewClient(&theClient)
	require.Nil(t, err)

	ctx := context.Background()
	domainOp := apigw.NewDomainOp(client)

	domain, err := domainOp.Create(ctx, &v1.Domain{DomainName: os.Getenv("SAKURA_TEST_DOMAIN")})
	require.Nil(t, err)

	// TODO: 証明書を作ってのテストも行うようにする
	// err = domainOp.Update(ctx, &v1.DomainPUT{CertificateId: v1.NewOptUUID(cert.ID.Value)}, domain.ID.Value)

	domains, err := domainOp.List(ctx)
	assert.Nil(t, err)
	assert.Greater(t, len(domains), 0)

	err = domainOp.Delete(ctx, domain.ID.Value)
	assert.Nil(t, err)
}

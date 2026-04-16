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

func TestCertificateAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET")(t)

	var theClient saclient.Client
	client, err := apigw.NewClient(&theClient)
	require.Nil(t, err)

	ctx := context.Background()
	certOp := apigw.NewCertificateOp(client)

	/*
	 * 以下の証明書ファイルはテスト用に適当な内容で生成したものです。
	 * % openssl genpkey -algorithm RSA -out rsa.key
	 * % openssl req -new -x509 -key rsa.key -out rsa.crt -days 1000
	 * % openssl ecparam -name prime256v1 -out server.key -genkey
	 * % openssl req -new -x509 -key ecdsa.key -out ecdsa.crt -days 1000
	 */
	crtRsa, _ := os.ReadFile("./testdata/rsa.crt")
	keyRsa, _ := os.ReadFile("./testdata/rsa.key")
	crtEcdsa, _ := os.ReadFile("./testdata/ecdsa.crt")
	keyEcdsa, _ := os.ReadFile("./testdata/ecdsa.key")
	certArg := v1.Certificate{
		Name: v1.NewOptName("test-cert"),
		Rsa: v1.NewOptCertificateDetails(v1.CertificateDetails{
			Cert: v1.NewOptString(string(crtRsa)),
			Key:  v1.NewOptString(string(keyRsa)),
		}),
		Ecdsa: v1.NewOptCertificateDetails(v1.CertificateDetails{
			Cert: v1.NewOptString(string(crtEcdsa)),
			Key:  v1.NewOptString(string(keyEcdsa)),
		}),
	}

	cert, err := certOp.Create(ctx, &certArg)
	require.Nil(t, err)

	err = certOp.Update(ctx, &certArg, cert.ID.Value)
	assert.Nil(t, err)

	certs, err := certOp.List(ctx)
	assert.Nil(t, err)
	assert.Greater(t, len(certs), 0)

	err = certOp.Delete(ctx, cert.ID.Value)
	assert.Nil(t, err)
}

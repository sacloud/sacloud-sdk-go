// Copyright 2026 The sacloud/iaas-api-go Authors
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

package fake

import (
	"context"
	"testing"

	"github.com/sacloud/sacloud-sdk-go/api/iaas"
	"github.com/stretchr/testify/require"
)

func TestProxyLBOpSetCertificates(t *testing.T) {
	// Exercise the public factory path where PrimaryCerts must populate the
	// differently named PrimaryCert field in the fake response.
	SwitchFactoryFuncToFake()

	ctx := context.Background()
	op := iaas.NewProxyLBOp(nil)
	proxyLB, err := op.Create(ctx, &iaas.ProxyLBCreateRequest{})
	require.NoError(t, err)

	request := &iaas.ProxyLBSetCertificatesRequest{
		PrimaryCerts: &iaas.ProxyLBPrimaryCert{
			ServerCertificate:       "server-certificate",
			IntermediateCertificate: "intermediate-certificate",
			PrivateKey:              "private-key",
		},
	}
	certificates, err := op.SetCertificates(ctx, proxyLB.ID, request)
	require.NoError(t, err)
	require.NotNil(t, certificates.PrimaryCert)
	require.Equal(t, request.PrimaryCerts.ServerCertificate, certificates.PrimaryCert.ServerCertificate)
	require.Equal(t, request.PrimaryCerts.IntermediateCertificate, certificates.PrimaryCert.IntermediateCertificate)
	require.Equal(t, request.PrimaryCerts.PrivateKey, certificates.PrimaryCert.PrivateKey)
	require.Equal(t, "dummy-common-name.org", certificates.PrimaryCert.CertificateCommonName)
	require.False(t, certificates.PrimaryCert.CertificateEndDate.IsZero())

	stored, err := op.GetCertificates(ctx, proxyLB.ID)
	require.NoError(t, err)
	require.Equal(t, certificates, stored)
}

func TestProxyLBOpSetCertificatesWithoutPrimaryCert(t *testing.T) {
	SwitchFactoryFuncToFake()

	ctx := context.Background()
	op := iaas.NewProxyLBOp(nil)
	proxyLB, err := op.Create(ctx, &iaas.ProxyLBCreateRequest{})
	require.NoError(t, err)

	certificates, err := op.SetCertificates(ctx, proxyLB.ID, &iaas.ProxyLBSetCertificatesRequest{})
	require.NoError(t, err)
	require.NotNil(t, certificates.PrimaryCert)
	require.Equal(t, "dummy-common-name.org", certificates.PrimaryCert.CertificateCommonName)
	require.False(t, certificates.PrimaryCert.CertificateEndDate.IsZero())
}

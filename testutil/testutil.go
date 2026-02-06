// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/go-faster/jx"
	"github.com/google/uuid"
	apprun_dedicated "github.com/sacloud/apprun-dedicated-api-go"
	v1 "github.com/sacloud/apprun-dedicated-api-go/apis/v1"
	super "github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
)

var theClient saclient.Client

//nolint:nakedret
func NewTestClient(
	v interface{ Encode(*jx.Encoder) },
	s ...int,
) (
	c *v1.Client,
	e error,
) {
	s = append(s, http.StatusOK)
	f := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(s[0])

		if s[0] == http.StatusNoContent {
			return
		}

		if v == nil {
			return
		}

		e := new(jx.Encoder)
		v.Encode(e)
		_, _ = w.Write(e.Bytes())
	}
	h := http.HandlerFunc(f)
	sv := httptest.NewServer(h)

	sa, e := theClient.DupWith(saclient.WithTestServer(sv))

	if e != nil {
		return
	}

	c, e = apprun_dedicated.NewClientWithAPIRootURL(sa, sv.URL)

	if e != nil {
		return
	}

	e = sa.Populate()

	if e != nil {
		c = nil
	}

	return
}

func FakeCN() string {
	h := super.RandomName("host", 32, super.CharSetAlphaNum)
	return fmt.Sprintf("%s.xn--eckwd4c7cu47r2wf.jp", h)
}

func FakeCertificate() (ret v1.ReadCertificate) {
	ret.SetFake()
	ret.SetSubjectAlternativeNames(make([]string, 3))
	for i := range ret.SubjectAlternativeNames {
		ret.SubjectAlternativeNames[i] = FakeCN()
	}
	return
}

func Fake400Error() (ret v1.Error) {
	ret.SetFake()
	ret.SetStatus(http.StatusBadRequest)
	return
}

func Fake403Error() (ret v1.Error) {
	ret.SetFake()
	ret.SetStatus(http.StatusForbidden)
	return
}

func Fake404Error() (ret v1.Error) {
	ret.SetFake()
	ret.SetStatus(http.StatusNotFound)
	return
}

//nolint:nakedret
func OreSign() (certPEM, keyPEM []byte, err error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	if err != nil {
		return
	}

	key, err := x509.MarshalECPrivateKey(priv)

	if err != nil {
		return
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))

	if err != nil {
		return
	}

	tpl := x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName: FakeCN(),
		},
		NotBefore: time.Now().Add(-time.Hour),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),

		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	der, err := x509.CreateCertificate(rand.Reader, &tpl, &tpl, &priv.PublicKey, priv)

	if err != nil {
		return
	}

	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: key})
	return
}

func FakeApplicationVersion() (ret v1.ReadApplicationVersionDetail) {
	ret.SetFake()
	ret.SetVersion(1)
	ret.SetCPU(128)
	ret.SetMemory(128)
	ret.FixedScale.Reset()
	ret.MinScale.Reset()
	ret.MaxScale.Reset()
	ret.ScaleInThreshold.Reset()
	ret.ScaleOutThreshold.Reset()
	return
}

func FakeAutoScalingGroupNodeInterface() (ret v1.AutoScalingGroupNodeInterface) {
	ret.SetInterfaceIndex(0)
	ret.SetUpstream("shared")
	ret.SetIpPool(make([]v1.IpRange, 0))
	ret.NetmaskLen.SetTo(24)
	ret.DefaultGateway.SetTo("192.168.1.1")
	ret.PacketFilterID.SetTo("filter-id")
	ret.SetConnectsToLB(true)
	return
}

func FakeAutoScalingGroup() (ret v1.ReadAutoScalingGroupDetail) {
	ret.SetAutoScalingGroupID(v1.AutoScalingGroupID(uuid.New()))
	ret.SetName("test-asg")
	ret.SetZone("is1a")
	ret.SetNameServers([]v1.IPv4{"133.242.0.3"})
	ret.SetWorkerServiceClassPath("/is1a/server/1/core-2-2")
	ret.SetMinNodes(1)
	ret.SetMaxNodes(3)
	ret.SetWorkerNodeCount(2)
	ret.SetDeleting(false)
	ret.SetInterfaces([]v1.AutoScalingGroupNodeInterface{FakeAutoScalingGroupNodeInterface()})
	return
}

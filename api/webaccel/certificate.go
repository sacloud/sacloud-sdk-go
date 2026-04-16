// Copyright 2022-2025 The webaccel-api-go authors
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

package webaccel

import (
	"encoding/json"
	"strings"
	"time"
)

// Certificates 証明書
type Certificates struct {
	Current *CurrentCertificate
	Old     []*OldCertificate
}

// UnmarshalJSON JSONアンマーシャル(配列、オブジェクトが混在するためここで対応)
func (w *Certificates) UnmarshalJSON(data []byte) error {
	targetData := strings.ReplaceAll(strings.ReplaceAll(string(data), " ", ""), "\n", "")
	if targetData == `[]` {
		return nil
	}

	type alias Certificates

	var tmp alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	*w = Certificates(tmp)
	return nil
}

// CurrentCertificate 現在有効な証明書
type CurrentCertificate struct {
	ID                string
	SiteID            string
	CertificateChain  string
	Key               string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	SerialNumber      string
	NotBefore         int64
	NotAfter          int64
	Issuer            *Issuer
	Subject           *Subject
	DNSNames          []string
	SHA256Fingerprint string
}

// Issuer .
type Issuer struct {
	Country            string
	Organization       string
	OrganizationalUnit string
	CommonName         string
}

// Subject .
type Subject struct {
	Country            string
	Organization       string
	OrganizationalUnit string
	Locality           string
	Province           string
	StreetAddress      string
	PostalCode         string
	SerialNumber       string
	CommonName         string
}

// OldCertificate .
type OldCertificate struct {
	ID                string
	SiteID            string
	CertificateChain  string
	Key               string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	SerialNumber      string
	NotBefore         int64
	NotAfter          int64
	Issuer            *Issuer
	Subject           *Subject
	DNSNames          []string
	SHA256Fingerprint string
}

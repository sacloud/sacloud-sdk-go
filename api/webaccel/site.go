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

import "time"

// Site サイト
type Site struct {
	ID   string
	Name string

	DomainType      string `validate:"oneof=own_domain subdomain"`
	Domain          string
	Subdomain       string
	ASCIIDomain     string
	RequestProtocol string // 0:http/https, 1:httpsのみ, 2:httpsにリダイレクト
	DefaultCacheTTL int    `validate:"min=-1,max=604800"` // -1:無効, 0 ～ 604800 の範囲内の数値: デフォルトのキャッシュ期間(秒)
	VarySupport     string // 0:無効, 1:有効
	// NOTE: List()だと空、Read()でのみ参照可能
	NormalizeAE string // 1:Accept-Encodingをgzipに正規化, 3:bzとgzipの組に正規化

	OriginType           string // 0:ウェブサーバ, 1:オブジェクトストレージ
	Origin               string
	OriginProtocol       string
	HostHeader           string
	OriginGuardToken     string
	NextOriginGuardToken string

	// オブジェクトストレージをオリジンにする場合
	BucketName string
	S3Endpoint string
	S3Region   string
	DocIndex   string // 0:無効, 1:有効

	// CORSRules ルール一覧、設定されている場合単一要素を持つ配列となる
	// NOTE: List()だと空、Read()でのみ参照可能
	CORSRules         []*CORSRule
	OnetimeURLSecrets []string

	Status             string `validate:"oneof=enabled disabled"`
	HasCertificate     bool
	HasOldCertificate  bool
	GibSentInLastWeek  int64
	CertValidNotBefore int64
	CertValidNotAfter  int64
	CreatedAt          time.Time
}

// CORSRule .
type CORSRule struct {
	AllowsAnyOrigin bool     // trueの場合全オリジンを許可
	AllowedOrigins  []string `validate:"omitempty,max=4"`
}

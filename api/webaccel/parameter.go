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

// DeleteAllCacheRequest .
type DeleteAllCacheRequest struct {
	Domain string
}

// DeleteCacheRequest .
type DeleteCacheRequest struct {
	URL []string
}

// DeleteCacheResult .
type DeleteCacheResult struct {
	URL    string
	Status int
	Result string
}

// ListSitesResult .
type ListSitesResult struct {
	Total int `json:",omitempty"` // Total count of target resources
	From  int `json:",omitempty"` // Current page number
	Count int `json:",omitempty"` // Count of current page

	Sites []*Site `json:",omitempty"`
}

// CreateOrUpdateCertificateRequest .
type CreateOrUpdateCertificateRequest struct {
	CertificateChain string
	Key              string
}

// CreateSiteRequest サイト作成リクエスト
type CreateSiteRequest struct {
	// 「オリジン種別」に関係なく設定できる共通項目
	Name            string `json:",omitempty"`
	OriginType      string `json:",omitempty" validate:"omitempty,oneof=0 1"` // 0:ウェブサーバ, 1:オブジェクトストレージ
	Domain          string `json:",omitempty" validate:"omitempty"`           // 独自ドメインを設定する場合のみ
	DomainType      string `json:",omitempty" validate:"omitempty,oneof=own_domain subdomain"`
	RequestProtocol string `json:",omitempty" validate:"omitempty,oneof=0 1 2"` // 0:http/https, 1:httpsのみ, 2:httpsにリダイレクト
	OriginProtocol  string `json:",omitempty" validate:"omitempty,oneof=http https"`
	DefaultCacheTTL *int   `json:",omitempty" validate:"omitempty,min=-1,max=604800"` // -1:無効, 0 ～ 604800 の範囲内の数値: デフォルトのキャッシュ期間(秒)
	VarySupport     string `json:",omitempty" validate:"omitempty,oneof=0 1"`         // 0:無効, 1:有効
	NormalizeAE     string `json:",omitempty" validate:"omitempty,oneof=1 3"`         // 1:Accept-Encodingをgzipに正規化, 3:bzとgzipの組に正規化

	// 「オリジン種別」が「ウェブサーバ」の場合に設定可能な項目
	Origin     string `json:",omitempty"`
	HostHeader string `json:",omitempty"`

	// 「オリジン種別」が「オブジェクトストレージ」の場合に設定可能な項目
	BucketName      string `json:",omitempty"`
	S3Endpoint      string `json:",omitempty"`
	S3Region        string `json:",omitempty"`
	DocIndex        string `json:",omitempty" validate:"omitempty,oneof=0 1"` // 0:無効, 1:有効
	AccessKeyID     string `json:",omitempty"`
	SecretAccessKey string `json:",omitempty"`
}

// UpdateSiteRequest サイト更新リクエスト
type UpdateSiteRequest struct {
	// 「オリジン種別」に関係なく設定できる共通項目
	Name            string `json:",omitempty"`
	OriginType      string `json:",omitempty" validate:"omitempty,oneof=0 1"`   // 0:ウェブサーバ, 1:オブジェクトストレージ
	RequestProtocol string `json:",omitempty" validate:"omitempty,oneof=0 1 2"` // 0:http/https, 1:httpsのみ, 2:httpsにリダイレクト
	OriginProtocol  string `json:",omitempty" validate:"omitempty,oneof=http https"`
	DefaultCacheTTL *int   `json:",omitempty" validate:"omitempty,min=-1,max=604800"` // -1:無効, 0 ～ 604800 の範囲内の数値: デフォルトのキャッシュ期間(秒)
	VarySupport     string `json:",omitempty" validate:"omitempty,oneof=0 1"`         // 0:無効, 1:有効
	NormalizeAE     string `json:",omitempty" validate:"omitempty,oneof=1 3"`         // 1:Accept-Encodingをgzipに正規化, 3:bzとgzipの組に正規化

	// CORSRules ルール一覧、設定されている場合単一要素を持つ配列となる
	CORSRules         *[]*CORSRule `json:",omitempty"`
	OnetimeURLSecrets *[]string    `json:",omitempty"`

	// 「オリジン種別」が「ウェブサーバ」の場合に設定可能な項目
	Origin     string `json:",omitempty"`
	HostHeader string `json:",omitempty"`

	// 「オリジン種別」が「オブジェクトストレージ」の場合に設定可能な項目
	BucketName      string `json:",omitempty"`
	S3Endpoint      string `json:",omitempty"`
	S3Region        string `json:",omitempty"`
	DocIndex        string `json:",omitempty" validate:"omitempty,oneof=0 1"` // 0:無効, 1:有効
	AccessKeyID     string `json:",omitempty"`
	SecretAccessKey string `json:",omitempty"`
}

// UpdateSiteStatusRequest サイトの有効化状態更新リクエスト
type UpdateSiteStatusRequest struct {
	Status string `validate:"oneof=enabled disabled"`
}

type ACLResult struct {
	ACL string
}

// OriginGuardTokenResponse サイトのオリジンガードトークン取得レスポンス
type OriginGuardTokenResponse struct {
	OriginGuardToken     string `json:","`
	NextOriginGuardToken string `json:","`
}

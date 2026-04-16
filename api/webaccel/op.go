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
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type APICaller interface {
	Do(ctx context.Context, method, uri string, body interface{}) ([]byte, error)
	RootURL() string
}

var _ API = (*Op)(nil)

// Op implements WebAccelAPI interface
type Op struct {
	Client APICaller
}

// NewOp creates new Op instance
func NewOp(caller APICaller) API {
	return &Op{Client: caller}
}

// Create 新規サイトの作成
func (o *Op) Create(ctx context.Context, param *CreateSiteRequest) (*Site, error) {
	url := o.Client.RootURL() + "site"

	// build request body
	type createRequest struct {
		Site *CreateSiteRequest
	}
	body := &createRequest{Site: param}

	// do request
	data, err := o.Client.Do(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}

	// build results
	type createResult struct {
		Site *Site
	}
	var results createResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return results.Site, nil
}

// List サイト一覧
//
// NOTE: 各サイトのCORSRulesはnullになる点に注意
func (o *Op) List(ctx context.Context) (*ListSitesResult, error) {
	url := o.Client.RootURL() + "site"

	// build request body
	var body interface{}

	// do request
	data, err := o.Client.Do(ctx, "GET", url, body)
	if err != nil {
		return nil, err
	}

	// build results
	var results ListSitesResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return &results, nil
}

// Read サイト詳細
func (o *Op) Read(ctx context.Context, id string) (*Site, error) {
	url := o.Client.RootURL() + fmt.Sprintf("site/%s", id)

	// build request body
	var body interface{}

	// do request
	data, err := o.Client.Do(ctx, "GET", url, body)
	if err != nil {
		return nil, err
	}

	// build results
	type readResult struct {
		Site *Site
	}
	var results readResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return results.Site, nil
}

// Update サイト更新
func (o *Op) Update(ctx context.Context, id string, param *UpdateSiteRequest) (*Site, error) {
	url := o.Client.RootURL() + fmt.Sprintf("site/%s", id)

	// build request body
	type updateRequest struct {
		Site *UpdateSiteRequest
	}
	body := &updateRequest{Site: param}

	// do request
	data, err := o.Client.Do(ctx, "PUT", url, body)
	if err != nil {
		return nil, err
	}

	// build results
	type updateResult struct {
		Site *Site
	}
	var results updateResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return results.Site, nil
}

// UpdateStatus サイト有効化状態の更新
func (o *Op) UpdateStatus(ctx context.Context, id string, param *UpdateSiteStatusRequest) (*Site, error) {
	url := o.Client.RootURL() + fmt.Sprintf("site/%s/status", id)

	// build request body
	type updateStatusRequest struct {
		Site *UpdateSiteStatusRequest
	}
	body := &updateStatusRequest{Site: param}

	// do request
	data, err := o.Client.Do(ctx, "PUT", url, body)
	if err != nil {
		return nil, err
	}

	// build results
	type updateStatusResult struct {
		Site *Site
	}
	var results updateStatusResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return results.Site, nil
}

// ReadACL サイトのACL取得
func (o *Op) ReadACL(ctx context.Context, id string) (*ACLResult, error) {
	url := o.Client.RootURL() + fmt.Sprintf("site/%s/acl", id)

	// build request body
	var body interface{}

	// do request
	data, err := o.Client.Do(ctx, "GET", url, body)
	if err != nil {
		return nil, err
	}

	// build results
	var result ACLResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UpsertACL サイトのACLの登録/更新
func (o *Op) UpsertACL(ctx context.Context, id string, acl string) (*ACLResult, error) {
	url := o.Client.RootURL() + fmt.Sprintf("site/%s/acl", id)

	// build request body
	type upsertACLRequest struct {
		ACL string `validate:"required"`
	}
	body := &upsertACLRequest{ACL: acl}

	// do request
	data, err := o.Client.Do(ctx, "PUT", url, body)
	if err != nil {
		return nil, err
	}

	// build results
	var result ACLResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteACL サイトのACLの削除
func (o *Op) DeleteACL(ctx context.Context, id string) error {
	url := o.Client.RootURL() + fmt.Sprintf("site/%s/acl", id)

	// build request body
	var body interface{}

	// do request
	_, err := o.Client.Do(ctx, "DELETE", url, body)
	return err
}

// ReadCertificate サイト証明書の参照
func (o *Op) ReadCertificate(ctx context.Context, id string) (*Certificates, error) {
	url := o.Client.RootURL() + fmt.Sprintf("site/%s/certificate", id)

	// build request body
	var body interface{}

	// do request
	data, err := o.Client.Do(ctx, "GET", url, body)
	if err != nil {
		return nil, err
	}

	type readCertificateResult struct {
		Certificate *Certificates
	}

	// build results
	var results readCertificateResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return results.Certificate, nil
}

// CreateCertificate サイトに証明書を登録
func (o *Op) CreateCertificate(ctx context.Context, id string, param *CreateOrUpdateCertificateRequest) (*Certificates, error) {
	url := o.Client.RootURL() + fmt.Sprintf("site/%s/certificate", id)

	// build request body
	type createRequest struct {
		Certificate *CreateOrUpdateCertificateRequest
	}
	body := &createRequest{Certificate: param}

	// do request
	data, err := o.Client.Do(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}

	type createCertificateResult struct {
		Certificate *Certificates
	}

	// build results
	var results createCertificateResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return results.Certificate, nil
}

// UpdateCertificate サイトの証明書を更新
func (o *Op) UpdateCertificate(ctx context.Context, id string, param *CreateOrUpdateCertificateRequest) (*Certificates, error) {
	url := o.Client.RootURL() + fmt.Sprintf("site/%s/certificate", id)

	// build request body
	type updateRequest struct {
		Certificate *CreateOrUpdateCertificateRequest
	}
	body := &updateRequest{Certificate: param}

	// do request
	data, err := o.Client.Do(ctx, "PUT", url, body)
	if err != nil {
		return nil, err
	}

	type updateCertificateResult struct {
		Certificate *Certificates
	}

	// build results
	var results updateCertificateResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return results.Certificate, nil
}

// CreateAutoCertUpdate Let's Encrypt 自動更新証明書有効化
func (o *Op) CreateAutoCertUpdate(ctx context.Context, id string) error {
	url := o.Client.RootURL() + fmt.Sprintf("site/%s/auto-cert-update", id)

	// build request body
	type autoCertUpdateRequest struct {
		Type string
	}
	body := &autoCertUpdateRequest{Type: "letsencrypt"}

	// do request
	_, err := o.Client.Do(ctx, "POST", url, body)
	return err
}

// DeleteAutoCertUpdate Let's Encrypt 自動更新証明書無効化
func (o *Op) DeleteAutoCertUpdate(ctx context.Context, id string) error {
	url := o.Client.RootURL() + fmt.Sprintf("site/%s/auto-cert-update", id)

	// build request body
	var body interface{}
	_, err := o.Client.Do(ctx, "DELETE", url, body)
	return err
}

// CreateOriginGuardToken オリジンガードトークンの発行・更新
func (o *Op) CreateOriginGuardToken(ctx context.Context, id string) (*OriginGuardTokenResponse, error) {
	url := o.Client.RootURL() + fmt.Sprintf("site/%s/origin-guard-token", id)

	var body interface{}
	// do request
	data, err := o.Client.Do(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}

	var results OriginGuardTokenResponse
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return &results, nil
}

// CreateNextOriginGuardToken オリジンガードトークンの次期トークン作成
func (o *Op) CreateNextOriginGuardToken(ctx context.Context, id string) (*OriginGuardTokenResponse, error) {
	url := o.Client.RootURL() + fmt.Sprintf("site/%s/origin-guard-token/next", id)

	var body interface{}
	// do request
	data, err := o.Client.Do(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}

	var results OriginGuardTokenResponse
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return &results, nil
}

// DeleteOriginGuardToken オリジンガードトークンの削除
func (o *Op) DeleteOriginGuardToken(ctx context.Context, id string) error {
	url := o.Client.RootURL() + fmt.Sprintf("site/%s/origin-guard-token", id)

	var body interface{}
	// do request
	_, err := o.Client.Do(ctx, "DELETE", url, body)
	if err != nil {
		return err
	}
	return nil
}

// DeleteNextOriginGuardToken オリジンガードトークンの次期トークン更新キャンセル
func (o *Op) DeleteNextOriginGuardToken(ctx context.Context, id string) error {
	url := o.Client.RootURL() + fmt.Sprintf("site/%s/origin-guard-token/next", id)

	var body interface{}
	// do request
	_, err := o.Client.Do(ctx, "DELETE", url, body)
	if err != nil {
		return err
	}
	return nil
}

// DeleteCertificate サイトの証明書を削除
func (o *Op) DeleteCertificate(ctx context.Context, id string) error {
	url := o.Client.RootURL() + fmt.Sprintf("site/%s/certificate", id)

	// build request body
	var body interface{}

	// do request
	_, err := o.Client.Do(ctx, "DELETE", url, body)
	return err
}

// DeleteAllCache 全キャッシュ削除
func (o *Op) DeleteAllCache(ctx context.Context, param *DeleteAllCacheRequest) error {
	url := o.Client.RootURL() + "deleteallcache"

	// build request body
	type deleteAllCacheRequest struct {
		Site *DeleteAllCacheRequest
	}
	body := &deleteAllCacheRequest{Site: param}

	// do request
	_, err := o.Client.Do(ctx, "POST", url, body)
	return err
}

// DeleteCache URLごとにキャッシュ削除
func (o *Op) DeleteCache(ctx context.Context, param *DeleteCacheRequest) ([]*DeleteCacheResult, error) {
	url := o.Client.RootURL() + "deletecache"

	// do request
	data, err := o.Client.Do(ctx, "POST", url, param)
	if err != nil {
		return nil, err
	}

	// build results
	type deleteCacheResult struct {
		Results []*DeleteCacheResult
	}
	var results deleteCacheResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return results.Results, nil
}

// MonthlyUsage クラウドアカウントに登録されている全サイトの月別使用量を取得する。
//
// targetフィールドの値は「yyyymm」形式で対象年月を指定する。
// (例: 2021年02月の場合は、「202102」と指定。)
// 指定がない場合は、今月の月別使用量を取得する。
func (o *Op) MonthlyUsage(ctx context.Context, targetYM string) (*MonthlyUsageResults, error) {
	params := url.Values{
		"target": {targetYM},
	}
	url := o.Client.RootURL() + "monthlyusage?" + params.Encode()

	// do request
	data, err := o.Client.Do(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// build results
	var results MonthlyUsageResults
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return &results, nil
}

// Delete サイト削除
func (o *Op) Delete(ctx context.Context, id string) (*Site, error) {
	url := o.Client.RootURL() + fmt.Sprintf("site/%s", id)

	// build request body
	var body interface{}

	// do request
	data, err := o.Client.Do(ctx, "DELETE", url, body)
	if err != nil {
		return nil, err
	}

	// build results
	type deleteResults struct {
		Site *Site
	}
	var results deleteResults
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return results.Site, nil
}

// ApplyLogUploadConfig ログアップロード設定を作成・更新
func (o *Op) ApplyLogUploadConfig(ctx context.Context, id string, param *LogUploadConfig) (*LogUploadConfig, error) {
	url := o.Client.RootURL() + fmt.Sprintf("site/%s/log-upload-config", id)

	type applyLogUploadConfigRequest struct {
		SiteLogUploadConfig *LogUploadConfig `json:",omitempty"`
	}
	req := applyLogUploadConfigRequest{
		SiteLogUploadConfig: param,
	}
	// do request
	data, err := o.Client.Do(ctx, "PUT", url, req)
	if err != nil {
		return nil, err
	}

	// build results
	type applyLogUploadConfigResponse applyLogUploadConfigRequest
	var result applyLogUploadConfigResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result.SiteLogUploadConfig, nil
}

// ReadLogUploadConfig ログアップロードを取得
func (o *Op) ReadLogUploadConfig(ctx context.Context, id string) (*LogUploadConfig, error) {
	url := o.Client.RootURL() + fmt.Sprintf("site/%s/log-upload-config", id)

	var body interface{}

	// do request
	data, err := o.Client.Do(ctx, "GET", url, body)
	if err != nil {
		return nil, err
	}

	// build results
	type ReadLogUploadConfigResponse struct {
		SiteLogUploadConfig *LogUploadConfig `json:",omitempty"`
	}

	var result ReadLogUploadConfigResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result.SiteLogUploadConfig, nil
}

// DeleteLogUploadConfig ログアップロードを削除
func (o *Op) DeleteLogUploadConfig(ctx context.Context, id string) error {
	url := o.Client.RootURL() + fmt.Sprintf("site/%s/log-upload-config", id)

	var body interface{}

	// do request
	_, err := o.Client.Do(ctx, "DELETE", url, body)
	if err != nil {
		return err
	}

	return nil
}

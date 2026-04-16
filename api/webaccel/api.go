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

import "context"

// API is interface for operate WebAccel resource
type API interface {
	List(ctx context.Context) (*ListSitesResult, error)
	Create(ctx context.Context, param *CreateSiteRequest) (*Site, error)
	Read(ctx context.Context, id string) (*Site, error)
	Update(ctx context.Context, id string, param *UpdateSiteRequest) (*Site, error)
	UpdateStatus(ctx context.Context, id string, param *UpdateSiteStatusRequest) (*Site, error)
	CreateOriginGuardToken(ctx context.Context, id string) (*OriginGuardTokenResponse, error)
	DeleteOriginGuardToken(ctx context.Context, id string) error
	CreateNextOriginGuardToken(ctx context.Context, id string) (*OriginGuardTokenResponse, error)
	DeleteNextOriginGuardToken(ctx context.Context, id string) error
	CreateAutoCertUpdate(ctx context.Context, id string) error
	DeleteAutoCertUpdate(ctx context.Context, id string) error
	ReadACL(ctx context.Context, id string) (*ACLResult, error)
	UpsertACL(ctx context.Context, id string, acl string) (*ACLResult, error)
	DeleteACL(ctx context.Context, id string) error
	ReadCertificate(ctx context.Context, id string) (*Certificates, error)
	CreateCertificate(ctx context.Context, id string, param *CreateOrUpdateCertificateRequest) (*Certificates, error)
	UpdateCertificate(ctx context.Context, id string, param *CreateOrUpdateCertificateRequest) (*Certificates, error)
	DeleteCertificate(ctx context.Context, id string) error
	DeleteAllCache(ctx context.Context, param *DeleteAllCacheRequest) error
	DeleteCache(ctx context.Context, param *DeleteCacheRequest) ([]*DeleteCacheResult, error)
	Delete(ctx context.Context, id string) (*Site, error)
	MonthlyUsage(ctx context.Context, targetYM string) (*MonthlyUsageResults, error)
	ApplyLogUploadConfig(ctx context.Context, id string, param *LogUploadConfig) (*LogUploadConfig, error)
	ReadLogUploadConfig(ctx context.Context, id string) (*LogUploadConfig, error)
	DeleteLogUploadConfig(ctx context.Context, id string) error
}

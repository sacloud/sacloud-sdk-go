// Copyright 2025- The sacloud/iam-api-go authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// without any warranties or conditions of any kind.

package serviceprincipal

import (
	"context"

	"github.com/google/uuid"
	v1 "github.com/sacloud/sacloud-sdk-go/api/iam/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/api/iam/common"
	"github.com/sacloud/sacloud-sdk-go/common/packages/into"
)

type ServicePrincipalAPI interface {
	List(ctx context.Context, params ListParams) (*v1.ListServicePrincipalsOK, error)
	Create(ctx context.Context, params CreateParams) (*v1.ServicePrincipal, error)
	Read(ctx context.Context, id int) (*v1.ServicePrincipal, error)
	Update(ctx context.Context, id int, params UpdateParams) (*v1.ServicePrincipal, error)
	Delete(ctx context.Context, id int) error

	ListKeys(ctx context.Context, id int, params ListKeysParams) (*v1.ListServicePrincipalKeysOK, error)
	UploadKey(ctx context.Context, id int, publicKey v1.ServiceprincipalKeyPublicKey) (*v1.ServicePrincipalKey, error)
	EnableKey(ctx context.Context, id int, keyID uuid.UUID) (*v1.ServicePrincipalKey, error)
	DisableKey(ctx context.Context, id int, keyID uuid.UUID) (*v1.ServicePrincipalKey, error)
	DeleteKey(ctx context.Context, id int, keyID uuid.UUID) error

	IssueToken(ctx context.Context, assertion string) (*v1.ServicePrincipalOAuth2AccessToken, error)
}

type servicePrincipalOp struct {
	client *v1.Client
}

func NewServicePrincipalOp(client *v1.Client) ServicePrincipalAPI {
	return &servicePrincipalOp{client: client}
}

type ListParams struct {
	Page      *int
	PerPage   *int
	ProjectID *int
	Ordering  *v1.ListServicePrincipalsOrdering
}

func (s *servicePrincipalOp) List(ctx context.Context, params ListParams) (*v1.ListServicePrincipalsOK, error) {
	return common.ErrorFromDecodedResponse[v1.ListServicePrincipalsOK]("ServicePrincipal.List", func() (any, error) {
		return s.client.ListServicePrincipals(ctx, v1.ListServicePrincipalsParams{
			Page:      into.Opt[v1.OptInt](params.Page),
			PerPage:   into.Opt[v1.OptInt](params.PerPage),
			ProjectID: into.Opt[v1.OptInt](params.ProjectID),
			Ordering:  into.Opt[v1.OptListServicePrincipalsOrdering](params.Ordering),
		})
	})
}

type CreateParams = v1.CreateServicePrincipalReq

func (s *servicePrincipalOp) Create(ctx context.Context, params CreateParams) (*v1.ServicePrincipal, error) {
	return common.ErrorFromDecodedResponse[v1.ServicePrincipal]("ServicePrincipal.Create", func() (any, error) {
		return s.client.CreateServicePrincipal(ctx, &params)
	})
}

func (s *servicePrincipalOp) Read(ctx context.Context, id int) (*v1.ServicePrincipal, error) {
	return common.ErrorFromDecodedResponse[v1.ServicePrincipal]("ServicePrincipal.Read", func() (any, error) {
		return s.client.ReadServicePrincipal(ctx, v1.ReadServicePrincipalParams{ServicePrincipalID: id})
	})
}

type UpdateParams = v1.UpdateServicePrincipalReq

func (s *servicePrincipalOp) Update(ctx context.Context, id int, params UpdateParams) (*v1.ServicePrincipal, error) {
	return common.ErrorFromDecodedResponse[v1.ServicePrincipal]("ServicePrincipal.Update", func() (any, error) {
		return s.client.UpdateServicePrincipal(ctx, &params, v1.UpdateServicePrincipalParams{ServicePrincipalID: id})
	})
}

func (s *servicePrincipalOp) Delete(ctx context.Context, id int) error {
	_, err := common.ErrorFromDecodedResponse[v1.DeleteServicePrincipalNoContent]("ServicePrincipal.Delete", func() (any, error) {
		return s.client.DeleteServicePrincipal(ctx, v1.DeleteServicePrincipalParams{ServicePrincipalID: id})
	})
	return err
}

type ListKeysParams struct {
	Page     *int
	PerPage  *int
	Ordering *v1.ListServicePrincipalKeysOrdering
}

func (s *servicePrincipalOp) ListKeys(ctx context.Context, id int, params ListKeysParams) (*v1.ListServicePrincipalKeysOK, error) {
	return common.ErrorFromDecodedResponse[v1.ListServicePrincipalKeysOK]("ServicePrincipal.ListKeys", func() (any, error) {
		return s.client.ListServicePrincipalKeys(ctx, v1.ListServicePrincipalKeysParams{
			ServicePrincipalID: id,
			Page:               into.Opt[v1.OptInt](params.Page),
			PerPage:            into.Opt[v1.OptInt](params.PerPage),
			Ordering:           into.Opt[v1.OptListServicePrincipalKeysOrdering](params.Ordering),
		})
	})
}

func (s *servicePrincipalOp) UploadKey(ctx context.Context, id int, publicKey v1.ServiceprincipalKeyPublicKey) (*v1.ServicePrincipalKey, error) {
	return common.ErrorFromDecodedResponse[v1.ServicePrincipalKey]("ServicePrincipal.UploadKey", func() (any, error) {
		request := v1.NewOptUploadServicePrincipalKeyReq(v1.UploadServicePrincipalKeyReq{PublicKey: publicKey})
		return s.client.UploadServicePrincipalKey(ctx, request, v1.UploadServicePrincipalKeyParams{ServicePrincipalID: id})
	})
}

func (s *servicePrincipalOp) EnableKey(ctx context.Context, id int, keyID uuid.UUID) (*v1.ServicePrincipalKey, error) {
	return common.ErrorFromDecodedResponse[v1.ServicePrincipalKey]("ServicePrincipal.EnableKey", func() (any, error) {
		return s.client.EnableServicePrincipalKey(ctx, v1.EnableServicePrincipalKeyParams{
			ServicePrincipalID:    id,
			ServicePrincipalKeyID: keyID,
		})
	})
}

func (s *servicePrincipalOp) DisableKey(ctx context.Context, id int, keyID uuid.UUID) (*v1.ServicePrincipalKey, error) {
	return common.ErrorFromDecodedResponse[v1.ServicePrincipalKey]("ServicePrincipal.DisableKey", func() (any, error) {
		return s.client.DisableServicePrincipalKey(ctx, v1.DisableServicePrincipalKeyParams{
			ServicePrincipalID:    id,
			ServicePrincipalKeyID: keyID,
		})
	})
}

func (s *servicePrincipalOp) DeleteKey(ctx context.Context, id int, keyID uuid.UUID) error {
	_, err := common.ErrorFromDecodedResponse[v1.DeleteServicePrincipalKeyNoContent]("ServicePrincipal.DeleteKey", func() (any, error) {
		return s.client.DeleteServicePrincipalKey(ctx, v1.DeleteServicePrincipalKeyParams{
			ServicePrincipalID:    id,
			ServicePrincipalKeyID: keyID,
		})
	})
	return err
}

func (s *servicePrincipalOp) IssueToken(ctx context.Context, assertion string) (*v1.ServicePrincipalOAuth2AccessToken, error) {
	return common.ErrorFromDecodedResponse[v1.ServicePrincipalOAuth2AccessToken]("ServicePrincipal.IssueToken", func() (any, error) {
		return s.client.IssueServicePrincipalToken(ctx, &v1.ServicePrincipalJWTGrantRequest{
			GrantType: v1.ServicePrincipalJWTGrantRequestGrantTypeUrnIetfParamsOAuthGrantTypeJwtBearer,
			Assertion: assertion,
		})
	})
}

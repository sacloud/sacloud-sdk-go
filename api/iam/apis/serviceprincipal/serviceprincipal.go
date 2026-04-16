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
	v1 "github.com/sacloud/iam-api-go/apis/v1"
	"github.com/sacloud/iam-api-go/common"
)

type ServicePrincipalAPI interface {
	List(ctx context.Context, params ListParams) (*v1.ServicePrincipalsGetOK, error)
	Create(ctx context.Context, params CreateParams) (*v1.ServicePrincipal, error)
	Read(ctx context.Context, id int) (*v1.ServicePrincipal, error)
	Update(ctx context.Context, id int, params UpdateParams) (*v1.ServicePrincipal, error)
	Delete(ctx context.Context, id int) error

	ListKeys(ctx context.Context, id int, params ListKeysParams) (*v1.ServicePrincipalsServicePrincipalIDKeysGetOK, error)
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
	Ordering  *v1.ServicePrincipalsGetOrdering
}

func (s *servicePrincipalOp) List(ctx context.Context, params ListParams) (*v1.ServicePrincipalsGetOK, error) {
	return common.ErrorFromDecodedResponse[v1.ServicePrincipalsGetOK]("ServicePrincipal.List", func() (any, error) {
		return s.client.ServicePrincipalsGet(ctx, v1.ServicePrincipalsGetParams{
			Page:      common.IntoOpt[v1.OptInt](params.Page),
			PerPage:   common.IntoOpt[v1.OptInt](params.PerPage),
			ProjectID: common.IntoOpt[v1.OptInt](params.ProjectID),
			Ordering:  common.IntoOpt[v1.OptServicePrincipalsGetOrdering](params.Ordering),
		})
	})
}

type CreateParams = v1.ServicePrincipalsPostReq

func (s *servicePrincipalOp) Create(ctx context.Context, params CreateParams) (*v1.ServicePrincipal, error) {
	return common.ErrorFromDecodedResponse[v1.ServicePrincipal]("ServicePrincipal.Create", func() (any, error) {
		return s.client.ServicePrincipalsPost(ctx, &params)
	})
}

func (s *servicePrincipalOp) Read(ctx context.Context, id int) (*v1.ServicePrincipal, error) {
	return common.ErrorFromDecodedResponse[v1.ServicePrincipal]("ServicePrincipal.Read", func() (any, error) {
		return s.client.ServicePrincipalsServicePrincipalIDGet(ctx, v1.ServicePrincipalsServicePrincipalIDGetParams{ServicePrincipalID: id})
	})
}

type UpdateParams = v1.ServicePrincipalsServicePrincipalIDPutReq

func (s *servicePrincipalOp) Update(ctx context.Context, id int, params UpdateParams) (*v1.ServicePrincipal, error) {
	return common.ErrorFromDecodedResponse[v1.ServicePrincipal]("ServicePrincipal.Update", func() (any, error) {
		return s.client.ServicePrincipalsServicePrincipalIDPut(ctx, &params, v1.ServicePrincipalsServicePrincipalIDPutParams{ServicePrincipalID: id})
	})
}

func (s *servicePrincipalOp) Delete(ctx context.Context, id int) error {
	_, err := common.ErrorFromDecodedResponse[v1.ServicePrincipalsServicePrincipalIDDeleteNoContent]("ServicePrincipal.Delete", func() (any, error) {
		return s.client.ServicePrincipalsServicePrincipalIDDelete(ctx, v1.ServicePrincipalsServicePrincipalIDDeleteParams{ServicePrincipalID: id})
	})
	return err
}

type ListKeysParams struct {
	Page     *int
	PerPage  *int
	Ordering *v1.ServicePrincipalsServicePrincipalIDKeysGetOrdering
}

func (s *servicePrincipalOp) ListKeys(ctx context.Context, id int, params ListKeysParams) (*v1.ServicePrincipalsServicePrincipalIDKeysGetOK, error) {
	return common.ErrorFromDecodedResponse[v1.ServicePrincipalsServicePrincipalIDKeysGetOK]("ServicePrincipal.ListKeys", func() (any, error) {
		return s.client.ServicePrincipalsServicePrincipalIDKeysGet(ctx, v1.ServicePrincipalsServicePrincipalIDKeysGetParams{
			ServicePrincipalID: id,
			Page:               common.IntoOpt[v1.OptInt](params.Page),
			PerPage:            common.IntoOpt[v1.OptInt](params.PerPage),
			Ordering:           common.IntoOpt[v1.OptServicePrincipalsServicePrincipalIDKeysGetOrdering](params.Ordering),
		})
	})
}

func (s *servicePrincipalOp) UploadKey(ctx context.Context, id int, publicKey v1.ServiceprincipalKeyPublicKey) (*v1.ServicePrincipalKey, error) {
	return common.ErrorFromDecodedResponse[v1.ServicePrincipalKey]("ServicePrincipal.UploadKey", func() (any, error) {
		request := v1.NewOptServicePrincipalsServicePrincipalIDUploadKeyPostReq(v1.ServicePrincipalsServicePrincipalIDUploadKeyPostReq{PublicKey: publicKey})
		params := v1.ServicePrincipalsServicePrincipalIDUploadKeyPostParams{ServicePrincipalID: id}
		return s.client.ServicePrincipalsServicePrincipalIDUploadKeyPost(ctx, request, params)
	})
}

func (s *servicePrincipalOp) EnableKey(ctx context.Context, id int, keyID uuid.UUID) (*v1.ServicePrincipalKey, error) {
	return common.ErrorFromDecodedResponse[v1.ServicePrincipalKey]("ServicePrincipal.EnableKey", func() (any, error) {
		return s.client.ServicePrincipalsServicePrincipalIDKeysServicePrincipalKeyIDEnablePost(ctx, v1.ServicePrincipalsServicePrincipalIDKeysServicePrincipalKeyIDEnablePostParams{
			ServicePrincipalID:    id,
			ServicePrincipalKeyID: keyID,
		})
	})
}

func (s *servicePrincipalOp) DisableKey(ctx context.Context, id int, keyID uuid.UUID) (*v1.ServicePrincipalKey, error) {
	return common.ErrorFromDecodedResponse[v1.ServicePrincipalKey]("ServicePrincipal.DisableKey", func() (any, error) {
		return s.client.ServicePrincipalsServicePrincipalIDKeysServicePrincipalKeyIDDisablePost(ctx, v1.ServicePrincipalsServicePrincipalIDKeysServicePrincipalKeyIDDisablePostParams{
			ServicePrincipalID:    id,
			ServicePrincipalKeyID: keyID,
		})
	})
}

func (s *servicePrincipalOp) DeleteKey(ctx context.Context, id int, keyID uuid.UUID) error {
	_, err := common.ErrorFromDecodedResponse[v1.ServicePrincipalsServicePrincipalIDKeysServicePrincipalKeyIDDeleteNoContent]("ServicePrincipal.DeleteKey", func() (any, error) {
		return s.client.ServicePrincipalsServicePrincipalIDKeysServicePrincipalKeyIDDelete(ctx, v1.ServicePrincipalsServicePrincipalIDKeysServicePrincipalKeyIDDeleteParams{
			ServicePrincipalID:    id,
			ServicePrincipalKeyID: keyID,
		})
	})
	return err
}

func (s *servicePrincipalOp) IssueToken(ctx context.Context, assertion string) (*v1.ServicePrincipalOAuth2AccessToken, error) {
	return common.ErrorFromDecodedResponse[v1.ServicePrincipalOAuth2AccessToken]("ServicePrincipal.IssueToken", func() (any, error) {
		return s.client.ServicePrincipalsOAuth2TokenPost(ctx, &v1.ServicePrincipalJWTGrantRequest{
			GrantType: v1.ServicePrincipalJWTGrantRequestGrantTypeUrnIetfParamsOAuthGrantTypeJwtBearer,
			Assertion: assertion,
		})
	})
}

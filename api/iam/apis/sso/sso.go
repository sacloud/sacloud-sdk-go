// Copyright 2025- The sacloud/iam-api-go authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is provided on an "AS IS" basis,
// without any warranties or conditions of any kind.

package sso

import (
	"context"

	v1 "github.com/sacloud/iam-api-go/apis/v1"
	"github.com/sacloud/iam-api-go/common"
)

type SSOAPI interface {
	List(ctx context.Context, page, perPage *int) (*v1.SSOProfilesGetOK, error)
	Create(ctx context.Context, params CreateParams) (*v1.SSOProfile, error)
	Read(ctx context.Context, id int) (*v1.SSOProfile, error)
	Update(ctx context.Context, id int, params UpdateParams) (*v1.SSOProfile, error)
	Delete(ctx context.Context, id int) error

	Link(ctx context.Context, id int) (*v1.SSOProfile, error)
	Unlink(ctx context.Context, id int) (*v1.SSOProfile, error)
}

type ssoOp struct {
	client *v1.Client
}

func NewSSOOp(client *v1.Client) SSOAPI { return &ssoOp{client: client} }

func (s *ssoOp) List(ctx context.Context, page, perPage *int) (*v1.SSOProfilesGetOK, error) {
	return common.ErrorFromDecodedResponse[v1.SSOProfilesGetOK]("SSO.List", func() (any, error) {
		return s.client.SSOProfilesGet(ctx, v1.SSOProfilesGetParams{
			Page:    common.IntoOpt[v1.OptInt](page),
			PerPage: common.IntoOpt[v1.OptInt](perPage),
		})
	})
}

type CreateParams = v1.SSOProfilesPostReq

func (s *ssoOp) Create(ctx context.Context, params CreateParams) (*v1.SSOProfile, error) {
	return common.ErrorFromDecodedResponse[v1.SSOProfile]("SSO.Create", func() (any, error) {
		return s.client.SSOProfilesPost(ctx, &params)
	})
}

func (s *ssoOp) Read(ctx context.Context, id int) (*v1.SSOProfile, error) {
	return common.ErrorFromDecodedResponse[v1.SSOProfile]("SSO.Read", func() (any, error) {
		return s.client.SSOProfilesSSOProfileIDGet(ctx, v1.SSOProfilesSSOProfileIDGetParams{SSOProfileID: id})
	})
}

type UpdateParams = v1.SSOProfilesSSOProfileIDPutReq

func (s *ssoOp) Update(ctx context.Context, id int, params UpdateParams) (*v1.SSOProfile, error) {
	return common.ErrorFromDecodedResponse[v1.SSOProfile]("SSO.Update", func() (any, error) {
		return s.client.SSOProfilesSSOProfileIDPut(ctx, &params, v1.SSOProfilesSSOProfileIDPutParams{SSOProfileID: id})
	})
}

func (s *ssoOp) Delete(ctx context.Context, id int) error {
	_, err := common.ErrorFromDecodedResponse[v1.SSOProfilesSSOProfileIDDeleteNoContent]("SSO.Delete", func() (any, error) {
		return s.client.SSOProfilesSSOProfileIDDelete(ctx, v1.SSOProfilesSSOProfileIDDeleteParams{SSOProfileID: id})
	})
	return err
}

func (s *ssoOp) Link(ctx context.Context, id int) (*v1.SSOProfile, error) {
	return common.ErrorFromDecodedResponse[v1.SSOProfile]("SSO.Link", func() (any, error) {
		return s.client.SSOProfilesSSOProfileIDAssignPost(ctx, v1.SSOProfilesSSOProfileIDAssignPostParams{SSOProfileID: id})
	})
}

func (s *ssoOp) Unlink(ctx context.Context, id int) (*v1.SSOProfile, error) {
	return common.ErrorFromDecodedResponse[v1.SSOProfile]("SSO.Unlink", func() (any, error) {
		return s.client.SSOProfilesSSOProfileIDUnassignPost(ctx, v1.SSOProfilesSSOProfileIDUnassignPostParams{SSOProfileID: id})
	})
}

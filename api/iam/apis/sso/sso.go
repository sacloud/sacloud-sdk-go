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

	v1 "github.com/sacloud/sacloud-sdk-go/api/iam/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/api/iam/common"
	"github.com/sacloud/sacloud-sdk-go/common/packages/into"
)

type SSOAPI interface {
	List(ctx context.Context, page, perPage *int) (*v1.ListSsoProfilesOK, error)
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

func (s *ssoOp) List(ctx context.Context, page, perPage *int) (*v1.ListSsoProfilesOK, error) {
	return common.ErrorFromDecodedResponse[v1.ListSsoProfilesOK]("SSO.List", func() (any, error) {
		return s.client.ListSsoProfiles(ctx, v1.ListSsoProfilesParams{
			Page:    into.Opt[v1.OptInt](page),
			PerPage: into.Opt[v1.OptInt](perPage),
		})
	})
}

type CreateParams = v1.CreateSsoProfileReq

func (s *ssoOp) Create(ctx context.Context, params CreateParams) (*v1.SSOProfile, error) {
	return common.ErrorFromDecodedResponse[v1.SSOProfile]("SSO.Create", func() (any, error) {
		return s.client.CreateSsoProfile(ctx, &params)
	})
}

func (s *ssoOp) Read(ctx context.Context, id int) (*v1.SSOProfile, error) {
	return common.ErrorFromDecodedResponse[v1.SSOProfile]("SSO.Read", func() (any, error) {
		return s.client.ReadSsoProfile(ctx, v1.ReadSsoProfileParams{SSOProfileID: id})
	})
}

type UpdateParams = v1.UpdateSsoProfileReq

func (s *ssoOp) Update(ctx context.Context, id int, params UpdateParams) (*v1.SSOProfile, error) {
	return common.ErrorFromDecodedResponse[v1.SSOProfile]("SSO.Update", func() (any, error) {
		return s.client.UpdateSsoProfile(ctx, &params, v1.UpdateSsoProfileParams{SSOProfileID: id})
	})
}

func (s *ssoOp) Delete(ctx context.Context, id int) error {
	_, err := common.ErrorFromDecodedResponse[v1.DeleteSsoProfileNoContent]("SSO.Delete", func() (any, error) {
		return s.client.DeleteSsoProfile(ctx, v1.DeleteSsoProfileParams{SSOProfileID: id})
	})
	return err
}

func (s *ssoOp) Link(ctx context.Context, id int) (*v1.SSOProfile, error) {
	return common.ErrorFromDecodedResponse[v1.SSOProfile]("SSO.Link", func() (any, error) {
		return s.client.AssignSsoProfile(ctx, v1.AssignSsoProfileParams{SSOProfileID: id})
	})
}

func (s *ssoOp) Unlink(ctx context.Context, id int) (*v1.SSOProfile, error) {
	return common.ErrorFromDecodedResponse[v1.SSOProfile]("SSO.Unlink", func() (any, error) {
		return s.client.UnassignSsoProfile(ctx, v1.UnassignSsoProfileParams{SSOProfileID: id})
	})
}

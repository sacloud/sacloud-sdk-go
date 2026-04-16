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
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package user2fa

import (
	"context"

	v1 "github.com/sacloud/iam-api-go/apis/v1"
	"github.com/sacloud/iam-api-go/common"
)

type User2FAAPI interface {
	DeactivateOTP(ctx context.Context) error

	ListTrustedDevices(ctx context.Context) (*v1.CompatUsersUserIDTrustedDevicesGetOK, error)
	DeleteTrustedDevice(ctx context.Context, trustedDeviceID int) error
	ClearTrustedDevices(ctx context.Context) error

	ListSecurityKeys(ctx context.Context) (*v1.CompatUsersUserIDSecurityKeysGetOK, error)
	ReadSecurityKey(ctx context.Context, securityKeyID int) (*v1.UserSecurityKey, error)
	UpdateSecurityKey(ctx context.Context, securityKeyID int, name string) (*v1.UserSecurityKey, error)
	DeleteSecurityKey(ctx context.Context, securityKeyID int) error
}

type user2faOp struct {
	client *v1.Client
	user   *v1.User
}

func NewUser2FAOp(client *v1.Client, user *v1.User) User2FAAPI {
	return &user2faOp{
		client: client,
		user:   user,
	}
}

func (u *user2faOp) getUserID() int { return u.user.GetID() }

func (u *user2faOp) DeactivateOTP(ctx context.Context) error {
	_, err := common.ErrorFromDecodedResponse[v1.CompatUsersUserIDDeactivateOtpPostNoContent]("User2FA.DeactivateOTP", func() (any, error) {
		return u.client.CompatUsersUserIDDeactivateOtpPost(ctx, v1.CompatUsersUserIDDeactivateOtpPostParams{UserID: u.getUserID()})
	})
	return err
}

func (u *user2faOp) ListTrustedDevices(ctx context.Context) (*v1.CompatUsersUserIDTrustedDevicesGetOK, error) {
	return common.ErrorFromDecodedResponse[v1.CompatUsersUserIDTrustedDevicesGetOK]("User2FA.ListTrustedDevices", func() (any, error) {
		return u.client.CompatUsersUserIDTrustedDevicesGet(ctx, v1.CompatUsersUserIDTrustedDevicesGetParams{UserID: u.getUserID()})
	})
}

func (u *user2faOp) DeleteTrustedDevice(ctx context.Context, trustedDeviceID int) error {
	_, err := common.ErrorFromDecodedResponse[v1.CompatUsersUserIDTrustedDevicesTrustedDeviceIDDeleteNoContent]("User2FA.DeleteTrustedDevice", func() (any, error) {
		return u.client.CompatUsersUserIDTrustedDevicesTrustedDeviceIDDelete(ctx, v1.CompatUsersUserIDTrustedDevicesTrustedDeviceIDDeleteParams{
			UserID:          u.getUserID(),
			TrustedDeviceID: trustedDeviceID,
		})
	})
	return err
}

func (u *user2faOp) ClearTrustedDevices(ctx context.Context) error {
	_, err := common.ErrorFromDecodedResponse[v1.CompatUsersUserIDClearTrustedDevicesPostNoContent]("User2FA.ClearTrustedDevices", func() (any, error) {
		return u.client.CompatUsersUserIDClearTrustedDevicesPost(ctx, v1.CompatUsersUserIDClearTrustedDevicesPostParams{UserID: u.getUserID()})
	})
	return err
}

func (u *user2faOp) ListSecurityKeys(ctx context.Context) (*v1.CompatUsersUserIDSecurityKeysGetOK, error) {
	return common.ErrorFromDecodedResponse[v1.CompatUsersUserIDSecurityKeysGetOK]("User2FA.ListSecurityKeys", func() (any, error) {
		return u.client.CompatUsersUserIDSecurityKeysGet(ctx, v1.CompatUsersUserIDSecurityKeysGetParams{UserID: u.getUserID()})
	})
}

func (u *user2faOp) ReadSecurityKey(ctx context.Context, securityKeyID int) (*v1.UserSecurityKey, error) {
	return common.ErrorFromDecodedResponse[v1.UserSecurityKey]("User2FA.ReadSecurityKey", func() (any, error) {
		return u.client.CompatUsersUserIDSecurityKeysSecurityKeyIDGet(ctx, v1.CompatUsersUserIDSecurityKeysSecurityKeyIDGetParams{
			UserID:        u.getUserID(),
			SecurityKeyID: securityKeyID,
		})
	})
}

func (u *user2faOp) UpdateSecurityKey(ctx context.Context, securityKeyID int, name string) (*v1.UserSecurityKey, error) {
	return common.ErrorFromDecodedResponse[v1.UserSecurityKey]("User2FA.UpdateSecurityKey", func() (any, error) {
		req := v1.NewOptCompatUsersUserIDSecurityKeysSecurityKeyIDPutReq(v1.CompatUsersUserIDSecurityKeysSecurityKeyIDPutReq{Name: name})
		params := v1.CompatUsersUserIDSecurityKeysSecurityKeyIDPutParams{
			UserID:        u.getUserID(),
			SecurityKeyID: securityKeyID,
		}
		return u.client.CompatUsersUserIDSecurityKeysSecurityKeyIDPut(ctx, req, params)
	})
}

func (u *user2faOp) DeleteSecurityKey(ctx context.Context, securityKeyID int) error {
	_, err := common.ErrorFromDecodedResponse[v1.CompatUsersUserIDSecurityKeysSecurityKeyIDDeleteNoContent]("User2FA.DeleteSecurityKey", func() (any, error) {
		return u.client.CompatUsersUserIDSecurityKeysSecurityKeyIDDelete(ctx, v1.CompatUsersUserIDSecurityKeysSecurityKeyIDDeleteParams{
			UserID:        u.getUserID(),
			SecurityKeyID: securityKeyID,
		})
	})
	return err
}

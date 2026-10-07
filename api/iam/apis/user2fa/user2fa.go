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

	v1 "github.com/sacloud/sacloud-sdk-go/api/iam/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/api/iam/common"
)

type User2FAAPI interface {
	DeactivateOTP(ctx context.Context) error

	ListTrustedDevices(ctx context.Context) (*v1.ListTrustedDevicesOK, error)
	DeleteTrustedDevice(ctx context.Context, trustedDeviceID int) error
	ClearTrustedDevices(ctx context.Context) error

	ListSecurityKeys(ctx context.Context) (*v1.ListSecurityKeysOK, error)
	ReadSecurityKey(ctx context.Context, securityKeyID int) (*v1.UserSecurityKey, error)
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
	_, err := common.ErrorFromDecodedResponse[v1.DeactivateOtpNoContent]("User2FA.DeactivateOTP", func() (any, error) {
		return u.client.DeactivateOtp(ctx, v1.DeactivateOtpParams{UserID: u.getUserID()})
	})
	return err
}

func (u *user2faOp) ListTrustedDevices(ctx context.Context) (*v1.ListTrustedDevicesOK, error) {
	return common.ErrorFromDecodedResponse[v1.ListTrustedDevicesOK]("User2FA.ListTrustedDevices", func() (any, error) {
		return u.client.ListTrustedDevices(ctx, v1.ListTrustedDevicesParams{UserID: u.getUserID()})
	})
}

func (u *user2faOp) DeleteTrustedDevice(ctx context.Context, trustedDeviceID int) error {
	_, err := common.ErrorFromDecodedResponse[v1.DeleteTrustedDeviceNoContent]("User2FA.DeleteTrustedDevice", func() (any, error) {
		return u.client.DeleteTrustedDevice(ctx, v1.DeleteTrustedDeviceParams{
			UserID:          u.getUserID(),
			TrustedDeviceID: trustedDeviceID,
		})
	})
	return err
}

func (u *user2faOp) ClearTrustedDevices(ctx context.Context) error {
	_, err := common.ErrorFromDecodedResponse[v1.ClearTrustedDevicesNoContent]("User2FA.ClearTrustedDevices", func() (any, error) {
		return u.client.ClearTrustedDevices(ctx, v1.ClearTrustedDevicesParams{UserID: u.getUserID()})
	})
	return err
}

func (u *user2faOp) ListSecurityKeys(ctx context.Context) (*v1.ListSecurityKeysOK, error) {
	return common.ErrorFromDecodedResponse[v1.ListSecurityKeysOK]("User2FA.ListSecurityKeys", func() (any, error) {
		return u.client.ListSecurityKeys(ctx, v1.ListSecurityKeysParams{UserID: u.getUserID()})
	})
}

func (u *user2faOp) ReadSecurityKey(ctx context.Context, securityKeyID int) (*v1.UserSecurityKey, error) {
	return common.ErrorFromDecodedResponse[v1.UserSecurityKey]("User2FA.ReadSecurityKey", func() (any, error) {
		return u.client.ReadSecurityKey(ctx, v1.ReadSecurityKeyParams{
			UserID:        u.getUserID(),
			SecurityKeyID: securityKeyID,
		})
	})
}

func (u *user2faOp) DeleteSecurityKey(ctx context.Context, securityKeyID int) error {
	_, err := common.ErrorFromDecodedResponse[v1.DeleteSecurityKeyNoContent]("User2FA.DeleteSecurityKey", func() (any, error) {
		return u.client.DeleteSecurityKey(ctx, v1.DeleteSecurityKeyParams{
			UserID:        u.getUserID(),
			SecurityKeyID: securityKeyID,
		})
	})
	return err
}

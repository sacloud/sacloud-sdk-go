// Copyright 2022-2026 The object-storage-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package objectstorage

import (
	"context"

	v2 "github.com/sacloud/sacloud-sdk-go/api/object-storage/apis/v2"
)

type AccountAPI interface {
	Create(ctx context.Context) (*v2.AccountData, error)
	Read(ctx context.Context) (*v2.AccountData, error)
	Delete(ctx context.Context) error

	ListAccessKeys(ctx context.Context) ([]v2.AccountKeysDataItem, error)
	// Secretはこの戻り値でのみ参照可能
	CreateAccessKey(ctx context.Context) (*v2.AccountKeyData, error)
	// Secretは読み取りAPIのレスポンスに含まれない
	ReadAccessKey(ctx context.Context, keyId string) (*v2.ReadAccountKeyData, error)
	DeleteAccessKey(ctx context.Context, keyId string) error
}

var _ AccountAPI = (*accountOp)(nil)

type accountOp struct {
	client *SiteClient
}

func NewAccountOp(client *SiteClient) AccountAPI {
	return &accountOp{client: client}
}

func (op *accountOp) Create(ctx context.Context) (*v2.AccountData, error) {
	res, err := errorFromDecodedResponse[v2.Account]("Accounts.Create", func() (any, error) { return op.client.client.CreateAccount(ctx) })
	if err != nil {
		return nil, err
	}

	return new(res.GetData()), nil
}

func (op *accountOp) Read(ctx context.Context) (*v2.AccountData, error) {
	res, err := errorFromDecodedResponse[v2.Account]("Accounts.Read", func() (any, error) { return op.client.client.ReadAccount(ctx) })
	if err != nil {
		return nil, err
	}

	return new(res.GetData()), nil
}

func (op *accountOp) Delete(ctx context.Context) error {
	return errorFromDecodedResponseE[v2.Account]("Accounts.Delete", func() (any, error) { return op.client.client.DeleteAccount(ctx) })
}

func (op *accountOp) ListAccessKeys(ctx context.Context) ([]v2.AccountKeysDataItem, error) {
	res, err := errorFromDecodedResponse[v2.AccountKeys]("Accounts.ListAccessKeys", func() (any, error) { return op.client.client.ListAccountKeys(ctx) })
	if err != nil {
		return nil, err
	}

	return res.GetData(), nil
}

func (op *accountOp) CreateAccessKey(ctx context.Context) (*v2.AccountKeyData, error) {
	res, err := errorFromDecodedResponse[v2.AccountKey]("Accounts.CreateAccessKey", func() (any, error) { return op.client.client.CreateAccountKey(ctx) })
	if err != nil {
		return nil, err
	}

	return new(res.GetData()), nil
}

func (op *accountOp) ReadAccessKey(ctx context.Context, keyId string) (*v2.ReadAccountKeyData, error) {
	res, err := errorFromDecodedResponse[v2.ReadAccountKey]("Accounts.ReadAccessKey", func() (any, error) { return op.client.client.ReadAccountKey(ctx, v2.ReadAccountKeyParams{ID: keyId}) })
	if err != nil {
		return nil, err
	}

	return new(res.GetData()), nil
}

func (op *accountOp) DeleteAccessKey(ctx context.Context, keyId string) error {
	return errorFromDecodedResponseE[v2.DeleteAccountKeyNoContent]("Accounts.DeleteAccessKey", func() (any, error) {
		return op.client.client.DeleteAccountKey(ctx, v2.DeleteAccountKeyParams{ID: keyId})
	})
}

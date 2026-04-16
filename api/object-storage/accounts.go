// Copyright 2022-2026 The object-storage-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package objectstorage

import (
	"context"
	"errors"

	v2 "github.com/sacloud/object-storage-api-go/apis/v2"
)

type AccountAPI interface {
	Create(ctx context.Context) (*v2.AccountData, error)
	Read(ctx context.Context) (*v2.AccountData, error)
	Delete(ctx context.Context) error

	ListAccessKeys(ctx context.Context) ([]v2.AccountKeysDataItem, error)
	// Secretはこの戻り値でのみ参照可能
	CreateAccessKey(ctx context.Context) (*v2.AccountKeyData, error)
	// Secretは常に空文字になっている
	ReadAccessKey(ctx context.Context, keyId string) (*v2.AccountKeyData, error)
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
	res, err := op.client.client.CreateAccount(ctx)
	if err != nil {
		return nil, NewAPIError("Accounts.Create", 0, err)
	}

	switch r := res.(type) {
	case *v2.Account:
		return &r.Data.Value, nil
	case *v2.Error401:
		return nil, NewAPIError("Accounts.Create", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.Error403:
		return nil, NewAPIError("Accounts.Create", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.Error409:
		return nil, NewAPIError("Accounts.Create", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.ErrorDefaultStatusCode:
		return nil, NewAPIError("Accounts.Create", r.StatusCode, errors.New(string(r.Response.Error.Value.Message.Value)))
	default:
		return nil, NewAPIError("Accounts.Create", 0, errors.New("unknown error"))
	}
}

func (op *accountOp) Read(ctx context.Context) (*v2.AccountData, error) {
	res, err := op.client.client.GetAccount(ctx)
	if err != nil {
		return nil, NewAPIError("Accounts.Read", 0, err)
	}

	switch r := res.(type) {
	case *v2.Account:
		return &r.Data.Value, nil
	case *v2.Error401:
		return nil, NewAPIError("Accounts.Read", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.Error404:
		return nil, NewAPIError("Accounts.Read", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.ErrorDefaultStatusCode:
		return nil, NewAPIError("Accounts.Read", r.StatusCode, errors.New(string(r.Response.Error.Value.Message.Value)))
	default:
		return nil, NewAPIError("Accounts.Read", 0, errors.New("unknown error"))
	}
}

func (op *accountOp) Delete(ctx context.Context) error {
	res, err := op.client.client.DeleteAccount(ctx)
	if err != nil {
		return NewAPIError("Accounts.Delete", 0, err)
	}

	switch r := res.(type) {
	case *v2.DeleteAccountNoContent:
		return nil
	case *v2.Error401:
		return NewAPIError("Accounts.Delete", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.Error409:
		return NewAPIError("Accounts.Delete", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.ErrorDefaultStatusCode:
		return NewAPIError("Accounts.Delete", r.StatusCode, errors.New(string(r.Response.Error.Value.Message.Value)))
	default:
		return NewAPIError("Accounts.Delete", 0, errors.New("unknown error"))
	}
}

func (op *accountOp) ListAccessKeys(ctx context.Context) ([]v2.AccountKeysDataItem, error) {
	res, err := op.client.client.GetAccountKeys(ctx)
	if err != nil {
		return nil, NewAPIError("Accounts.ListAccessKeys", 0, err)
	}

	switch r := res.(type) {
	case *v2.AccountKeys:
		return r.Data, nil
	case *v2.Error401:
		return nil, NewAPIError("Accounts.ListAccessKeys", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.Error404:
		return nil, NewAPIError("Accounts.ListAccessKeys", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.ErrorDefaultStatusCode:
		return nil, NewAPIError("Accounts.ListAccessKeys", r.StatusCode, errors.New(string(r.Response.Error.Value.Message.Value)))
	default:
		return nil, NewAPIError("Accounts.ListAccessKeys", 0, errors.New("unknown error"))
	}
}

func (op *accountOp) CreateAccessKey(ctx context.Context) (*v2.AccountKeyData, error) {
	res, err := op.client.client.CreateAccountKey(ctx)
	if err != nil {
		return nil, NewAPIError("Accounts.CreateAccessKey", 0, err)
	}

	switch r := res.(type) {
	case *v2.AccountKey:
		return &r.Data.Value, nil
	case *v2.Error401:
		return nil, NewAPIError("Accounts.CreateAccessKey", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.Error404:
		return nil, NewAPIError("Accounts.CreateAccessKey", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.Error409:
		return nil, NewAPIError("Accounts.CreateAccessKey", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.ErrorDefaultStatusCode:
		return nil, NewAPIError("Accounts.CreateAccessKey", r.StatusCode, errors.New(string(r.Response.Error.Value.Message.Value)))
	default:
		return nil, NewAPIError("Accounts.CreateAccessKey", 0, errors.New("unknown error"))
	}
}

func (op *accountOp) ReadAccessKey(ctx context.Context, keyId string) (*v2.AccountKeyData, error) {
	res, err := op.client.client.GetAccountKey(ctx, v2.GetAccountKeyParams{ID: keyId})
	if err != nil {
		return nil, NewAPIError("Accounts.ReadAccessKey", 0, err)
	}

	switch r := res.(type) {
	case *v2.AccountKey:
		return &r.Data.Value, nil
	case *v2.Error401:
		return nil, NewAPIError("Accounts.ReadAccessKey", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.Error404:
		return nil, NewAPIError("Accounts.ReadAccessKey", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.ErrorDefaultStatusCode:
		return nil, NewAPIError("Accounts.ReadAccessKey", r.StatusCode, errors.New(string(r.Response.Error.Value.Message.Value)))
	default:
		return nil, NewAPIError("Accounts.ReadAccessKey", 0, errors.New("unknown error"))
	}
}

func (op *accountOp) DeleteAccessKey(ctx context.Context, keyId string) error {
	res, err := op.client.client.DeleteAccountKey(ctx, v2.DeleteAccountKeyParams{ID: keyId})
	if err != nil {
		return NewAPIError("Accounts.DeleteAccessKey", 0, err)
	}

	switch r := res.(type) {
	case *v2.DeleteAccountKeyNoContent:
		return nil
	case *v2.Error401:
		return NewAPIError("Accounts.DeleteAccessKey", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.ErrorDefaultStatusCode:
		return NewAPIError("Accounts.DeleteAccessKey", r.StatusCode, errors.New(string(r.Response.Error.Value.Message.Value)))
	default:
		return NewAPIError("Accounts.DeleteAccessKey", 0, errors.New("unknown error"))
	}
}

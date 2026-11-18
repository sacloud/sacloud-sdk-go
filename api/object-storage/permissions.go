// Copyright 2022-2026 The object-storage-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package objectstorage

import (
	"context"
	"errors"
	"strconv"

	v2 "github.com/sacloud/sacloud-sdk-go/api/object-storage/apis/v2"
)

type PermissionsAPI interface {
	List(ctx context.Context) ([]v2.PermissionsDataItem, error)
	Create(ctx context.Context, displayName string, controls v2.BucketControlsBody) (*v2.PermissionData, error)
	Read(ctx context.Context, permissionId string) (*v2.PermissionData, error)
	Update(ctx context.Context, permissionId string, displayName string, controls v2.BucketControlsBody) (*v2.PermissionData, error)
	Delete(ctx context.Context, permissionId string) error

	ListAccessKeys(ctx context.Context, permissionId string) ([]v2.PermissionKeyReadResponseBody, error)
	// Secretはこの戻り値でのみ参照可能
	CreateAccessKey(ctx context.Context, permissionId string) (*v2.PermissionKeyCreateResponseBody, error)
	// Secretは読み取りAPIのレスポンスに含まれない
	ReadAccessKey(ctx context.Context, permissionId string, accessKeyId string) (*v2.PermissionKeyReadResponseBody, error)
	DeleteAccessKey(ctx context.Context, permissionId string, accessKeyId string) error
}

var _ PermissionsAPI = (*permissionOp)(nil)

type permissionOp struct {
	client *SiteClient
}

// NewPermissionOp パーミッション関連API
func NewPermissionOp(client *SiteClient) PermissionsAPI {
	return &permissionOp{client: client}
}

func (op *permissionOp) List(ctx context.Context) ([]v2.PermissionsDataItem, error) {
	res, err := op.client.client.ListPermissions(ctx)
	if err != nil {
		return nil, NewAPIError("Permissions.List", 0, err)
	}

	switch r := res.(type) {
	case *v2.Permissions:
		return r.Data, nil
	case *v2.Error401:
		return nil, NewAPIError("Permissions.List", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.ErrorDefaultStatusCode:
		return nil, NewAPIError("Permissions.List", r.StatusCode, errors.New(string(r.Response.Error.Value.Message.Value)))
	default:
		return nil, unknownError("Permissions.List", r)
	}
}

func (op *permissionOp) Create(ctx context.Context, displayName string, controls v2.BucketControlsBody) (*v2.PermissionData, error) {
	res, err := op.client.client.CreatePermission(ctx, &v2.PermissionBucketControlsBody{
		DisplayName:    v2.DisplayName(displayName),
		BucketControls: controls,
	})
	if err != nil {
		return nil, NewAPIError("Permissions.Create", 0, err)
	}

	switch r := res.(type) {
	case *v2.Permission:
		return &r.Data, nil
	case *v2.Error401:
		return nil, NewAPIError("Permissions.Create", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.Error404:
		return nil, NewAPIError("Permissions.Create", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.Error409:
		return nil, NewAPIError("Permissions.Create", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.ErrorDefaultStatusCode:
		return nil, NewAPIError("Permissions.Create", r.StatusCode, errors.New(string(r.Response.Error.Value.Message.Value)))
	default:
		return nil, unknownError("Permissions.Create", r)
	}
}

func (op *permissionOp) Read(ctx context.Context, permissionId string) (*v2.PermissionData, error) {
	id, err := strconv.Atoi(permissionId)
	if err != nil {
		return nil, NewAPIError("Permissions.Read", 0, err)
	}
	res, err := op.client.client.ReadPermission(ctx, v2.ReadPermissionParams{ID: id})
	if err != nil {
		return nil, NewAPIError("Permissions.Read", 0, err)
	}

	switch r := res.(type) {
	case *v2.Permission:
		return &r.Data, nil
	case *v2.Error401:
		return nil, NewAPIError("Permissions.Read", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.Error404:
		return nil, NewAPIError("Permissions.Read", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.ErrorDefaultStatusCode:
		return nil, NewAPIError("Permissions.Read", r.StatusCode, errors.New(string(r.Response.Error.Value.Message.Value)))
	default:
		return nil, unknownError("Permissions.Read", r)
	}
}

func (op *permissionOp) Update(ctx context.Context, permissionId, name string, controls v2.BucketControlsBody) (*v2.PermissionData, error) {
	id, err := strconv.Atoi(permissionId)
	if err != nil {
		return nil, NewAPIError("Permissions.Update", 0, err)
	}
	res, err := op.client.client.UpdatePermission(ctx, &v2.PermissionBucketControlsBody{
		DisplayName:    v2.DisplayName(name),
		BucketControls: controls,
	}, v2.UpdatePermissionParams{ID: id})
	if err != nil {
		return nil, NewAPIError("Permissions.Update", 0, err)
	}

	switch r := res.(type) {
	case *v2.Permission:
		return &r.Data, nil
	case *v2.Error401:
		return nil, NewAPIError("Permissions.Update", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.Error404:
		return nil, NewAPIError("Permissions.Update", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.ErrorDefaultStatusCode:
		return nil, NewAPIError("Permissions.Update", r.StatusCode, errors.New(string(r.Response.Error.Value.Message.Value)))
	default:
		return nil, unknownError("Permissions.Update", r)
	}
}

func (op *permissionOp) Delete(ctx context.Context, permissionId string) error {
	id, err := strconv.Atoi(permissionId)
	if err != nil {
		return NewAPIError("Permissions.Delete", 0, err)
	}
	res, err := op.client.client.DeletePermission(ctx, v2.DeletePermissionParams{ID: id})
	if err != nil {
		return NewAPIError("Permissions.Delete", 0, err)
	}

	switch r := res.(type) {
	case *v2.DeletePermissionNoContent:
		return nil
	case *v2.Error401:
		return NewAPIError("Permissions.Delete", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.ErrorDefaultStatusCode:
		return NewAPIError("Permissions.Delete", r.StatusCode, errors.New(string(r.Response.Error.Value.Message.Value)))
	default:
		return unknownError("Permissions.Delete", r)
	}
}

func (op *permissionOp) ListAccessKeys(ctx context.Context, permissionId string) ([]v2.PermissionKeyReadResponseBody, error) {
	id, err := strconv.Atoi(permissionId)
	if err != nil {
		return nil, NewAPIError("Permissions.ListAccessKeys", 0, err)
	}
	res, err := op.client.client.ListPermissionKeys(ctx, v2.ListPermissionKeysParams{ID: id})
	if err != nil {
		return nil, NewAPIError("Permissions.ListAccessKeys", 0, err)
	}

	switch r := res.(type) {
	case *v2.HandlerListPermissionKeysResponse:
		return r.Data, nil
	case *v2.Error401:
		return nil, NewAPIError("Permissions.ListAccessKeys", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.ErrorDefaultStatusCode:
		return nil, NewAPIError("Permissions.ListAccessKeys", r.StatusCode, errors.New(string(r.Response.Error.Value.Message.Value)))
	default:
		return nil, unknownError("Permissions.ListAccessKeys", r)
	}
}

func (op *permissionOp) CreateAccessKey(ctx context.Context, permissionId string) (*v2.PermissionKeyCreateResponseBody, error) {
	id, err := strconv.Atoi(permissionId)
	if err != nil {
		return nil, NewAPIError("Permissions.CreateAccessKey", 0, err)
	}
	res, err := op.client.client.CreatePermissionKey(ctx, v2.CreatePermissionKeyParams{ID: id})
	if err != nil {
		return nil, NewAPIError("Permissions.CreateAccessKey", 0, err)
	}

	switch r := res.(type) {
	case *v2.HandlerCreatePermissionKeyResponse:
		return &r.Data, nil
	case *v2.Error401:
		return nil, NewAPIError("Permissions.CreateAccessKey", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.ErrorDefaultStatusCode:
		return nil, NewAPIError("Permissions.CreateAccessKey", r.StatusCode, errors.New(string(r.Response.Error.Value.Message.Value)))
	default:
		return nil, unknownError("Permissions.CreateAccessKey", r)
	}
}

func (op *permissionOp) ReadAccessKey(ctx context.Context, permissionId, accessKeyId string) (*v2.PermissionKeyReadResponseBody, error) {
	id, err := strconv.Atoi(permissionId)
	if err != nil {
		return nil, NewAPIError("Permissions.ReadAccessKey", 0, err)
	}
	res, err := op.client.client.ReadPermissionKey(ctx, v2.ReadPermissionKeyParams{ID: id, KeyID: accessKeyId})
	if err != nil {
		return nil, NewAPIError("Permissions.ReadAccessKey", 0, err)
	}

	switch r := res.(type) {
	case *v2.HandlerReadPermissionKeyResponse:
		return &r.Data, nil
	case *v2.Error401:
		return nil, NewAPIError("Permissions.ReadAccessKey", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.ErrorDefaultStatusCode:
		return nil, NewAPIError("Permissions.ReadAccessKey", r.StatusCode, errors.New(string(r.Response.Error.Value.Message.Value)))
	default:
		return nil, unknownError("Permissions.ReadAccessKey", r)
	}
}

func (op *permissionOp) DeleteAccessKey(ctx context.Context, permissionId, accessKeyId string) error {
	id, err := strconv.Atoi(permissionId)
	if err != nil {
		return NewAPIError("Permissions.DeleteAccessKey", 0, err)
	}
	res, err := op.client.client.DeletePermissionKey(ctx, v2.DeletePermissionKeyParams{ID: id, KeyID: accessKeyId})
	if err != nil {
		return NewAPIError("Permissions.DeleteAccessKey", 0, err)
	}

	switch r := res.(type) {
	case *v2.DeletePermissionKeyNoContent:
		return nil
	case *v2.Error401:
		return NewAPIError("Permissions.DeleteAccessKey", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.ErrorDefaultStatusCode:
		return NewAPIError("Permissions.DeleteAccessKey", r.StatusCode, errors.New(string(r.Response.Error.Value.Message.Value)))
	default:
		return unknownError("Permissions.DeleteAccessKey", r)
	}
}

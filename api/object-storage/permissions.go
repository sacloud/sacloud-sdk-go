// Copyright 2022-2026 The object-storage-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package objectstorage

import (
	"context"
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
	res, err := errorFromDecodedResponse[v2.Permissions]("Permissions.List", func() (any, error) { return op.client.client.ListPermissions(ctx) })
	if err != nil {
		return nil, err
	}

	return res.GetData(), nil
}

func (op *permissionOp) Create(ctx context.Context, displayName string, controls v2.BucketControlsBody) (*v2.PermissionData, error) {
	res, err := errorFromDecodedResponse[v2.Permission]("Permissions.Create", func() (any, error) {
		return op.client.client.CreatePermission(ctx, &v2.PermissionBucketControlsBody{
			DisplayName:    v2.DisplayName(displayName),
			BucketControls: controls,
		})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetData()), nil
}

func (op *permissionOp) Read(ctx context.Context, permissionId string) (*v2.PermissionData, error) {
	id, err := strconv.Atoi(permissionId)
	if err != nil {
		return nil, NewAPIError("Permissions.Read", 0, err)
	}
	res, err := errorFromDecodedResponse[v2.Permission]("Permissions.Read", func() (any, error) { return op.client.client.ReadPermission(ctx, v2.ReadPermissionParams{ID: id}) })
	if err != nil {
		return nil, err
	}

	return new(res.GetData()), nil
}

func (op *permissionOp) Update(ctx context.Context, permissionId, name string, controls v2.BucketControlsBody) (*v2.PermissionData, error) {
	id, err := strconv.Atoi(permissionId)
	if err != nil {
		return nil, NewAPIError("Permissions.Update", 0, err)
	}
	res, err := errorFromDecodedResponse[v2.Permission]("Permissions.Update", func() (any, error) {
		return op.client.client.UpdatePermission(ctx, &v2.PermissionBucketControlsBody{
			DisplayName:    v2.DisplayName(name),
			BucketControls: controls,
		}, v2.UpdatePermissionParams{ID: id})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetData()), nil
}

func (op *permissionOp) Delete(ctx context.Context, permissionId string) error {
	id, err := strconv.Atoi(permissionId)
	if err != nil {
		return NewAPIError("Permissions.Delete", 0, err)
	}
	return errorFromDecodedResponseE[v2.DeletePermissionNoContent]("Permissions.Delete", func() (any, error) { return op.client.client.DeletePermission(ctx, v2.DeletePermissionParams{ID: id}) })
}

func (op *permissionOp) ListAccessKeys(ctx context.Context, permissionId string) ([]v2.PermissionKeyReadResponseBody, error) {
	id, err := strconv.Atoi(permissionId)
	if err != nil {
		return nil, NewAPIError("Permissions.ListAccessKeys", 0, err)
	}
	res, err := errorFromDecodedResponse[v2.HandlerListPermissionKeysResponse]("Permissions.ListAccessKeys", func() (any, error) {
		return op.client.client.ListPermissionKeys(ctx, v2.ListPermissionKeysParams{ID: id})
	})
	if err != nil {
		return nil, err
	}

	return res.GetData(), nil
}

func (op *permissionOp) CreateAccessKey(ctx context.Context, permissionId string) (*v2.PermissionKeyCreateResponseBody, error) {
	id, err := strconv.Atoi(permissionId)
	if err != nil {
		return nil, NewAPIError("Permissions.CreateAccessKey", 0, err)
	}
	res, err := errorFromDecodedResponse[v2.HandlerCreatePermissionKeyResponse]("Permissions.CreateAccessKey", func() (any, error) {
		return op.client.client.CreatePermissionKey(ctx, v2.CreatePermissionKeyParams{ID: id})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetData()), nil
}

func (op *permissionOp) ReadAccessKey(ctx context.Context, permissionId, accessKeyId string) (*v2.PermissionKeyReadResponseBody, error) {
	id, err := strconv.Atoi(permissionId)
	if err != nil {
		return nil, NewAPIError("Permissions.ReadAccessKey", 0, err)
	}
	res, err := errorFromDecodedResponse[v2.HandlerReadPermissionKeyResponse]("Permissions.ReadAccessKey", func() (any, error) {
		return op.client.client.ReadPermissionKey(ctx, v2.ReadPermissionKeyParams{ID: id, KeyID: accessKeyId})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetData()), nil
}

func (op *permissionOp) DeleteAccessKey(ctx context.Context, permissionId, accessKeyId string) error {
	id, err := strconv.Atoi(permissionId)
	if err != nil {
		return NewAPIError("Permissions.DeleteAccessKey", 0, err)
	}
	return errorFromDecodedResponseE[v2.DeletePermissionKeyNoContent]("Permissions.DeleteAccessKey", func() (any, error) {
		return op.client.client.DeletePermissionKey(ctx, v2.DeletePermissionKeyParams{ID: id, KeyID: accessKeyId})
	})
}

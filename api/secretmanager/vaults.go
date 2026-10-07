// Copyright 2025- The sacloud/secretmanager-api-go authors
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

package secretmanager

import (
	"context"

	v1 "github.com/sacloud/sacloud-sdk-go/api/secretmanager/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/packages/into"
)

// VaultAPIはVaultの操作をCRUD+Lで行うためのインターフェース
type VaultAPI interface {
	List(ctx context.Context, count, from *int) ([]v1.Vault, error)
	Read(ctx context.Context, id string) (*v1.Vault, error)
	Create(ctx context.Context, request CreateVaultParams) (*v1.CreateVault, error)
	Update(ctx context.Context, id string, request UpdateVaultParams) (*v1.Vault, error)
	Delete(ctx context.Context, id string) error
}

var _ VaultAPI = (*vaultOp)(nil)

type vaultOp struct {
	client *v1.Client
}

func NewVaultOp(client *v1.Client) VaultAPI {
	return &vaultOp{client: client}
}

func (op *vaultOp) List(ctx context.Context, count, from *int) ([]v1.Vault, error) {
	res, err := op.client.ListVaults(ctx, v1.ListVaultsParams{
		Count: into.Opt[v1.OptInt](count),
		From:  into.Opt[v1.OptInt](from),
	})
	if err != nil {
		return nil, createAPIError("List", err)
	}

	return res.Vaults, nil
}

func (op *vaultOp) Read(ctx context.Context, id string) (*v1.Vault, error) {
	res, err := op.client.ReadVault(ctx, v1.ReadVaultParams{ResourceID: id})
	if err != nil {
		return nil, createAPIError("Read", err)
	}

	return &res.Vault, nil
}

type CreateVaultParams struct {
	Name        string
	Description *string
	KmsKeyID    string
	Tags        []string
}

func (op *vaultOp) Create(ctx context.Context, request CreateVaultParams) (*v1.CreateVault, error) {
	res, err := op.client.CreateVault(ctx, &v1.WrappedCreateVaultRequest{
		Vault: v1.CreateVaultRequest{
			Name:        request.Name,
			Description: into.Opt[v1.OptString](request.Description),
			KmsKeyID:    request.KmsKeyID,
			Tags:        into.OptNilArray[v1.OptNilStringArray](&request.Tags),
		},
	})
	if err != nil {
		return nil, createAPIError("Create", err)
	}

	return &res.Vault, nil
}

type UpdateVaultParams struct {
	Name        string
	Description *string
	Tags        []string
}

func (op *vaultOp) Update(ctx context.Context, id string, request UpdateVaultParams) (*v1.Vault, error) {
	res, err := op.client.UpdateVault(ctx, &v1.WrappedVaultRequest{
		Vault: v1.VaultRequest{
			Name:        request.Name,
			Description: into.Opt[v1.OptString](request.Description),
			Tags:        into.OptNilArray[v1.OptNilStringArray](&request.Tags),
		},
	}, v1.UpdateVaultParams{ResourceID: id})
	if err != nil {
		return nil, createAPIError("Update", err)
	}

	return &res.Vault, nil
}

func (op *vaultOp) Delete(ctx context.Context, id string) error {
	err := op.client.DeleteVault(ctx, v1.DeleteVaultParams{ResourceID: id})
	if err != nil {
		return createAPIError("Delete", err)
	}
	return nil
}

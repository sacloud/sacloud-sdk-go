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

	v1 "github.com/sacloud/secretmanager-api-go/apis/v1"
)

// VaultAPIはVaultの操作をCRUD+Lで行うためのインターフェース
type VaultAPI interface {
	List(ctx context.Context) ([]v1.Vault, error)
	Read(ctx context.Context, id string) (*v1.Vault, error)
	Create(ctx context.Context, request v1.CreateVault) (*v1.CreateVault, error)
	Update(ctx context.Context, id string, request v1.Vault) (*v1.Vault, error)
	Delete(ctx context.Context, id string) error
}

var _ VaultAPI = (*vaultOp)(nil)

type vaultOp struct {
	client *v1.Client
}

func NewVaultOp(client *v1.Client) VaultAPI {
	return &vaultOp{client: client}
}

func (op *vaultOp) List(ctx context.Context) ([]v1.Vault, error) {
	res, err := op.client.SecretmanagerVaultsList(ctx)
	if err != nil {
		return nil, createAPIError("List", err)
	}

	return res.Vaults, nil
}

func (op *vaultOp) Read(ctx context.Context, id string) (*v1.Vault, error) {
	res, err := op.client.SecretmanagerVaultsRetrieve(ctx, v1.SecretmanagerVaultsRetrieveParams{ResourceID: id})
	if err != nil {
		return nil, createAPIError("Read", err)
	}

	return &res.Vault, nil
}

func (op *vaultOp) Create(ctx context.Context, request v1.CreateVault) (*v1.CreateVault, error) {
	res, err := op.client.SecretmanagerVaultsCreate(ctx, &v1.WrappedCreateVault{
		Vault: request,
	})
	if err != nil {
		return nil, createAPIError("Create", err)
	}

	return &res.Vault, nil
}

func (op *vaultOp) Update(ctx context.Context, id string, request v1.Vault) (*v1.Vault, error) {
	res, err := op.client.SecretmanagerVaultsUpdate(ctx, &v1.WrappedVault{
		Vault: request,
	}, v1.SecretmanagerVaultsUpdateParams{ResourceID: id})
	if err != nil {
		return nil, createAPIError("Update", err)
	}

	return &res.Vault, nil
}

func (op *vaultOp) Delete(ctx context.Context, id string) error {
	err := op.client.SecretmanagerVaultsDestroy(ctx, v1.SecretmanagerVaultsDestroyParams{ResourceID: id})
	if err != nil {
		return createAPIError("Delete", err)
	}
	return nil
}

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

// SecretAPIはSecretの操作をCRUD+Lで行うためのインターフェース. READは未実装
type SecretAPI interface {
	List(ctx context.Context) ([]v1.Secret, error)
	// Read(ctx context.Context, id string) (*v1.Secret, error)
	Create(ctx context.Context, request v1.CreateSecret) (*v1.Secret, error)
	Update(ctx context.Context, request v1.CreateSecret) (*v1.Secret, error)
	Delete(ctx context.Context, request v1.DeleteSecret) error
	Unveil(ctx context.Context, request v1.Unveil) (*v1.Unveil, error)
}

var _ SecretAPI = (*secretOp)(nil)

type secretOp struct {
	client  *v1.Client
	vaultId string
}

func NewSecretOp(client *v1.Client, id string) SecretAPI {
	return &secretOp{client: client, vaultId: id}
}

func (op *secretOp) List(ctx context.Context) ([]v1.Secret, error) {
	res, err := op.client.SecretmanagerVaultsSecretsList(ctx,
		v1.SecretmanagerVaultsSecretsListParams{VaultResourceID: op.vaultId})
	if err != nil {
		return nil, createAPIError("List", err)
	}

	return res.Secrets, nil
}

func (op *secretOp) Create(ctx context.Context, request v1.CreateSecret) (*v1.Secret, error) {
	res, err := op.client.SecretmanagerVaultsSecretsCreate(ctx, &v1.WrappedCreateSecret{
		Secret: request,
	}, v1.SecretmanagerVaultsSecretsCreateParams{VaultResourceID: op.vaultId})
	if err != nil {
		return nil, createAPIError("Create", err)
	}

	return &res.Secret, nil
}

// Create / Updateは同じAPIを使うためUpdateは内部でCreateを呼び出すだけ
func (op *secretOp) Update(ctx context.Context, request v1.CreateSecret) (*v1.Secret, error) {
	return op.Create(ctx, request)
}

func (op *secretOp) Unveil(ctx context.Context, request v1.Unveil) (*v1.Unveil, error) {
	res, err := op.client.SecretmanagerVaultsSecretsUnveil(ctx, &v1.WrappedUnveil{
		Secret: request,
	}, v1.SecretmanagerVaultsSecretsUnveilParams{VaultResourceID: op.vaultId})
	if err != nil {
		return nil, createAPIError("Unveil", err)
	}

	return &res.Secret, nil
}

func (op *secretOp) Delete(ctx context.Context, request v1.DeleteSecret) error {
	err := op.client.SecretmanagerVaultsSecretsDestroy(ctx, &v1.WrappedDeleteSecret{
		Secret: request,
	}, v1.SecretmanagerVaultsSecretsDestroyParams{VaultResourceID: op.vaultId})
	if err != nil {
		return createAPIError("Delete", err)
	}
	return nil
}

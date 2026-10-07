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

// SecretAPIはSecretの操作をCRUD+Lで行うためのインターフェース. READは未実装
type SecretAPI interface {
	List(ctx context.Context, count, from *int) ([]v1.SecretResponse, error)
	// Read(ctx context.Context, id string) (*v1.CreateSecretResponse, error)
	Create(ctx context.Context, request CreateSecretParams) (*v1.CreateSecretResponse, error)
	Update(ctx context.Context, request UpdateSecretParams) (*v1.CreateSecretResponse, error)
	Delete(ctx context.Context, request string) error
	Unveil(ctx context.Context, request UnveilParams) (*v1.UnveilResponse, error)
}

var _ SecretAPI = (*secretOp)(nil)

type secretOp struct {
	client  *v1.Client
	vaultId string
}

func NewSecretOp(client *v1.Client, id string) SecretAPI {
	return &secretOp{client: client, vaultId: id}
}

func (op *secretOp) List(ctx context.Context, count, from *int) ([]v1.SecretResponse, error) {
	res, err := op.client.ListVaultSecrets(ctx, v1.ListVaultSecretsParams{
		VaultResourceID: op.vaultId,
		Count:           into.Opt[v1.OptInt](count),
		From:            into.Opt[v1.OptInt](from),
	})
	if err != nil {
		return nil, createAPIError("List", err)
	}

	return res.Secrets, nil
}

type CreateSecretParams = v1.CreateSecretRequest

func (op *secretOp) Create(ctx context.Context, request CreateSecretParams) (*v1.CreateSecretResponse, error) {
	res, err := op.client.CreateVaultSecret(ctx, &v1.WrappedCreateSecretRequest{
		Secret: request,
	}, v1.CreateVaultSecretParams{VaultResourceID: op.vaultId})
	if err != nil {
		return nil, createAPIError("Create", err)
	}

	return &res.Secret, nil
}

type UpdateSecretParams = v1.CreateSecretRequest

// Create / Updateは同じAPIを使うためUpdateは内部でCreateを呼び出すだけ
func (op *secretOp) Update(ctx context.Context, request UpdateSecretParams) (*v1.CreateSecretResponse, error) {
	return op.Create(ctx, request)
}

type UnveilParams struct {
	Name    string
	Version *int
}

func (op *secretOp) Unveil(ctx context.Context, request UnveilParams) (*v1.UnveilResponse, error) {
	res, err := op.client.UnveilSecret(ctx, &v1.WrappedUnveilRequest{
		Secret: v1.UnveilRequest{
			Name:    request.Name,
			Version: into.OptNil[v1.OptNilInt](request.Version),
		},
	}, v1.UnveilSecretParams{VaultResourceID: op.vaultId})
	if err != nil {
		return nil, createAPIError("Unveil", err)
	}

	return &res.Secret, nil
}

func (op *secretOp) Delete(ctx context.Context, request string) error {
	err := op.client.DeleteVaultSecret(ctx, &v1.WrappedDeleteSecretRequest{
		Secret: v1.WrappedDeleteSecretRequestSecret{
			Name: request,
		},
	}, v1.DeleteVaultSecretParams{VaultResourceID: op.vaultId})
	if err != nil {
		return createAPIError("Delete", err)
	}
	return nil
}

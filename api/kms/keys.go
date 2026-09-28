// Copyright 2025- The sacloud/kms-api-go authors
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

package kms

import (
	"context"
	"encoding/base64"
	"errors"

	v1 "github.com/sacloud/sacloud-sdk-go/api/kms/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/packages/into"
)

type KeyAPI interface {
	List(ctx context.Context, count, from *int) ([]v1.Key, error)
	Read(ctx context.Context, id string) (*v1.Key, error)
	Create(ctx context.Context, request CreateParams) (*v1.CreateKeyResponse, error)
	Update(ctx context.Context, id string, request UpdateParams) (*v1.Key, error)
	Delete(ctx context.Context, id string) error

	Rotate(ctx context.Context, id string) (*v1.Key, error)
	ChangeStatus(ctx context.Context, id string, status v1.ChangeKeyStateRequestStatus) (v1.ChangeKeyStateStatus, error)
	ScheduleDestruction(ctx context.Context, id string, pendingDays int) (v1.KeyScheduledDestruction, error)

	Encrypt(ctx context.Context, id string, plain []byte, algo v1.EncryptionRequestAlgo) (string, error)
	Decrypt(ctx context.Context, id, cipher string) ([]byte, error)
}

var _ KeyAPI = (*keyOp)(nil)

type keyOp struct {
	client *v1.Client
}

func NewKeyOp(client *v1.Client) KeyAPI {
	return &keyOp{client: client}
}

func (op *keyOp) List(ctx context.Context, count, from *int) ([]v1.Key, error) {
	res, err := op.client.ListKeys(ctx, v1.ListKeysParams{
		Count: into.Opt[v1.OptInt](count),
		From:  into.Opt[v1.OptInt](from),
	})
	if err != nil {
		return nil, createAPIError("Key.List", err)
	}
	return res.Keys, nil
}

func (op *keyOp) Read(ctx context.Context, id string) (*v1.Key, error) {
	res, err := op.client.ReadKey(ctx, v1.ReadKeyParams{ResourceID: id})
	if err != nil {
		return nil, createAPIError("Key.Read", err)
	}
	return &res.Key, nil
}

type CreateParams struct {
	Name        string
	Description *string
	Tags        []string
	PlainKey    *string
}

func (op *keyOp) Create(ctx context.Context, request CreateParams) (*v1.CreateKeyResponse, error) {
	res, err := op.client.CreateKey(ctx, &v1.WrappedCreateKeyRequest{Key: v1.CreateKeyRequest{
		Name:        request.Name,
		Description: into.Opt[v1.OptString](request.Description),
		Tags:        into.OptNilArray[v1.OptNilStringArray](&request.Tags),
		PlainKey:    into.Opt[v1.OptString](request.PlainKey),
	}})
	if err != nil {
		return nil, createAPIError("Key.Create", err)
	}
	return &res.Key, nil
}

type UpdateParams struct {
	Name        string
	Description *string
	Tags        []string
}

func (op *keyOp) Update(ctx context.Context, id string, request UpdateParams) (*v1.Key, error) {
	res, err := op.client.UpdateKey(ctx, &v1.WrappedKeyRequest{Key: v1.KeyRequest{
		Name:        request.Name,
		Description: into.Opt[v1.OptString](request.Description),
		Tags:        into.OptNilArray[v1.OptNilStringArray](&request.Tags),
	}}, v1.UpdateKeyParams{ResourceID: id})
	if err != nil {
		return nil, createAPIError("Key.Update", err)
	}
	return &res.Key, nil
}

func (op *keyOp) Delete(ctx context.Context, id string) error {
	err := op.client.DeleteKey(ctx, v1.DeleteKeyParams{ResourceID: id})
	if err != nil {
		return createAPIError("Key.Delete", err)
	}
	return nil
}

func (op *keyOp) Rotate(ctx context.Context, id string) (*v1.Key, error) {
	res, err := op.client.RotateKey(ctx, v1.RotateKeyParams{ResourceID: id})
	if err != nil {
		return nil, createAPIError("Key.Rotate", err)
	}

	switch p := any(res).(type) {
	case *v1.WrappedKey:
		return &p.Key, nil
	default:
		return nil, NewAPIError("Key.Rotate", 0, nil)
	}
}

func (op *keyOp) ChangeStatus(ctx context.Context, id string, status v1.ChangeKeyStateRequestStatus) (v1.ChangeKeyStateStatus, error) {
	state, err := op.client.ChangeKeyStatus(ctx, &v1.WrappedChangeKeyStateRequest{Key: v1.ChangeKeyStateRequest{Status: status}}, v1.ChangeKeyStatusParams{ResourceID: id})
	if err != nil {
		return "", createAPIError("Key.ChangeStatus", err)
	}
	return state.GetKey().Status, nil
}

func (op *keyOp) ScheduleDestruction(ctx context.Context, id string, pendingDays int) (v1.KeyScheduledDestruction, error) {
	var zero v1.KeyScheduledDestruction

	if pendingDays < 7 || pendingDays > 90 {
		return zero, NewError("Key.ScheduleDestruction", errors.New("pending days must be between 7 and 90 days"))
	}

	status, err := op.client.ScheduleKeyDestruction(ctx, &v1.WrappedScheduleDestructionKeyRequest{
		Key: v1.ScheduleDestructionKeyRequest{PendingDays: v1.NewOptInt(pendingDays)},
	}, v1.ScheduleKeyDestructionParams{ResourceID: id})
	if err != nil {
		return zero, createAPIError("Key.ScheduleDestruction", err)
	}
	return status.GetKey(), nil
}

func (op *keyOp) Encrypt(ctx context.Context, id string, plain []byte, algo v1.EncryptionRequestAlgo) (string, error) {
	// APIドキュメントではAlgoはRequiredになっていないが、実際にはwriteOnlyの必須フィールドとなっている
	res, err := op.client.EncryptDataWithKey(ctx, &v1.WrappedEncryptionRequest{
		Key: v1.EncryptionRequest{Plain: base64.StdEncoding.EncodeToString(plain), Algo: v1.NewOptEncryptionRequestAlgo(algo)},
	}, v1.EncryptDataWithKeyParams{ResourceID: id})
	if err != nil {
		return "", createAPIError("Key.Encrypt", err)
	}
	return res.Key.Cipher, nil
}

func (op *keyOp) Decrypt(ctx context.Context, id, cipher string) ([]byte, error) {
	res, err := op.client.DecryptDataWithKey(ctx, &v1.WrappedDecryptionRequest{Key: v1.DecryptionRequest{Cipher: cipher}}, v1.DecryptDataWithKeyParams{ResourceID: id})
	if err != nil {
		return nil, createAPIError("Key.Decrypt", err)
	}

	plain, err := base64.StdEncoding.DecodeString(res.Key.Plain)
	if err != nil {
		return nil, NewError("Key.Decrypt", errors.New("got broken base64-encoded plain"))
	}
	return plain, nil
}

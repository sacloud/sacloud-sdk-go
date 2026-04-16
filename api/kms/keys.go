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

	v1 "github.com/sacloud/kms-api-go/apis/v1"
)

type KeyAPI interface {
	List(ctx context.Context) ([]v1.Key, error)
	Read(ctx context.Context, id string) (*v1.Key, error)
	Create(ctx context.Context, request v1.CreateKey) (*v1.CreateKey, error)
	Update(ctx context.Context, id string, request v1.Key) (*v1.Key, error)
	Delete(ctx context.Context, id string) error

	Rotate(ctx context.Context, id string) (*v1.Key, error)
	ChangeStatus(ctx context.Context, id string, status v1.ChangeKeyStatusStatus) error
	ScheduleDestruction(ctx context.Context, id string, pendingDays int) error

	Encrypt(ctx context.Context, id string, plain []byte, algo v1.KeyEncryptAlgoEnum) (string, error)
	Decrypt(ctx context.Context, id, cipher string) ([]byte, error)
}

var _ KeyAPI = (*keyOp)(nil)

type keyOp struct {
	client *v1.Client
}

func NewKeyOp(client *v1.Client) KeyAPI {
	return &keyOp{client: client}
}

func (op *keyOp) List(ctx context.Context) ([]v1.Key, error) {
	res, err := op.client.KmsKeysList(ctx)
	if err != nil {
		return nil, createAPIError("Key.List", err)
	}
	return res.Keys, nil
}

func (op *keyOp) Read(ctx context.Context, id string) (*v1.Key, error) {
	res, err := op.client.KmsKeysRetrieve(ctx, v1.KmsKeysRetrieveParams{ResourceID: id})
	if err != nil {
		return nil, createAPIError("Key.Read", err)
	}
	return &res.Key, nil
}

func (op *keyOp) Create(ctx context.Context, request v1.CreateKey) (*v1.CreateKey, error) {
	res, err := op.client.KmsKeysCreate(ctx, &v1.WrappedCreateKey{
		Key: request,
	})
	if err != nil {
		return nil, createAPIError("Key.Create", err)
	}
	return &res.Key, nil
}

func (op *keyOp) Update(ctx context.Context, id string, request v1.Key) (*v1.Key, error) {
	res, err := op.client.KmsKeysUpdate(ctx, &v1.WrappedKey{
		Key: request,
	}, v1.KmsKeysUpdateParams{ResourceID: id})
	if err != nil {
		return nil, createAPIError("Key.Update", err)
	}
	return &res.Key, nil
}

func (op *keyOp) Delete(ctx context.Context, id string) error {
	err := op.client.KmsKeysDestroy(ctx, v1.KmsKeysDestroyParams{ResourceID: id})
	if err != nil {
		return createAPIError("Key.Delete", err)
	}
	return nil
}

func (op *keyOp) Rotate(ctx context.Context, id string) (*v1.Key, error) {
	res, err := op.client.KmsKeysRotate(ctx, v1.KmsKeysRotateParams{ResourceID: id})
	if err != nil {
		return nil, createAPIError("Key.Rotate", err)
	}

	switch p := res.(type) {
	case *v1.WrappedKey:
		return &p.Key, nil
	case *v1.KmsKeysRotateForbidden:
		return nil, NewAPIError("Key.Rotate", 403, errors.New("forbidden - Key is not available for rotation"))
	default:
		return nil, NewAPIError("Key.Rotate", 0, nil)
	}
}

func (op *keyOp) ChangeStatus(ctx context.Context, id string, status v1.ChangeKeyStatusStatus) error {
	err := op.client.KmsKeysStatus(ctx, &v1.WrappedChangeKeyStatus{
		Key: v1.ChangeKeyStatus{Status: v1.NewOptChangeKeyStatusStatus(status)},
	}, v1.KmsKeysStatusParams{ResourceID: id})
	if err != nil {
		return createAPIError("Key.ChangeStatus", err)
	}
	return nil
}

func (op *keyOp) ScheduleDestruction(ctx context.Context, id string, pendingDays int) error {
	if pendingDays < 7 || pendingDays > 90 {
		return NewError("Key.ScheduleDestruction", errors.New("pending days must be between 7 and 90 days"))
	}

	err := op.client.KmsKeysScheduleDestruction(ctx, &v1.WrappedScheduleDestructionKey{
		Key: v1.ScheduleDestructionKey{PendingDays: pendingDays},
	}, v1.KmsKeysScheduleDestructionParams{ResourceID: id})
	if err != nil {
		return createAPIError("Key.ScheduleDestruction", err)
	}
	return nil
}

func (op *keyOp) Encrypt(ctx context.Context, id string, plain []byte, algo v1.KeyEncryptAlgoEnum) (string, error) {
	// APIドキュメントではAlgoはRequiredになっていないが、実際にはwriteOnlyの必須フィールドとなっている
	res, err := op.client.KmsKeysEncrypt(ctx, &v1.WrappedKeyPlain{
		Key: v1.KeyPlain{Plain: base64.StdEncoding.EncodeToString(plain), Algo: v1.NewOptKeyEncryptAlgoEnum(algo)},
	}, v1.KmsKeysEncryptParams{ResourceID: id})
	if err != nil {
		return "", createAPIError("Key.Encrypt", err)
	}
	return res.Key.Cipher, nil
}

func (op *keyOp) Decrypt(ctx context.Context, id, cipher string) ([]byte, error) {
	res, err := op.client.KmsKeysDecrypt(ctx, &v1.WrappedKeyCipher{Key: v1.KeyCipher{Cipher: cipher}}, v1.KmsKeysDecryptParams{ResourceID: id})
	if err != nil {
		return nil, createAPIError("Key.Decrypt", err)
	}

	plain, err := base64.StdEncoding.DecodeString(res.Key.Plain)
	if err != nil {
		return nil, NewError("Key.Decrypt", errors.New("got broken base64-encoded plain"))
	}
	return plain, nil
}

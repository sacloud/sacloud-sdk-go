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
	res, err := errorFromDecodedResponse("Key.List", func() (*v1.PaginatedKeyList, error) {
		return op.client.ListKeys(ctx, v1.ListKeysParams{
			Count: into.Opt[v1.OptInt](count),
			From:  into.Opt[v1.OptInt](from),
		})
	})
	if err != nil {
		return nil, err
	}

	return res.GetKeys(), nil
}

func (op *keyOp) Read(ctx context.Context, id string) (*v1.Key, error) {
	res, err := errorFromDecodedResponse("Key.Read", func() (*v1.WrappedKey, error) { return op.client.ReadKey(ctx, v1.ReadKeyParams{ResourceID: id}) })
	if err != nil {
		return nil, err
	}

	return new(res.GetKey()), nil
}

type CreateParams struct {
	Name        string
	Description *string
	Tags        []string
	PlainKey    *string
}

func (op *keyOp) Create(ctx context.Context, request CreateParams) (*v1.CreateKeyResponse, error) {
	res, err := errorFromDecodedResponse("Key.Create", func() (*v1.WrappedCreateKeyResponse, error) {
		return op.client.CreateKey(ctx, &v1.WrappedCreateKeyRequest{Key: v1.CreateKeyRequest{
			Name:        request.Name,
			Description: into.Opt[v1.OptString](request.Description),
			Tags:        into.OptNonnilArray[v1.OptNilStringArray](&request.Tags),
			PlainKey:    into.Opt[v1.OptString](request.PlainKey),
		}})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetKey()), nil
}

type UpdateParams struct {
	Name        string
	Description *string
	Tags        []string
}

func (op *keyOp) Update(ctx context.Context, id string, request UpdateParams) (*v1.Key, error) {
	res, err := errorFromDecodedResponse("Key.Update", func() (*v1.WrappedKey, error) {
		return op.client.UpdateKey(ctx, &v1.WrappedKeyRequest{Key: v1.KeyRequest{
			Name:        request.Name,
			Description: into.Opt[v1.OptString](request.Description),
			Tags:        into.OptNonnilArray[v1.OptNilStringArray](&request.Tags),
		}}, v1.UpdateKeyParams{ResourceID: id})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetKey()), nil
}

func (op *keyOp) Rotate(ctx context.Context, id string) (*v1.Key, error) {
	res, err := errorFromDecodedResponse("Key.Rotate", func() (*v1.WrappedKey, error) { return op.client.RotateKey(ctx, v1.RotateKeyParams{ResourceID: id}) })
	if err != nil {
		return nil, err
	}

	return new(res.GetKey()), nil
}

func (op *keyOp) ChangeStatus(ctx context.Context, id string, status v1.ChangeKeyStateRequestStatus) (v1.ChangeKeyStateStatus, error) {
	res, err := errorFromDecodedResponse("Key.ChangeStatus", func() (*v1.WrappedChangeKeyState, error) {
		return op.client.ChangeKeyStatus(ctx,
			&v1.WrappedChangeKeyStateRequest{Key: v1.ChangeKeyStateRequest{Status: status}},
			v1.ChangeKeyStatusParams{ResourceID: id},
		)
	})
	if err != nil {
		var zero v1.ChangeKeyStateStatus
		return zero, err
	}

	return res.GetKey().Status, nil
}

func (op *keyOp) ScheduleDestruction(ctx context.Context, id string, pendingDays int) (v1.KeyScheduledDestruction, error) {
	var zero v1.KeyScheduledDestruction

	if pendingDays < 7 || pendingDays > 90 {
		return zero, NewError("Key.ScheduleDestruction", errors.New("pending days must be between 7 and 90 days"))
	}

	res, err := errorFromDecodedResponse("Key.ScheduleDestruction", func() (*v1.WrappedKeyScheduledDestruction, error) {
		return op.client.ScheduleKeyDestruction(ctx,
			&v1.WrappedScheduleDestructionKeyRequest{Key: v1.ScheduleDestructionKeyRequest{PendingDays: v1.NewOptInt(pendingDays)}},
			v1.ScheduleKeyDestructionParams{ResourceID: id},
		)
	})
	if err != nil {
		return zero, err
	}

	return res.GetKey(), nil
}

func (op *keyOp) Encrypt(ctx context.Context, id string, plain []byte, algo v1.EncryptionRequestAlgo) (string, error) {
	// APIドキュメントではAlgoはRequiredになっていないが、実際にはwriteOnlyの必須フィールドとなっている
	res, err := errorFromDecodedResponse("Key.Encrypt", func() (*v1.WrappedKeyCipher, error) {
		return op.client.EncryptDataWithKey(ctx,
			&v1.WrappedEncryptionRequest{Key: v1.EncryptionRequest{
				Plain: base64.StdEncoding.EncodeToString(plain),
				Algo:  v1.NewOptEncryptionRequestAlgo(algo),
			}},
			v1.EncryptDataWithKeyParams{ResourceID: id},
		)
	})
	if err != nil {
		// この場合、エラーを素直にreturnしてよいかに関してセキュリティ上の懸念がある。
		// 普通、Encryptの引数には機密性の高い情報が含まれていると想定するのが自然であり、
		// レスポンスが正常であれば、暗号化されて受信できていることに異論はなかろうが、
		// あくまでエラーである以上はその前提を置くのは剣呑ではないか?
		// 多分、マスクした方がよいだろう。
		return "", NewError("Key.Encrypt", errors.New("<<details redacted>>"))
	}

	return new(res.GetKey()).GetCipher(), nil
}

func (op *keyOp) Decrypt(ctx context.Context, id, cipher string) ([]byte, error) {
	res, err := errorFromDecodedResponse("Key.Decrypt", func() (*v1.WrappedKeyPlain, error) {
		return op.client.DecryptDataWithKey(ctx, &v1.WrappedDecryptionRequest{Key: v1.DecryptionRequest{Cipher: cipher}}, v1.DecryptDataWithKeyParams{ResourceID: id})
	})
	if err != nil {
		// Encryptと同様の懸念がある。Decryptの応答はもっと直接的に平文である。
		// エラーに含まれる情報が機密性を保持している前提を置くのは危険ではないか。
		// 多分、マスクした方がよいだろう。
		return nil, NewError("Key.Decrypt", errors.New("<<details redacted>>"))
	}

	plain, err := base64.StdEncoding.DecodeString(new(res.GetKey()).GetPlain())
	if err != nil {
		return nil, NewError("Key.Decrypt", errors.New("got broken base64-encoded plain"))
	}

	return plain, nil
}

// Copyright 2016-2025 The terraform-provider-sakura Authors
// SPDX-License-Identifier: Apache-2.0

package nosql

import (
	"context"
	"errors"

	"github.com/google/uuid"
	v1 "github.com/sacloud/sacloud-sdk-go/api/nosql/apis/v1"
)

type BackupAPI interface {
	List(ctx context.Context) ([]v1.NosqlBackup, error)
	Create(ctx context.Context) error
	Restore(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

var _ BackupAPI = (*backupOp)(nil)

type backupOp struct {
	client *v1.Client
	dbId   string
}

func NewBackupOp(client *v1.Client, dbId string) BackupAPI {
	return &backupOp{client: client, dbId: dbId}
}

func (op *backupOp) List(ctx context.Context) ([]v1.NosqlBackup, error) {
	res, err := errorFromDecodedResponse[v1.NosqlBackupResponse]("Backup.List", func() (any, error) {
		return op.client.ListNosqlDatabaseBackups(ctx, v1.ListNosqlDatabaseBackupsParams{ApplianceID: op.dbId})
	})
	if err != nil {
		return nil, err
	}

	n, ok := res.GetNosql().Get()
	if !ok {
		return nil, NewAPIError("Backup.List", 0, errors.New("no Nosql data"))
	}

	return n.GetBackups(), nil
}

func (op *backupOp) Create(ctx context.Context) error {
	return errorFromDecodedResponseE[v1.NosqlOkResponse]("Backup.Create", func() (any, error) {
		return op.client.CreateNosqlDatabaseBackup(ctx, v1.CreateNosqlDatabaseBackupParams{ApplianceID: op.dbId})
	})
}

func (op *backupOp) Restore(ctx context.Context, id uuid.UUID) error {
	return errorFromDecodedResponseE[v1.NosqlOkResponse]("Backup.Restore", func() (any, error) {
		return op.client.RestoreNosqlDatabaseBackup(ctx, v1.RestoreNosqlDatabaseBackupParams{ApplianceID: op.dbId, BackupID: id})
	})
}

func (op *backupOp) Delete(ctx context.Context, id uuid.UUID) error {
	return errorFromDecodedResponseE[v1.NosqlOkResponse]("Backup.Delete", func() (any, error) {
		return op.client.DeleteNosqlDatabaseBackup(ctx, v1.DeleteNosqlDatabaseBackupParams{ApplianceID: op.dbId, BackupID: id})
	})
}

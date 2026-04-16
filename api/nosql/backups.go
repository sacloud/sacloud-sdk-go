// Copyright 2016-2025 The terraform-provider-sakura Authors
// SPDX-License-Identifier: Apache-2.0

package nosql

import (
	"context"
	"errors"

	"github.com/google/uuid"
	v1 "github.com/sacloud/nosql-api-go/apis/v1"
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
	res, err := op.client.GetBackupByApplianceID(ctx, v1.GetBackupByApplianceIDParams{ApplianceID: op.dbId})
	if err != nil {
		return nil, NewAPIError("Backup.List", 0, err)
	}

	switch p := res.(type) {
	case *v1.NosqlBackupResponse:
		return p.Nosql.Value.Backups, nil
	case *v1.BadRequestResponse:
		return nil, NewAPIError("Backup.List", 400, errors.New(p.ErrorMsg.Value))
	case *v1.UnauthorizedResponse:
		return nil, NewAPIError("Backup.List", 401, errors.New(p.ErrorMsg.Value))
	case *v1.ServerErrorResponse:
		return nil, NewAPIError("Backup.List", 500, errors.New(p.ErrorMsg.Value))
	default:
		return nil, NewAPIError("Backup.List", 0, nil)
	}
}

func (op *backupOp) Create(ctx context.Context) error {
	res, err := op.client.CreateBackup(ctx, v1.CreateBackupParams{ApplianceID: op.dbId})
	if err != nil {
		return NewAPIError("Backup.Create", 0, err)
	}

	switch p := res.(type) {
	case *v1.NosqlOkResponse:
		return nil
	case *v1.BadRequestResponse:
		return NewAPIError("Backup.Create", 400, errors.New(p.ErrorMsg.Value))
	case *v1.UnauthorizedResponse:
		return NewAPIError("Backup.Create", 401, errors.New(p.ErrorMsg.Value))
	case *v1.NotFoundResponse:
		return NewAPIError("Backup.Create", 404, errors.New(p.ErrorMsg.Value))
	case *v1.ServerErrorResponse:
		return NewAPIError("Backup.Create", 500, errors.New(p.ErrorMsg.Value))
	default:
		return NewAPIError("Backup.Create", 0, nil)
	}
}

func (op *backupOp) Restore(ctx context.Context, id uuid.UUID) error {
	res, err := op.client.RestoreBackup(ctx, v1.RestoreBackupParams{ApplianceID: op.dbId, BackupID: id})
	if err != nil {
		return NewAPIError("Backup.Restore", 0, err)
	}

	switch p := res.(type) {
	case *v1.NosqlOkResponse:
		return nil
	case *v1.BadRequestResponse:
		return NewAPIError("Backup.Restore", 400, errors.New(p.ErrorMsg.Value))
	case *v1.UnauthorizedResponse:
		return NewAPIError("Backup.Restore", 401, errors.New(p.ErrorMsg.Value))
	case *v1.NotFoundResponse:
		return NewAPIError("Backup.Restore", 404, errors.New(p.ErrorMsg.Value))
	case *v1.ServerErrorResponse:
		return NewAPIError("Backup.Restore", 500, errors.New(p.ErrorMsg.Value))
	default:
		return NewAPIError("Backup.Restore", 0, nil)
	}
}

func (op *backupOp) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := op.client.DeleteBackup(ctx, v1.DeleteBackupParams{ApplianceID: op.dbId, BackupID: id})
	if err != nil {
		return NewAPIError("Backup.Delete", 0, err)
	}

	switch p := res.(type) {
	case *v1.NosqlOkResponse:
		return nil
	case *v1.BadRequestResponse:
		return NewAPIError("Backup.Delete", 400, errors.New(p.ErrorMsg.Value))
	case *v1.UnauthorizedResponse:
		return NewAPIError("Backup.Delete", 401, errors.New(p.ErrorMsg.Value))
	case *v1.NotFoundResponse:
		return NewAPIError("Backup.Delete", 404, errors.New(p.ErrorMsg.Value))
	case *v1.ServerErrorResponse:
		return NewAPIError("Backup.Delete", 500, errors.New(p.ErrorMsg.Value))
	default:
		return NewAPIError("Backup.Delete", 0, nil)
	}
}

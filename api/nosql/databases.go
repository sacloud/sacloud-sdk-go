// Copyright 2016-2025 The terraform-provider-sakura Authors
// SPDX-License-Identifier: Apache-2.0

package nosql

import (
	"context"
	"errors"

	v1 "github.com/sacloud/nosql-api-go/apis/v1"
)

type DatabaseAPI interface {
	List(ctx context.Context) ([]v1.GetNosqlAppliance, error)
	Create(ctx context.Context, plan Plan, request v1.NosqlCreateRequestAppliance) (*v1.NosqlAppliance, error)
	Read(ctx context.Context, id string) (*v1.GetNosqlAppliance, error)
	Update(ctx context.Context, id string, request v1.NosqlUpdateRequestAppliance) error
	Delete(ctx context.Context, id string) error
	ApplyChanges(ctx context.Context, id string) error
	GetStatus(ctx context.Context, id string) (*v1.NosqlStatusResponseApplianceSettingsResponseNosql, error)
}

var _ DatabaseAPI = (*databaseOp)(nil)

type databaseOp struct {
	client *v1.Client
}

func NewDatabaseOp(client *v1.Client) DatabaseAPI {
	return &databaseOp{client: client}
}

func (op *databaseOp) List(ctx context.Context) ([]v1.GetNosqlAppliance, error) {
	res, err := op.client.ListDB(ctx, v1.ListDBParams{FilterClass: "nosql"})
	if err != nil {
		return nil, NewAPIError("Database.List", 0, err)
	}

	switch p := res.(type) {
	case *v1.NosqlListResponse:
		return p.Appliances, nil
	case *v1.BadRequestResponse:
		return nil, NewAPIError("Database.List", 400, errors.New(p.ErrorMsg.Value))
	case *v1.UnauthorizedResponse:
		return nil, NewAPIError("Database.List", 401, errors.New(p.ErrorMsg.Value))
	case *v1.ServerErrorResponse:
		return nil, NewAPIError("Database.List", 500, errors.New(p.ErrorMsg.Value))
	default:
		return nil, NewAPIError("Database.List", 0, nil)
	}
}

func (op *databaseOp) Create(ctx context.Context, plan Plan, request v1.NosqlCreateRequestAppliance) (*v1.NosqlAppliance, error) {
	request.Class = "nosql"
	request.Plan = v1.Plan{ID: plan.GetPlanID()}
	request.ServiceClass = plan.GetServiceClass()
	// TODO: 不具合が修正されたら削除
	request.Remark.Nosql.DiskSize = v1.NewOptNilInt(plan.GetDiskSizeMB())
	request.Remark.Nosql.Memory = v1.NewOptNilInt(plan.GetMemoryMB())
	request.Remark.Nosql.Nodes = v1.NewOptNilInt(plan.GetNodes())
	request.Remark.Nosql.Virtualcore = v1.NewOptNilInt(plan.GetVirtualCore())

	res, err := op.client.CreateDB(ctx, &v1.NosqlCreateRequest{Appliance: request})
	if err != nil {
		return nil, NewAPIError("Database.Create", 0, err)
	}

	switch p := res.(type) {
	case *v1.NosqlCreateResponse:
		return &p.Appliance, nil
	case *v1.BadRequestResponse:
		return nil, NewAPIError("Database.Create", 400, errors.New(p.ErrorMsg.Value))
	case *v1.UnauthorizedResponse:
		return nil, NewAPIError("Database.Create", 401, errors.New(p.ErrorMsg.Value))
	case *v1.ConflictErrorResponse:
		return nil, NewAPIError("Database.Create", 409, errors.New(p.ErrorMsg.Value))
	case *v1.ServerErrorResponse:
		return nil, NewAPIError("Database.Create", 500, errors.New(p.ErrorMsg.Value))
	default:
		return nil, NewAPIError("Database.Create", 0, nil)
	}
}

func (op *databaseOp) Read(ctx context.Context, id string) (*v1.GetNosqlAppliance, error) {
	res, err := op.client.GetDB(ctx, v1.GetDBParams{ApplianceID: id})
	if err != nil {
		return nil, NewAPIError("Database.Read", 0, err)
	}

	switch p := res.(type) {
	case *v1.NosqlGetResponse:
		return &p.Appliance, nil
	case *v1.BadRequestResponse:
		return nil, NewAPIError("Database.Read", 400, errors.New(p.ErrorMsg.Value))
	case *v1.UnauthorizedResponse:
		return nil, NewAPIError("Database.Read", 401, errors.New(p.ErrorMsg.Value))
	case *v1.NotFoundResponse:
		return nil, NewAPIError("Database.Read", 404, errors.New(p.ErrorMsg.Value))
	case *v1.ServerErrorResponse:
		return nil, NewAPIError("Database.Read", 500, errors.New(p.ErrorMsg.Value))
	default:
		return nil, NewAPIError("Database.Read", 0, nil)
	}
}

func (op *databaseOp) Update(ctx context.Context, id string, request v1.NosqlUpdateRequestAppliance) error {
	res, err := op.client.UpdateDB(ctx, &v1.NosqlUpdateRequest{Appliance: request},
		v1.UpdateDBParams{ApplianceID: id})
	if err != nil {
		return NewAPIError("Database.Update", 0, err)
	}

	switch p := res.(type) {
	case *v1.NosqlSuccessResponse:
		return nil
	case *v1.BadRequestResponse:
		return NewAPIError("Database.Update", 400, errors.New(p.ErrorMsg.Value))
	case *v1.UnauthorizedResponse:
		return NewAPIError("Database.Update", 401, errors.New(p.ErrorMsg.Value))
	case *v1.NotFoundResponse:
		return NewAPIError("Database.Update", 404, errors.New(p.ErrorMsg.Value))
	case *v1.ConflictErrorResponse:
		return NewAPIError("Database.Update", 409, errors.New(p.ErrorMsg.Value))
	case *v1.ServerErrorResponse:
		return NewAPIError("Database.Update", 500, errors.New(p.ErrorMsg.Value))
	default:
		return NewAPIError("Database.Update", 0, nil)
	}
}

func (op *databaseOp) Delete(ctx context.Context, id string) error {
	res, err := op.client.DeleteDB(ctx, v1.DeleteDBParams{ApplianceID: id})
	if err != nil {
		return NewAPIError("Database.Delete", 0, err)
	}

	switch p := res.(type) {
	case *v1.NosqlSuccessResponse:
		return nil
	case *v1.BadRequestResponse:
		return NewAPIError("Database.Delete", 400, errors.New(p.ErrorMsg.Value))
	case *v1.UnauthorizedResponse:
		return NewAPIError("Database.Delete", 401, errors.New(p.ErrorMsg.Value))
	case *v1.NotFoundResponse:
		return NewAPIError("Database.Delete", 404, errors.New(p.ErrorMsg.Value))
	case *v1.ConflictErrorResponse:
		return NewAPIError("Database.Delete", 409, errors.New(p.ErrorMsg.Value))
	case *v1.ServerErrorResponse:
		return NewAPIError("Database.Delete", 500, errors.New(p.ErrorMsg.Value))
	default:
		return NewAPIError("Database.Delete", 0, nil)
	}
}

func (op *databaseOp) ApplyChanges(ctx context.Context, id string) error {
	res, err := op.client.UpdateConfigDB(ctx, v1.UpdateConfigDBParams{ApplianceID: id})
	if err != nil {
		return NewAPIError("Database.ApplyChanges", 0, err)
	}

	switch p := res.(type) {
	case *v1.NosqlIsOkResponse:
		return nil
	case *v1.BadRequestResponse:
		return NewAPIError("Database.ApplyChanges", 400, errors.New(p.ErrorMsg.Value))
	default:
		return NewAPIError("Database.ApplyChanges", 0, nil)
	}
}

func (op *databaseOp) GetStatus(ctx context.Context, id string) (*v1.NosqlStatusResponseApplianceSettingsResponseNosql, error) {
	res, err := op.client.ConfirmStatusDB(ctx, v1.ConfirmStatusDBParams{ApplianceID: id})
	if err != nil {
		return nil, NewAPIError("Database.GetStatus", 0, err)
	}

	switch p := res.(type) {
	case *v1.NosqlStatusResponse:
		return &p.Appliance.Value.SettingsResponse.Value.Nosql.Value, nil
	case *v1.BadRequestResponse:
		return nil, NewAPIError("Database.GetStatus", 400, errors.New(p.ErrorMsg.Value))
	case *v1.UnauthorizedResponse:
		return nil, NewAPIError("Database.GetStatus", 401, errors.New(p.ErrorMsg.Value))
	case *v1.ServerErrorResponse:
		return nil, NewAPIError("Database.GetStatus", 500, errors.New(p.ErrorMsg.Value))
	default:
		return nil, NewAPIError("Database.GetStatus", 0, nil)
	}
}

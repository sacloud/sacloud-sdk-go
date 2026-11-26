// Copyright 2016-2025 The terraform-provider-sakura Authors
// SPDX-License-Identifier: Apache-2.0

package nosql

import (
	"context"
	"errors"

	v1 "github.com/sacloud/sacloud-sdk-go/api/nosql/apis/v1"
)

type DatabaseAPI interface {
	List(ctx context.Context) ([]v1.NosqlApplianceResponse, error)
	Create(ctx context.Context, plan Plan, request v1.NosqlCreateRequestAppliance) (*v1.NosqlAppliance, error)
	Read(ctx context.Context, id string) (*v1.NosqlApplianceResponse, error)
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

func (op *databaseOp) List(ctx context.Context) ([]v1.NosqlApplianceResponse, error) {
	res, err := errorFromDecodedResponse[v1.NosqlListResponse]("Database.List", func() (any, error) {
		return op.client.ListNosqlDatabases(ctx, v1.ListNosqlDatabasesParams{FilterClass: "nosql"})
	})
	if err != nil {
		return nil, err
	}

	return res.GetAppliances(), nil
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

	res, err := errorFromDecodedResponse[v1.NosqlCreateResponse]("Database.Create", func() (any, error) {
		return op.client.CreateNosqlDatabase(ctx, &v1.NosqlCreateRequest{Appliance: request})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetAppliance()), nil
}

func (op *databaseOp) Read(ctx context.Context, id string) (*v1.NosqlApplianceResponse, error) {
	res, err := errorFromDecodedResponse[v1.NosqlReadResponse]("Database.Read", func() (any, error) {
		return op.client.ReadNosqlDatabase(ctx, v1.ReadNosqlDatabaseParams{ApplianceID: id})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetAppliance()), nil
}

func (op *databaseOp) Update(ctx context.Context, id string, request v1.NosqlUpdateRequestAppliance) error {
	return errorFromDecodedResponseE[v1.NosqlSuccessResponse]("Database.Update", func() (any, error) {
		return op.client.UpdateNosqlDatabase(ctx, &v1.NosqlUpdateRequest{Appliance: request}, v1.UpdateNosqlDatabaseParams{ApplianceID: id})
	})
}

func (op *databaseOp) Delete(ctx context.Context, id string) error {
	return errorFromDecodedResponseE[v1.NosqlSuccessResponse]("Database.Delete", func() (any, error) {
		return op.client.DeleteNosqlDatabase(ctx, v1.DeleteNosqlDatabaseParams{ApplianceID: id})
	})
}

func (op *databaseOp) ApplyChanges(ctx context.Context, id string) error {
	return errorFromDecodedResponseE[v1.NosqlIsOkResponse]("Database.ApplyChanges", func() (any, error) {
		return op.client.ApplyNosqlDatabaseConfig(ctx, v1.ApplyNosqlDatabaseConfigParams{ApplianceID: id})
	})
}

func (op *databaseOp) GetStatus(ctx context.Context, id string) (*v1.NosqlStatusResponseApplianceSettingsResponseNosql, error) {
	res, err := errorFromDecodedResponse[v1.NosqlStatusResponse]("Database.GetStatus", func() (any, error) {
		return op.client.ReadNosqlDatabaseStatus(ctx, v1.ReadNosqlDatabaseStatusParams{ApplianceID: id})
	})
	if err != nil {
		return nil, NewAPIError("Database.GetStatus", 0, err)
	}

	a, ok := res.GetAppliance().Get()
	if !ok {
		return nil, NewAPIError("Database.GetStatus", 0, errors.New("no appliance data"))
	}

	s, ok := a.GetSettingsResponse().Get()
	if !ok {
		return nil, NewAPIError("Database.GetStatus", 0, errors.New("no settings response"))
	}

	n, ok := s.GetNosql().Get()
	if !ok {
		return nil, NewAPIError("Database.GetStatus", 0, errors.New("no Nosql data"))
	}

	return &n, nil
}

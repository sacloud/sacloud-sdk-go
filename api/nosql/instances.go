// Copyright 2016-2025 The terraform-provider-sakura Authors
// SPDX-License-Identifier: Apache-2.0

package nosql

import (
	"context"
	"errors"

	v1 "github.com/sacloud/sacloud-sdk-go/api/nosql/apis/v1"
)

type InstanceAPI interface {
	GetVersion(ctx context.Context) (*v1.NosqlVersionsResponseNosql, error)
	UpgradeVersion(ctx context.Context, version string) error
	GetParameters(ctx context.Context) ([]v1.NosqlParameter, error)
	SetParameters(ctx context.Context, params []v1.NosqlParameterSetting) error
	GetNodeHealth(ctx context.Context) (v1.NodeHealthNosqlStatus, error)
	AddNodes(ctx context.Context, plan Plan, request v1.NosqlCreateRequestAppliance) (*v1.NosqlAppliance, error)
	Recover(ctx context.Context) (string, error)
	Repair(ctx context.Context, repairType string) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

var _ InstanceAPI = (*instanceOp)(nil)

type instanceOp struct {
	client *v1.Client
	dbId   string
	zone   string
}

func NewInstanceOp(client *v1.Client, dbId string, zone string) InstanceAPI {
	return &instanceOp{client: client, dbId: dbId, zone: zone}
}

func (op *instanceOp) GetVersion(ctx context.Context) (*v1.NosqlVersionsResponseNosql, error) {
	res, err := op.client.ListNosqlDatabaseVersions(ctx, v1.ListNosqlDatabaseVersionsParams{ApplianceID: op.dbId})
	if err != nil {
		return nil, NewAPIError("Instance.GetVersion", 0, err)
	}

	switch p := res.(type) {
	case *v1.NosqlVersionsResponse:
		return &p.Nosql.Value, nil
	case *v1.BadRequestResponse:
		return nil, NewAPIError("Instance.GetVersion", 400, errors.New(p.ErrorMsg.Value))
	case *v1.UnauthorizedResponse:
		return nil, NewAPIError("Instance.GetVersion", 401, errors.New(p.ErrorMsg.Value))
	case *v1.ServerErrorResponse:
		return nil, NewAPIError("Instance.GetVersion", 500, errors.New(p.ErrorMsg.Value))
	default:
		return nil, NewAPIError("Instance.GetVersion", 0, nil)
	}
}

func (op *instanceOp) UpgradeVersion(ctx context.Context, version string) error {
	res, err := op.client.UpdateNosqlDatabaseVersion(ctx, &v1.NosqlUpdateVersionRequest{
		Nosql: v1.NosqlVersion{Version: v1.NewOptString(version)}},
		v1.UpdateNosqlDatabaseVersionParams{ApplianceID: op.dbId})
	if err != nil {
		return NewAPIError("Instance.UpgradeVersion", 0, err)
	}

	switch p := res.(type) {
	case *v1.NosqlUpdateVersionResponse:
		return nil
	case *v1.BadRequestResponse:
		return NewAPIError("Instance.UpgradeVersion", 400, errors.New(p.ErrorMsg.Value))
	case *v1.UnauthorizedResponse:
		return NewAPIError("Instance.UpgradeVersion", 401, errors.New(p.ErrorMsg.Value))
	case *v1.ConflictErrorResponse:
		return NewAPIError("Instance.UpgradeVersion", 409, errors.New(p.ErrorMsg.Value))
	case *v1.ServerErrorResponse:
		return NewAPIError("Instance.UpgradeVersion", 500, errors.New(p.ErrorMsg.Value))
	default:
		return NewAPIError("Instance.UpgradeVersion", 0, nil)
	}
}

func (op *instanceOp) GetParameters(ctx context.Context) ([]v1.NosqlParameter, error) {
	res, err := op.client.ListNosqlDatabaseParameters(ctx, v1.ListNosqlDatabaseParametersParams{ApplianceID: op.dbId})
	if err != nil {
		return nil, NewAPIError("Instance.GetParameters", 0, err)
	}

	switch p := res.(type) {
	case *v1.NosqlParametersResponse:
		return p.Nosql.Value.Parameters, nil
	case *v1.BadRequestResponse:
		return nil, NewAPIError("Instance.GetParameters", 400, errors.New(p.ErrorMsg.Value))
	case *v1.UnauthorizedResponse:
		return nil, NewAPIError("Instance.GetParameters", 401, errors.New(p.ErrorMsg.Value))
	case *v1.ServerErrorResponse:
		return nil, NewAPIError("Instance.GetParameters", 500, errors.New(p.ErrorMsg.Value))
	default:
		return nil, NewAPIError("Instance.GetParameters", 0, nil)
	}
}

func (op *instanceOp) SetParameters(ctx context.Context, params []v1.NosqlParameterSetting) error {
	res, err := op.client.UpdateNosqlDatabaseParameters(ctx, &v1.UpdateNosqlParametersRequest{
		Nosql: v1.UpdateNosqlParametersRequestNosql{Parameters: params}},
		v1.UpdateNosqlDatabaseParametersParams{ApplianceID: op.dbId})
	if err != nil {
		return NewAPIError("Instance.SetParameters", 0, err)
	}

	switch p := res.(type) {
	case *v1.UpdateNosqlParametersResponse:
		return nil
	case *v1.BadRequestResponse:
		return NewAPIError("Instance.SetParameters", 400, errors.New(p.ErrorMsg.Value))
	case *v1.UnauthorizedResponse:
		return NewAPIError("Instance.SetParameters", 401, errors.New(p.ErrorMsg.Value))
	case *v1.ConflictErrorResponse:
		return NewAPIError("Instance.SetParameters", 409, errors.New(p.ErrorMsg.Value))
	case *v1.ServerErrorResponse:
		return NewAPIError("Instance.SetParameters", 500, errors.New(p.ErrorMsg.Value))
	default:
		return NewAPIError("Instance.SetParameters", 0, nil)
	}
}

func (op *instanceOp) GetNodeHealth(ctx context.Context) (v1.NodeHealthNosqlStatus, error) {
	res, err := op.client.ReadNosqlDatabaseNodeHealth(ctx, v1.ReadNosqlDatabaseNodeHealthParams{ApplianceID: op.dbId})
	if err != nil {
		return v1.NodeHealthNosqlStatus(""), NewAPIError("Instance.GetNodeHealth", 0, err)
	}

	switch p := res.(type) {
	case *v1.NodeHealth:
		return p.Nosql.Value.Status.Value, nil
	case *v1.BadRequestResponse:
		return v1.NodeHealthNosqlStatus(""), NewAPIError("Instance.GetNodeHealth", 400, errors.New(p.ErrorMsg.Value))
	case *v1.UnauthorizedResponse:
		return v1.NodeHealthNosqlStatus(""), NewAPIError("Instance.GetNodeHealth", 401, errors.New(p.ErrorMsg.Value))
	case *v1.NotFoundResponse:
		return v1.NodeHealthNosqlStatus(""), NewAPIError("Instance.GetNodeHealth", 404, errors.New(p.ErrorMsg.Value))
	case *v1.ServerErrorResponse:
		return v1.NodeHealthNosqlStatus(""), NewAPIError("Instance.GetNodeHealth", 500, errors.New(p.ErrorMsg.Value))
	default:
		return v1.NodeHealthNosqlStatus(""), NewAPIError("Instance.GetNodeHealth", 0, nil)
	}
}

func (op *instanceOp) AddNodes(ctx context.Context, plan Plan, request v1.NosqlCreateRequestAppliance) (*v1.NosqlAppliance, error) {
	if op.zone == "" {
		return nil, NewError("Instance.AddNodes", errors.New("zone must be specified via NewInstanceOpWithZone"))
	}

	request.Class = "nosql"
	request.Plan = v1.Plan{ID: plan.GetPlanIDforNodes()}
	request.ServiceClass = plan.GetServiceClassForNodes()
	request.Remark.Nosql.PrimaryNodes = v1.NewOptNosqlRemarkNosqlPrimaryNodes(v1.NosqlRemarkNosqlPrimaryNodes{
		Appliance: v1.NosqlRemarkNosqlPrimaryNodesAppliance{ID: op.dbId, Zone: v1.NosqlRemarkNosqlPrimaryNodesApplianceZone{Name: op.zone}},
	})
	res, err := op.client.CreateNosqlDatabase(ctx, &v1.NosqlCreateRequest{Appliance: request})
	if err != nil {
		return nil, NewAPIError("Instance.AddNodes", 0, err)
	}

	switch p := res.(type) {
	case *v1.NosqlCreateResponse:
		return &p.Appliance, nil
	case *v1.BadRequestResponse:
		return nil, NewAPIError("Instance.AddNodes", 400, errors.New(p.ErrorMsg.Value))
	case *v1.UnauthorizedResponse:
		return nil, NewAPIError("Instance.AddNodes", 401, errors.New(p.ErrorMsg.Value))
	case *v1.ConflictErrorResponse:
		return nil, NewAPIError("Instance.AddNodes", 409, errors.New(p.ErrorMsg.Value))
	case *v1.ServerErrorResponse:
		return nil, NewAPIError("Instance.AddNodes", 500, errors.New(p.ErrorMsg.Value))
	default:
		return nil, NewAPIError("Instance.AddNodes", 0, nil)
	}
}

func (op *instanceOp) Recover(ctx context.Context) (string, error) {
	res, err := op.client.RecoverNosqlDatabaseNode(ctx, v1.RecoverNosqlDatabaseNodeParams{ApplianceID: op.dbId})
	if err != nil {
		return "", NewAPIError("Instance.Recover", 0, err)
	}

	switch p := res.(type) {
	case *v1.RecoverNosqlDatabaseNodeOK:
		return "ok", nil
	case *v1.RecoverNosqlDatabaseNodeAccepted:
		return "in_progress", nil
	case *v1.BadRequestResponse:
		return "", NewAPIError("Instance.Recover", 400, errors.New(p.ErrorMsg.Value))
	case *v1.UnauthorizedResponse:
		return "", NewAPIError("Instance.Recover", 401, errors.New(p.ErrorMsg.Value))
	case *v1.NotFoundResponse:
		return "", NewAPIError("Instance.Recover", 404, errors.New(p.ErrorMsg.Value))
	case *v1.ServerErrorResponse:
		return "", NewAPIError("Instance.Recover", 500, errors.New(p.ErrorMsg.Value))
	default:
		return "", NewAPIError("Instance.Recover", 0, nil)
	}
}

func (op *instanceOp) Repair(ctx context.Context, repairType string) error {
	res, err := op.client.RepairNosqlDatabase(ctx, &v1.NosqlRepairRequest{
		Nosql: v1.NewOptNosqlRepairRequestNosql(v1.NosqlRepairRequestNosql{RepairType: v1.NewOptNosqlRepairRequestNosqlRepairType(v1.NosqlRepairRequestNosqlRepairType(repairType))})},
		v1.RepairNosqlDatabaseParams{ApplianceID: op.dbId})
	if err != nil {
		return NewAPIError("Instance.Repair", 0, err)
	}

	switch p := res.(type) {
	case *v1.NosqlRepairResponse:
		return nil
	case *v1.BadRequestResponse:
		return NewAPIError("Instance.Repair", 400, errors.New(p.ErrorMsg.Value))
	case *v1.UnauthorizedResponse:
		return NewAPIError("Instance.Repair", 401, errors.New(p.ErrorMsg.Value))
	case *v1.NotFoundResponse:
		return NewAPIError("Instance.Repair", 404, errors.New(p.ErrorMsg.Value))
	case *v1.ConflictErrorResponse:
		return NewAPIError("Instance.Repair", 409, errors.New(p.ErrorMsg.Value))
	case *v1.ServerErrorResponse:
		return NewAPIError("Instance.Repair", 500, errors.New(p.ErrorMsg.Value))
	default:
		return NewAPIError("Instance.Repair", 0, nil)
	}
}

func (op *instanceOp) Start(ctx context.Context) error {
	res, err := op.client.PowerOnNosqlDatabase(ctx, v1.PowerOnNosqlDatabaseParams{ApplianceID: op.dbId})
	if err != nil {
		return NewAPIError("Instance.Start", 0, err)
	}

	switch p := res.(type) {
	case *v1.SuccessResponse:
		return nil
	case *v1.BadRequestResponse:
		return NewAPIError("Instance.Start", 400, errors.New(p.ErrorMsg.Value))
	case *v1.UnauthorizedResponse:
		return NewAPIError("Instance.Start", 401, errors.New(p.ErrorMsg.Value))
	case *v1.NotFoundResponse:
		return NewAPIError("Instance.Start", 404, errors.New(p.ErrorMsg.Value))
	case *v1.ConflictErrorResponse:
		return NewAPIError("Instance.Start", 409, errors.New(p.ErrorMsg.Value))
	case *v1.ServerErrorResponse:
		return NewAPIError("Instance.Start", 500, errors.New(p.ErrorMsg.Value))
	default:
		return NewAPIError("Instance.Start", 0, nil)
	}
}

func (op *instanceOp) Stop(ctx context.Context) error {
	res, err := op.client.PowerOffNosqlDatabase(ctx, v1.PowerOffNosqlDatabaseParams{ApplianceID: op.dbId})
	if err != nil {
		return NewAPIError("Instance.Stop", 0, err)
	}

	switch p := res.(type) {
	case *v1.SuccessResponse:
		return nil
	case *v1.BadRequestResponse:
		return NewAPIError("Instance.Stop", 400, errors.New(p.ErrorMsg.Value))
	case *v1.UnauthorizedResponse:
		return NewAPIError("Instance.Stop", 401, errors.New(p.ErrorMsg.Value))
	case *v1.NotFoundResponse:
		return NewAPIError("Instance.Stop", 404, errors.New(p.ErrorMsg.Value))
	case *v1.ConflictErrorResponse:
		return NewAPIError("Instance.Stop", 409, errors.New(p.ErrorMsg.Value))
	case *v1.ServerErrorResponse:
		return NewAPIError("Instance.Stop", 500, errors.New(p.ErrorMsg.Value))
	default:
		return NewAPIError("Instance.Stop", 0, nil)
	}
}

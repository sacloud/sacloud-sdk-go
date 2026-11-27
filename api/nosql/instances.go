// Copyright 2016-2025 The terraform-provider-sakura Authors
// SPDX-License-Identifier: Apache-2.0

package nosql

import (
	"context"
	"errors"
	"fmt"

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
	res, err := errorFromDecodedResponse[v1.NosqlVersionsResponse]("Instance.GetVersion", func() (any, error) {
		return op.client.ListNosqlDatabaseVersions(ctx, v1.ListNosqlDatabaseVersionsParams{ApplianceID: op.dbId})
	})
	if err != nil {
		return nil, err
	}

	n, ok := res.GetNosql().Get()
	if !ok {
		return nil, NewAPIError("Instance.GetVersion", 0, errors.New("no Nosql data"))
	}

	return &n, nil
}

func (op *instanceOp) UpgradeVersion(ctx context.Context, version string) error {
	return errorFromDecodedResponseE[v1.NosqlUpdateVersionResponse]("Instance.UpgradeVersion", func() (any, error) {
		return op.client.UpdateNosqlDatabaseVersion(ctx, &v1.NosqlUpdateVersionRequest{Nosql: v1.NosqlVersion{Version: v1.NewOptString(version)}}, v1.UpdateNosqlDatabaseVersionParams{ApplianceID: op.dbId})
	})
}

func (op *instanceOp) GetParameters(ctx context.Context) ([]v1.NosqlParameter, error) {
	res, err := errorFromDecodedResponse[v1.NosqlParametersResponse]("Instance.GetParameters", func() (any, error) {
		return op.client.ListNosqlDatabaseParameters(ctx, v1.ListNosqlDatabaseParametersParams{ApplianceID: op.dbId})
	})
	if err != nil {
		return nil, err
	}

	n, ok := res.GetNosql().Get()
	if !ok {
		return nil, NewAPIError("Instance.GetParameters", 0, errors.New("no Nosql data"))
	}

	return n.GetParameters(), nil
}

func (op *instanceOp) SetParameters(ctx context.Context, params []v1.NosqlParameterSetting) error {
	return errorFromDecodedResponseE[v1.UpdateNosqlParametersResponse]("Instance.SetParameters", func() (any, error) {
		return op.client.UpdateNosqlDatabaseParameters(ctx, &v1.UpdateNosqlParametersRequest{Nosql: v1.UpdateNosqlParametersRequestNosql{Parameters: params}}, v1.UpdateNosqlDatabaseParametersParams{ApplianceID: op.dbId})
	})
}

func (op *instanceOp) GetNodeHealth(ctx context.Context) (v1.NodeHealthNosqlStatus, error) {
	var zero v1.NodeHealthNosqlStatus
	res, err := errorFromDecodedResponse[v1.NodeHealth]("Instance.GetNodeHealth", func() (any, error) {
		return op.client.ReadNosqlDatabaseNodeHealth(ctx, v1.ReadNosqlDatabaseNodeHealthParams{ApplianceID: op.dbId})
	})
	if err != nil {
		return zero, err
	}

	n, ok := res.GetNosql().Get()
	if !ok {
		return zero, NewAPIError("Instance.GetNodeHealth", 0, errors.New("no Nosql data"))
	}

	s, ok := n.GetStatus().Get()
	if !ok {
		return zero, NewAPIError("Instance.GetNodeHealth", 0, errors.New("no status data"))
	}

	return s, nil
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

	res, err := errorFromDecodedResponse[v1.NosqlCreateResponse]("Instance.AddNodes", func() (any, error) {
		return op.client.CreateNosqlDatabase(ctx, &v1.NosqlCreateRequest{Appliance: request})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetAppliance()), nil
}

func (op *instanceOp) Recover(ctx context.Context) (string, error) {
	res, err := errorFromDecodedResponse2[v1.RecoverNosqlDatabaseNodeOK, v1.RecoverNosqlDatabaseNodeAccepted]("Instance.Recover", func() (any, error) {
		return op.client.RecoverNosqlDatabaseNode(ctx, v1.RecoverNosqlDatabaseNodeParams{ApplianceID: op.dbId})
	})
	if err != nil {
		return "", err
	}

	switch p := res.(type) {
	case *v1.RecoverNosqlDatabaseNodeOK:
		return "ok", nil
	case *v1.RecoverNosqlDatabaseNodeAccepted:
		return "in_progress", nil
	default:
		return "", NewAPIError("Instance.Recover", 0, fmt.Errorf("unknown type %#v", p))
	}
}

func (op *instanceOp) Repair(ctx context.Context, repairType string) error {
	return errorFromDecodedResponseE[v1.NosqlRepairResponse]("Instance.Repair", func() (any, error) {
		return op.client.RepairNosqlDatabase(ctx, &v1.NosqlRepairRequest{Nosql: v1.NewOptNosqlRepairRequestNosql(v1.NosqlRepairRequestNosql{RepairType: v1.NewOptNosqlRepairRequestNosqlRepairType(v1.NosqlRepairRequestNosqlRepairType(repairType))})}, v1.RepairNosqlDatabaseParams{ApplianceID: op.dbId})
	})
}

func (op *instanceOp) Start(ctx context.Context) error {
	return errorFromDecodedResponseE[v1.SuccessResponse]("Instance.Start", func() (any, error) {
		return op.client.PowerOnNosqlDatabase(ctx, v1.PowerOnNosqlDatabaseParams{ApplianceID: op.dbId})
	})
}

func (op *instanceOp) Stop(ctx context.Context) error {
	return errorFromDecodedResponseE[v1.SuccessResponse]("Instance.Stop", func() (any, error) {
		return op.client.PowerOffNosqlDatabase(ctx, v1.PowerOffNosqlDatabaseParams{ApplianceID: op.dbId})
	})
}

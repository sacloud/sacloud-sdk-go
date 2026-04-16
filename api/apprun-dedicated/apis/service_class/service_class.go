// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package service_class

import (
	"context"

	v1 "github.com/sacloud/apprun-dedicated-api-go/apis/v1"
	"github.com/sacloud/apprun-dedicated-api-go/common"
)

type ServiceClassAPI interface {
	ListLB(ctx context.Context) (classes []v1.ReadLbServiceClass, err error)
	ListWorker(ctx context.Context) (classes []v1.ReadWorkerServiceClass, err error)
}

type ServiceClassOp struct{ *v1.Client }

func NewServiceClassOp(client *v1.Client) *ServiceClassOp { return &ServiceClassOp{Client: client} }

func (op *ServiceClassOp) ListLB(ctx context.Context) (classes []v1.ReadLbServiceClass, err error) {
	res, err := common.ErrorFromDecodedResponse("ServiceClass.ListLB", func() (*v1.ListLbServiceClassResponse, error) {
		return op.Client.ListLbServiceClasses(ctx)
	})

	if res != nil {
		classes = res.LbServiceClasses
	}

	return
}

func (op *ServiceClassOp) ListWorker(ctx context.Context) (classes []v1.ReadWorkerServiceClass, err error) {
	res, err := common.ErrorFromDecodedResponse("ServiceClass.ListWorker", func() (*v1.ListWorkerServiceClassResponse, error) {
		return op.Client.ListWorkerServiceClasses(ctx)
	})

	if res != nil {
		classes = res.WorkerServiceClasses
	}

	return
}

var _ ServiceClassAPI = (*ServiceClassOp)(nil)

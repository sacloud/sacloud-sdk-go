// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package service_class

import (
	"context"

	v1 "github.com/sacloud/sacloud-sdk-go/api/apprun-dedicated/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/api/apprun-dedicated/common"
)

type ServiceClassAPI interface {
	ListLB(ctx context.Context) (classes []v1.LoadBalancerServiceClass, err error)
	ListWorker(ctx context.Context) (classes []v1.WorkerServiceClass, err error)
}

type ServiceClassOp struct{ *v1.Client }

func NewServiceClassOp(client *v1.Client) *ServiceClassOp { return &ServiceClassOp{Client: client} }

func (op *ServiceClassOp) ListLB(ctx context.Context) (classes []v1.LoadBalancerServiceClass, err error) {
	res, err := common.ErrorFromDecodedResponse("ServiceClass.ListLB", func() (*v1.ListLoadBalancerServiceClassesResponse, error) {
		return op.Client.ListLoadBalancerServiceClasses(ctx)
	})

	if res != nil {
		classes = res.LbServiceClasses
	}

	return
}

func (op *ServiceClassOp) ListWorker(ctx context.Context) (classes []v1.WorkerServiceClass, err error) {
	res, err := common.ErrorFromDecodedResponse("ServiceClass.ListWorker", func() (*v1.ListWorkerServiceClassesResponse, error) {
		return op.Client.ListWorkerServiceClasses(ctx)
	})

	if res != nil {
		classes = res.WorkerServiceClasses
	}

	return
}

var _ ServiceClassAPI = (*ServiceClassOp)(nil)

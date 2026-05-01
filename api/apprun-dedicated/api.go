// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package apprun_dedicated

import (
	"github.com/sacloud/sacloud-sdk-go/api/apprun-dedicated/apis/application"
	"github.com/sacloud/sacloud-sdk-go/api/apprun-dedicated/apis/autoscalinggroup"
	"github.com/sacloud/sacloud-sdk-go/api/apprun-dedicated/apis/certificate"
	"github.com/sacloud/sacloud-sdk-go/api/apprun-dedicated/apis/cluster"
	"github.com/sacloud/sacloud-sdk-go/api/apprun-dedicated/apis/loadbalancer"
	"github.com/sacloud/sacloud-sdk-go/api/apprun-dedicated/apis/service_class"
	"github.com/sacloud/sacloud-sdk-go/api/apprun-dedicated/apis/version"
	"github.com/sacloud/sacloud-sdk-go/api/apprun-dedicated/apis/workernode"
)

type ApplicationAPI = application.ApplicationAPI
type AutoScalingGroupAPI = autoscalinggroup.AutoScalingGroupAPI
type CertificateAPI = certificate.CertificateAPI
type ClusterAPI = cluster.ClusterAPI
type LoadBalancerAPI = loadbalancer.LoadBalancerAPI
type ServiceClassAPI = service_class.ServiceClassAPI
type VersionAPI = version.VersionAPI
type WorkerNodeAPI = workernode.WorkerNodeAPI

var NewApplicationOp = application.NewApplicationOp
var NewAutoScalingGroupOp = autoscalinggroup.NewAutoScalingGroupOp
var NewCertificateOp = certificate.NewCertificateOp
var NewClusterOp = cluster.NewClusterOp
var NewLoadBalancerOp = loadbalancer.NewLoadBalancerOp
var NewServiceClassOp = service_class.NewServiceClassOp
var NewVersionOp = version.NewVersionOp
var NewWorkerNodeOp = workernode.NewWorkerNodeOp

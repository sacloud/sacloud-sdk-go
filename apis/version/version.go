// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package version

import (
	"context"

	v1 "github.com/sacloud/apprun-dedicated-api-go/apis/v1"
	"github.com/sacloud/apprun-dedicated-api-go/common"
	"github.com/sacloud/saclient-go"
)

type VersionAPI interface {
	// List returns the list of ApplicationVersions, paginated.
	// Pass nil to `cursor` to get the first page, or
	// previously returned `nextCursor` to get the next page.
	List(ctx context.Context, elems int64, cursor *v1.ApplicationVersionNumber) (list []v1.ApplicationVersionDeploymentStatus, nextCursor *v1.ApplicationVersionNumber, err error)
	Create(ctx context.Context, params CreateParams) (version *v1.ReadApplicationVersionSummary, err error)
	Read(ctx context.Context, id v1.ApplicationVersionNumber) (version *VersionDetail, err error)
	Delete(ctx context.Context, id v1.ApplicationVersionNumber) error
}

type VersionOp struct {
	client        *v1.Client
	applicationID v1.ApplicationID
}

func NewVersionOp(client *v1.Client, applicationID v1.ApplicationID) *VersionOp {
	return &VersionOp{
		client:        client,
		applicationID: applicationID,
	}
}

func (op *VersionOp) List(ctx context.Context, maxItems int64, cursor *v1.ApplicationVersionNumber) (versions []v1.ApplicationVersionDeploymentStatus, nextCursor *v1.ApplicationVersionNumber, err error) {
	res, err := common.ErrorFromDecodedResponse("Version.List", func() (*v1.ListApplicationVersionResponse, error) {
		return op.client.ListApplicationVersions(ctx, v1.ListApplicationVersionsParams{
			ApplicationID: op.applicationID,
			Cursor:        common.IntoOpt[v1.OptApplicationVersionNumber](cursor),
			MaxItems:      maxItems,
		})
	})

	if res != nil {
		versions = res.Versions
		nextCursor = common.FromOpt(res.NextCursor)
	}

	return
}

func (op *VersionOp) Create(ctx context.Context, params CreateParams) (ver *v1.ReadApplicationVersionSummary, err error) {
	res, err := common.ErrorFromDecodedResponse("Version.Create", func() (*v1.CreateApplicationVersionResponse, error) {
		return op.client.CreateApplicationVersion(ctx, saclient.Ptr(params.into()), v1.CreateApplicationVersionParams{
			ApplicationID: op.applicationID,
		})
	})

	if res != nil {
		ver = &res.ApplicationVersion
	}

	return
}

func (op *VersionOp) Read(ctx context.Context, id v1.ApplicationVersionNumber) (ver *VersionDetail, err error) {
	res, err := common.ErrorFromDecodedResponse("Version.Read", func() (*v1.GetApplicationVersionResponse, error) {
		return op.client.GetApplicationVersion(ctx, v1.GetApplicationVersionParams{
			ApplicationID: op.applicationID,
			Version:       id,
		})
	})

	if res != nil {
		var detail VersionDetail
		detail.from(&res.ApplicationVersion)
		ver = &detail
	}

	return
}

func (op *VersionOp) Delete(ctx context.Context, id v1.ApplicationVersionNumber) error {
	return common.ErrorFromDecodedResponseE("Version.Delete", func() error {
		return op.client.DeleteApplicationVersion(ctx, v1.DeleteApplicationVersionParams{
			ApplicationID: op.applicationID,
			Version:       id,
		})
	})
}

var _ VersionAPI = (*VersionOp)(nil)

type ExposedPort struct {
	TargetPort       v1.Port
	LoadBalancerPort *v1.Port
	UseLetsEncrypt   bool
	Host             []string
	HealthCheck      *v1.HealthCheck
}

func (p ExposedPort) into() (ret v1.ExposedPort) {
	ret.SetTargetPort(p.TargetPort)
	ret.SetLoadBalancerPort(common.IntoNullable[v1.NilPort](p.LoadBalancerPort))
	ret.SetUseLetsEncrypt(p.UseLetsEncrypt)
	ret.SetHost(p.Host)
	ret.SetHealthCheck(common.IntoNullable[v1.NilHealthCheck](p.HealthCheck))

	return
}

func (p *ExposedPort) From(res *v1.ExposedPort) {
	p.TargetPort = res.GetTargetPort()
	p.LoadBalancerPort = common.FromOpt(res.GetLoadBalancerPort())
	p.UseLetsEncrypt = res.GetUseLetsEncrypt()
	p.Host = res.GetHost()
	p.HealthCheck = common.FromOpt(res.GetHealthCheck())
}

type EnvironmentVariable struct {
	Key    string
	Value  *string
	Secret bool
}

func (e EnvironmentVariable) into() (ret v1.CreateEnvironmentVariable) {
	ret.SetKey(e.Key)
	ret.SetValue(common.IntoOpt[v1.OptString](e.Value))
	ret.SetSecret(e.Secret)

	return
}

func (e *EnvironmentVariable) From(res *v1.ReadEnvironmentVariable) {
	e.Key = res.GetKey()
	e.Value = common.FromOpt(res.GetValue())
	e.Secret = res.GetSecret()
}

type CreateParams struct {
	CPU                    int64
	Memory                 int64
	ScalingMode            v1.ScalingMode
	FixedScale             *int32
	MinScale               *int32
	MaxScale               *int32
	ScaleInThreshold       *int32
	ScaleOutThreshold      *int32
	Image                  string
	Cmd                    []string
	RegistryUsername       *string
	RegistryPassword       *string
	RegistryPasswordAction v1.RegistryPasswordAction
	ExposedPorts           []ExposedPort
	EnvVar                 []EnvironmentVariable
}

func (c *CreateParams) into() (ret v1.CreateApplicationVersion) {
	ret.SetCPU(c.CPU)
	ret.SetMemory(c.Memory)
	ret.SetScalingMode(c.ScalingMode)
	ret.SetFixedScale(common.IntoOpt[v1.OptInt32](c.FixedScale))
	ret.SetMinScale(common.IntoOpt[v1.OptInt32](c.MinScale))
	ret.SetMaxScale(common.IntoOpt[v1.OptInt32](c.MaxScale))
	ret.SetScaleInThreshold(common.IntoOpt[v1.OptInt32](c.ScaleInThreshold))
	ret.SetScaleOutThreshold(common.IntoOpt[v1.OptInt32](c.ScaleOutThreshold))
	ret.SetImage(c.Image)
	ret.SetCmd(c.Cmd)
	ret.SetRegistryUsername(common.IntoNullable[v1.NilString](c.RegistryUsername))
	ret.SetRegistryPassword(common.IntoNullable[v1.NilString](c.RegistryPassword))
	ret.SetRegistryPasswordAction(c.RegistryPasswordAction)
	ret.SetExposedPorts(common.MapSlice(c.ExposedPorts, ExposedPort.into))
	ret.SetEnv(common.MapSlice(c.EnvVar, EnvironmentVariable.into))

	return
}

type VersionDetail struct {
	Version           v1.ApplicationVersionNumber
	CPU               int64
	Memory            int64
	ScalingMode       v1.ScalingMode
	FixedScale        *int32
	MinScale          *int32
	MaxScale          *int32
	ScaleInThreshold  *int32
	ScaleOutThreshold *int32
	Image             string
	Cmd               []string
	RegistryUsername  *string
	RegistryPassword  *string
	ActiveNodeCount   int64
	Created           int
	ExposedPorts      []ExposedPort
	EnvVar            []EnvironmentVariable
}

func (v *VersionDetail) from(res *v1.ReadApplicationVersionDetail) {
	v.Version = res.GetVersion()
	v.CPU = res.GetCPU()
	v.Memory = res.GetMemory()
	v.ScalingMode = res.GetScalingMode()
	v.FixedScale = common.FromOpt(res.GetFixedScale())
	v.MinScale = common.FromOpt(res.GetMinScale())
	v.MaxScale = common.FromOpt(res.GetMaxScale())
	v.ScaleInThreshold = common.FromOpt(res.GetScaleInThreshold())
	v.ScaleOutThreshold = common.FromOpt(res.GetScaleOutThreshold())
	v.Image = res.GetImage()
	v.Cmd = res.GetCmd()
	v.RegistryUsername = common.FromOpt(res.GetRegistryUsername())
	v.RegistryPassword = common.FromOpt(res.GetRegistryPassword())
	v.ActiveNodeCount = res.GetActiveNodeCount()
	v.Created = res.GetCreated()
	v.ExposedPorts = common.MapSlice(res.GetExposedPorts(), common.ConvertFrom[v1.ExposedPort, ExposedPort]())
	v.EnvVar = common.MapSlice(res.GetEnv(), common.ConvertFrom[v1.ReadEnvironmentVariable, EnvironmentVariable]())
}

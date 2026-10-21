// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package version

import (
	"context"

	v1 "github.com/sacloud/sacloud-sdk-go/api/apprun-dedicated/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/api/apprun-dedicated/common"
	"github.com/sacloud/sacloud-sdk-go/common/packages/into"
)

type VersionAPI interface {
	// List returns the list of ApplicationVersions, paginated.
	// Pass nil to `cursor` to get the first page, or
	// previously returned `nextCursor` to get the next page.
	List(ctx context.Context, maxItems int64, cursor *v1.ApplicationVersionNumber) (list []v1.ApplicationVersionSummary, nextCursor *v1.ApplicationVersionNumber, err error)
	Create(ctx context.Context, params CreateParams) (version *v1.ApplicationVersionIdentifier, err error)
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

func (op *VersionOp) List(ctx context.Context, maxItems int64, cursor *v1.ApplicationVersionNumber) (list []v1.ApplicationVersionSummary, nextCursor *v1.ApplicationVersionNumber, err error) {
	res, err := common.ErrorFromDecodedResponse("Version.List", func() (*v1.ListApplicationVersionsResponse, error) {
		return op.client.ListApplicationVersions(ctx, v1.ListApplicationVersionsParams{
			ApplicationID: op.applicationID,
			Cursor:        into.Opt[v1.OptApplicationVersionNumber](cursor),
			MaxItems:      maxItems,
		})
	})

	if res != nil {
		list = res.Versions
		nextCursor = common.FromOpt(res.NextCursor)
	}

	return
}

func (op *VersionOp) Create(ctx context.Context, params CreateParams) (ver *v1.ApplicationVersionIdentifier, err error) {
	res, err := common.ErrorFromDecodedResponse("Version.Create", func() (*v1.CreateApplicationVersionResponse, error) {
		return op.client.CreateApplicationVersion(ctx, new(params.into()), v1.CreateApplicationVersionParams{
			ApplicationID: op.applicationID,
		})
	})

	if res != nil {
		ver = &res.ApplicationVersion
	}

	return
}

func (op *VersionOp) Read(ctx context.Context, id v1.ApplicationVersionNumber) (ver *VersionDetail, err error) {
	res, err := common.ErrorFromDecodedResponse("Version.Read", func() (*v1.ReadApplicationVersionResponse, error) {
		return op.client.ReadApplicationVersion(ctx, v1.ReadApplicationVersionParams{
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
	TargetPort       v1.Port         `json:"targetPort"`
	LoadBalancerPort *v1.Port        `json:"loadBalancerPort,omitempty"`
	UseLetsEncrypt   bool            `json:"useLetsEncrypt"`
	Host             []string        `json:"host,omitempty"`
	HealthCheck      *v1.HealthCheck `json:"healthCheck,omitempty"`
}

func (p ExposedPort) into() (ret v1.ExposedPort) {
	ret.SetTargetPort(p.TargetPort)
	ret.SetLoadBalancerPort(into.Nil[v1.NilPort](p.LoadBalancerPort))
	ret.SetUseLetsEncrypt(p.UseLetsEncrypt)
	var hosts []v1.Hostname
	switch {
	case p.Host == nil:
		hosts = nil
	case len(p.Host) == 0:
		hosts = make([]v1.Hostname, 0)
	default:
		for _, host := range p.Host {
			hosts = append(hosts, v1.Hostname(host))
		}
	}
	ret.SetHost(hosts)
	ret.SetHealthCheck(into.Nil[v1.NilHealthCheck](p.HealthCheck))

	return
}

func (p *ExposedPort) From(res *v1.ExposedPort) {
	p.TargetPort = res.GetTargetPort()
	p.LoadBalancerPort = common.FromOpt(res.GetLoadBalancerPort())
	p.UseLetsEncrypt = res.GetUseLetsEncrypt()
	var hosts []string
	switch {
	case res.GetHost() == nil:
		hosts = nil
	case len(res.GetHost()) == 0:
		hosts = make([]string, 0)
	default:
		for _, host := range res.GetHost() {
			hosts = append(hosts, string(host))
		}
	}
	p.Host = hosts
	p.HealthCheck = common.FromOpt(res.GetHealthCheck())
}

type EnvironmentVariable struct {
	Key    string  `json:"key"`
	Value  *string `json:"value,omitempty"`
	Secret bool    `json:"secret"`
}

func (e EnvironmentVariable) into() (ret v1.CreateEnvironmentVariable) {
	var v *v1.EnvironmentVariableValue
	if e.Value != nil {
		v = new(v1.EnvironmentVariableValue(*e.Value))
	}
	ret.SetKey(e.Key)
	ret.SetValue(into.Opt[v1.OptEnvironmentVariableValue](v))
	ret.SetSecret(e.Secret)

	return
}

func (e *EnvironmentVariable) From(res *v1.ReadEnvironmentVariable) {
	var v *string
	if w := common.FromOpt(res.GetValue()); w != nil {
		v = new(string(*w))
	}
	e.Key = res.GetKey()
	e.Value = v
	e.Secret = res.GetSecret()
}

type CreateParams struct {
	CPU                    int64                     `json:"cpu"`
	Memory                 int64                     `json:"memory"`
	ScalingMode            v1.ScalingMode            `json:"scalingMode"`
	FixedScale             *int32                    `json:"fixedScale,omitempty"`
	MinScale               *int32                    `json:"minScale,omitempty"`
	MaxScale               *int32                    `json:"maxScale,omitempty"`
	ScaleInThreshold       *int32                    `json:"scaleInThreshold,omitempty"`
	ScaleOutThreshold      *int32                    `json:"scaleOutThreshold,omitempty"`
	Image                  string                    `json:"image"`
	Cmd                    []string                  `json:"cmd,omitempty"`
	RegistryUsername       *string                   `json:"registryUsername,omitempty"`
	RegistryPassword       *string                   `json:"registryPassword,omitempty"`
	RegistryPasswordAction v1.RegistryPasswordAction `json:"registryPasswordAction,omitempty"`
	ExposedPorts           []ExposedPort             `json:"exposedPorts,omitempty"`
	EnvVars                []EnvironmentVariable     `json:"envVars,omitempty"`
}

func (c *CreateParams) into() (ret v1.CreateApplicationVersionRequest) {
	ret.SetCPU(c.CPU)
	ret.SetMemory(c.Memory)
	ret.SetScalingMode(c.ScalingMode)
	ret.SetFixedScale(into.Opt[v1.OptInt32](c.FixedScale))
	ret.SetMinScale(into.Opt[v1.OptInt32](c.MinScale))
	ret.SetMaxScale(into.Opt[v1.OptInt32](c.MaxScale))
	ret.SetScaleInThreshold(into.Opt[v1.OptInt32](c.ScaleInThreshold))
	ret.SetScaleOutThreshold(into.Opt[v1.OptInt32](c.ScaleOutThreshold))
	ret.SetImage(c.Image)
	ret.SetCmd(c.Cmd)
	ret.SetRegistryUsername(into.Nil[v1.NilString](c.RegistryUsername))
	ret.SetRegistryPassword(into.Nil[v1.NilString](c.RegistryPassword))
	ret.SetRegistryPasswordAction(c.RegistryPasswordAction)
	ret.SetExposedPorts(common.MapSlice(c.ExposedPorts, ExposedPort.into))
	ret.SetEnv(common.MapSlice(c.EnvVars, EnvironmentVariable.into))

	return
}

type VersionDetail struct {
	Version           v1.ApplicationVersionNumber `json:"version"`
	CPU               int64                       `json:"cpu"`
	Memory            int64                       `json:"memory"`
	ScalingMode       v1.ScalingMode              `json:"scalingMode"`
	FixedScale        *int32                      `json:"fixedScale,omitempty"`
	MinScale          *int32                      `json:"minScale,omitempty"`
	MaxScale          *int32                      `json:"maxScale,omitempty"`
	ScaleInThreshold  *int32                      `json:"scaleInThreshold,omitempty"`
	ScaleOutThreshold *int32                      `json:"scaleOutThreshold,omitempty"`
	Image             string                      `json:"image"`
	Cmd               []string                    `json:"cmd,omitempty"`
	RegistryUsername  *string                     `json:"registryUsername,omitempty"`
	RegistryPassword  *string                     `json:"registryPassword,omitempty"`
	ActiveNodeCount   int64                       `json:"activeNodeCount"`
	Created           int                         `json:"created"`
	ExposedPorts      []ExposedPort               `json:"exposedPorts,omitempty"`
	EnvVars           []EnvironmentVariable       `json:"envVars,omitempty"`
}

func (v *VersionDetail) from(res *v1.ReadApplicationVersion) {
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
	v.EnvVars = common.MapSlice(res.GetEnv(), common.ConvertFrom[v1.ReadEnvironmentVariable, EnvironmentVariable]())
}

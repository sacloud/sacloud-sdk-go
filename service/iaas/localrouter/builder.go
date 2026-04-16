// Copyright 2022-2025 The sacloud/iaas-service-go Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package localrouter

import (
	"context"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	localrouter "github.com/sacloud/iaas-service-go/localrouter/builder"
)

// Builder ローカルルータの構築を行う
type Builder struct {
	ID          types.ID
	Name        string
	Description string
	Tags        types.Tags
	IconID      types.ID

	Switch       *iaas.LocalRouterSwitch
	Interface    *iaas.LocalRouterInterface
	Peers        []*iaas.LocalRouterPeer
	StaticRoutes []*iaas.LocalRouterStaticRoute

	SettingsHash string

	Caller iaas.APICaller
}

func BuilderFromResource(ctx context.Context, caller iaas.APICaller, id types.ID) (*Builder, error) {
	client := iaas.NewLocalRouterOp(caller)
	current, err := client.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	return &Builder{
		Name:         current.Name,
		Description:  current.Description,
		Tags:         current.Tags,
		IconID:       current.IconID,
		Switch:       current.Switch,
		Interface:    current.Interface,
		Peers:        current.Peers,
		StaticRoutes: current.StaticRoutes,
		Caller:       caller,
	}, nil
}

func (b *Builder) Build(ctx context.Context) (*iaas.LocalRouter, error) {
	builder := &localrouter.Builder{
		Name:         b.Name,
		Description:  b.Description,
		Tags:         b.Tags,
		IconID:       b.IconID,
		Switch:       b.Switch,
		Interface:    b.Interface,
		Peers:        b.Peers,
		StaticRoutes: b.StaticRoutes,
		SettingsHash: b.SettingsHash,
		Client:       localrouter.NewAPIClient(b.Caller),
	}

	if b.ID.IsEmpty() {
		return builder.Build(ctx)
	}
	return builder.Update(ctx, b.ID)
}

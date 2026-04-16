// Copyright 2016-2022 The sacloud/iaas-api-go Authors
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

package query

import (
	"context"
	"fmt"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/plans"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/types"
)

func findByPreviousIDCondition(id types.ID) *iaas.FindCondition {
	return &iaas.FindCondition{
		Filter: search.Filter{
			search.Key("Tags.Name"): search.TagsAndEqual(fmt.Sprintf("%s=%s", plans.PreviousIDTagName, id)),
		},
	}
}

// ReadServer 指定のIDでサーバを検索、IDで見つからなかった場合は@previous-idタグで検索し見つかったサーバリソースを返す
//
// 対象が見つからなかった場合はiaas.NoResultsErrorを返す
func ReadServer(ctx context.Context, caller iaas.APICaller, zone string, id types.ID) (*iaas.Server, error) {
	serverOp := iaas.NewServerOp(caller)

	server, err := serverOp.Read(ctx, zone, id)
	if err != nil {
		if !iaas.IsNotFoundError(err) {
			return nil, err
		}

		found, err := serverOp.Find(ctx, zone, findByPreviousIDCondition(id))
		if err != nil {
			return nil, err
		}
		if len(found.Servers) == 0 {
			return nil, iaas.NewNoResultsError()
		}

		// 複数ヒットした場合でも先頭だけ返す
		server = found.Servers[0]
	}

	return server, nil
}

// ReadRouter 指定のIDでルータを検索、IDで見つからなかった場合は@previous-idタグで検索し見つかったリソースを返す
//
// 対象が見つからなかった場合はiaas.NoResultsErrorを返す
func ReadRouter(ctx context.Context, caller iaas.APICaller, zone string, id types.ID) (*iaas.Internet, error) {
	routerOp := iaas.NewInternetOp(caller)

	router, err := routerOp.Read(ctx, zone, id)
	if err != nil {
		if !iaas.IsNotFoundError(err) {
			return nil, err
		}

		found, err := routerOp.Find(ctx, zone, findByPreviousIDCondition(id))
		if err != nil {
			return nil, err
		}
		if len(found.Internet) == 0 {
			return nil, iaas.NewNoResultsError()
		}

		// 複数ヒットした場合でも先頭だけ返す
		router = found.Internet[0]
	}

	return router, nil
}

// ReadProxyLB 指定のIDでELBを検索、IDで見つからなかった場合は@previous-idタグで検索し見つかったリソースを返す
//
// 対象が見つからなかった場合はiaas.NoResultsErrorを返す
func ReadProxyLB(ctx context.Context, caller iaas.APICaller, id types.ID) (*iaas.ProxyLB, error) {
	elbOp := iaas.NewProxyLBOp(caller)

	elb, err := elbOp.Read(ctx, id)
	if err != nil {
		if !iaas.IsNotFoundError(err) {
			return nil, err
		}

		found, err := elbOp.Find(ctx, findByPreviousIDCondition(id))
		if err != nil {
			return nil, err
		}
		if len(found.ProxyLBs) == 0 {
			return nil, iaas.NewNoResultsError()
		}

		// 複数ヒットした場合でも先頭だけ返す
		elb = found.ProxyLBs[0]
	}

	return elb, nil
}

// Copyright 2022-2025 The sacloud/iaas-api-go Authors
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

package test

import (
	"context"
	"fmt"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/search"
)

func lookupDNSByName(caller iaas.APICaller, zoneName string) (*iaas.DNS, error) {
	dnsOp := iaas.NewDNSOp(caller)
	searched, err := dnsOp.Find(context.Background(), &iaas.FindCondition{
		Count: 1,
		Filter: search.Filter{
			search.Key("Name"): zoneName,
		},
	})
	if err != nil {
		return nil, err
	}
	if searched.Count == 0 {
		return nil, fmt.Errorf("dns zone %q is not found", zoneName)
	}

	// 部分一致などにより予期せぬゾーンとマッチしていないかチェック
	if searched.DNS[0].Name != zoneName {
		return nil, fmt.Errorf("fetched dns zone does not match to desired: param: %s, actual: %s", zoneName, searched.DNS[0].Name)
	}

	return searched.DNS[0], nil
}

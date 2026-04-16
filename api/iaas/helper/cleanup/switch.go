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

package cleanup

import (
	"context"
	"fmt"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/query"
	"github.com/sacloud/iaas-api-go/types"
)

// DeleteSwitch 他のリソースから参照されていないかを確認した上で削除する
func DeleteSwitch(ctx context.Context, caller iaas.APICaller, zone string, id types.ID, option query.CheckReferencedOption) error {
	if err := query.WaitWhileSwitchIsReferenced(ctx, caller, zone, id, option); err != nil {
		return fmt.Errorf("switch[%s] is still being used by other resources: %s", id, err)
	}
	return iaas.NewSwitchOp(caller).Delete(ctx, zone, id)
}

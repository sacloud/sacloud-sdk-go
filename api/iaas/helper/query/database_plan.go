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
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/search/keys"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/packages-go/size"
)

type databaseSystemInfoEnvelope struct {
	Products []interface{}   // 利用しない
	Backup   interface{}     // 利用しない
	Plans    []*DatabasePlan `json:"AppliancePlans"`
}

type DatabasePlan struct {
	Class     string
	Model     string
	CPU       int
	MemoryMB  int
	DiskSizes []*DatabaseDiskPlan
}

type DatabaseDiskPlan struct {
	SizeMB       int // 実際のディスクのサイズ? DisplaySizeとは必ずしも一致しない(例: SizeMB:102400, DisplaySize: 90)
	DisplaySize  int // GB単位、コンパネに表示されるのはこの値
	PlanID       types.ID
	ServiceClass string
}

// ListDatabasePlan データベースアプライアンスのプラン情報一覧を取得
//
// modelには以下を指定
//   - Standard	: 標準プラン(非冗長化)
//   - Proxy	: 冗長化プラン
func ListDatabasePlan(ctx context.Context, finder NoteFinder, model string) ([]*DatabasePlan, error) {
	return listDatabasePlan(ctx, finder, model)
}

// GetProxyDatabasePlan 冗長化プランの指定のモデル/CPU/メモリ(GB)/ディスクサイズ(GB)からプランID/サービスクラスを返す
//
// cpu/memoryGB/diskSizeGBに対応するプランが存在しない場合はゼロ値を返す(errorは返さない)
//
// diskSizeGBはDatabaseDiskPlanのDisplaySizeと比較される
func GetProxyDatabasePlan(ctx context.Context, finder NoteFinder, cpu int, memoryGB int, diskSizeGB int) (types.ID, string, error) {
	plans, err := ListDatabasePlan(ctx, finder, "Proxy")
	if err != nil {
		return types.ID(0), "", err
	}
	for _, plan := range plans {
		if plan.CPU == cpu && plan.MemoryMB == memoryGB*size.GiB {
			for _, diskPlan := range plan.DiskSizes {
				if diskPlan.DisplaySize == diskSizeGB {
					return diskPlan.PlanID, diskPlan.ServiceClass, nil
				}
			}
		}
	}
	return types.ID(0), "", nil // plan not found
}

func listDatabasePlan(ctx context.Context, finder NoteFinder, model string) ([]*DatabasePlan, error) {
	if model != "Standard" && model != "Proxy" {
		return nil, fmt.Errorf("unsupported database plan model: %s", model)
	}

	// find note
	searched, err := finder.Find(ctx, &iaas.FindCondition{
		Filter: search.Filter{
			search.Key(keys.Name): "sys-database",
			search.Key("Class"):   "json",
			search.Key("Scope"):   "shared",
		},
	})
	if err != nil {
		return nil, err
	}
	if searched.Count == 0 || len(searched.Notes) == 0 {
		return nil, errors.New("note[sys-database] not found")
	}
	note := searched.Notes[0]

	// parse note's content
	var envelope databaseSystemInfoEnvelope
	if err := json.Unmarshal([]byte(note.Content), &envelope); err != nil {
		return nil, err
	}

	var plans []*DatabasePlan
	for i := range envelope.Plans {
		if envelope.Plans[i].Model == model {
			plans = append(plans, envelope.Plans[i])
		}
	}
	return plans, nil
}

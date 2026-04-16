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
	"errors"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/packages-go/size"
)

// FindServerPlanRequest サーバプラン検索パラメータ
type FindServerPlanRequest struct {
	CPU            int
	MemoryGB       int
	GPU            int
	GPUModel       string
	CPUModel       string
	Commitment     types.ECommitment
	Generation     types.EPlanGeneration
	ConfidentialVM bool
}

func (f *FindServerPlanRequest) findCondition() *iaas.FindCondition {
	confidentialVM := 0
	if f.ConfidentialVM {
		confidentialVM = 1
	}

	cond := &iaas.FindCondition{
		Sort: search.SortKeys{
			{Key: "Generation", Order: search.SortDesc},
		},
		Filter: search.Filter{
			search.Key("Commitment"):     types.Commitments.Standard,
			search.Key("ConfidentialVM"): confidentialVM, // HACK: bool型は0/1として指定する必要がある模様
		},
		Count: 1000,
	}
	if f.CPU > 0 {
		cond.Filter[search.Key("CPU")] = f.CPU
	}
	if f.MemoryGB > 0 {
		cond.Filter[search.Key("MemoryMB")] = size.GiBToMiB(f.MemoryGB)
	}

	if f.GPU > 0 {
		cond.Filter[search.Key("GPU")] = f.GPU
	}
	if f.Generation != types.PlanGenerations.Default {
		cond.Filter[search.Key("Generation")] = f.Generation
	}
	if f.Commitment != types.Commitments.Unknown && f.Commitment != types.Commitments.Standard {
		cond.Filter[search.Key("Commitment")] = f.Commitment
	}
	if f.CPUModel != "" {
		cond.Filter[search.Key("CPUModel")] = f.CPUModel
	}
	if f.GPUModel != "" {
		cond.Filter[search.Key("GPUModel")] = f.GPUModel
	}
	return cond
}

// FindServerPlan サーバプラン検索
func FindServerPlan(ctx context.Context, finder ServerPlanFinder, zone string, param *FindServerPlanRequest) (*iaas.ServerPlan, error) {
	var cond *iaas.FindCondition
	if param != nil {
		cond = param.findCondition()
	}

	searched, err := finder.Find(ctx, zone, cond)
	if err != nil {
		return nil, err
	}
	if searched.Count == 0 || len(searched.ServerPlans) == 0 {
		return nil, errors.New("server plan not found")
	}
	return searched.ServerPlans[0], nil
}

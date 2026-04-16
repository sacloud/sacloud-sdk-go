// Copyright 2021-2024 The sacloud/apprun-api-go authors
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

package fake

import v1 "github.com/sacloud/apprun-api-go/apis/v1"

func (engine *Engine) ListTraffics(appId string) (*v1.HandlerListTraffics, error) {
	if _, ok := engine.Traffics[appId]; !ok {
		return nil, newError(
			ErrorTypeNotFound, "traffic", nil,
			"アプリケーションが見つかりませんでした。")
	}

	var ts []v1.Traffic
	for _, t := range engine.Traffics[appId] {
		ts = append(ts, *t)
	}

	return &v1.HandlerListTraffics{
		Meta: nil,
		Data: ts,
	}, nil
}

func (engine *Engine) UpdateTraffic(appId string, body *v1.PutTrafficsBody) (*v1.HandlerPutTraffics, error) {
	if _, ok := engine.Traffics[appId]; !ok {
		return nil, newError(
			ErrorTypeNotFound, "traffic", nil,
			"アプリケーションが見つかりませんでした。")
	}

	var ts []*v1.Traffic
	var data []v1.Traffic
	total := 0
	for _, v := range *body {
		withVersion, _ := v.AsTrafficWithVersionName()
		withLatest, _ := v.AsTrafficWithLatestVersion()

		tr := &v1.Traffic{}
		if withVersion.VersionName != "" && withVersion.Percent > 0 {
			total += withVersion.Percent
			if err := tr.FromTrafficWithVersionName(withVersion); err != nil {
				return nil, err
			}
		} else {
			total += withLatest.Percent
			if err := tr.FromTrafficWithLatestVersion(withLatest); err != nil {
				return nil, err
			}
		}

		ts = append(ts, tr)
		data = append(data, *tr)
	}

	if total != 100 {
		return nil, newError(
			ErrorTypeInvalidRequest, "traffic", nil,
			"トラフィック分散の割合が合計100になりません")
	}

	engine.Traffics[appId] = ts

	return &v1.HandlerPutTraffics{
		Data: data,
		Meta: nil,
	}, nil
}

func (engine *Engine) initTraffic(app *v1.Application) {
	withLatest := &v1.Traffic{}
	if err := withLatest.FromTrafficWithLatestVersion(v1.TrafficWithLatestVersion{
		IsLatestVersion: true,
		Percent:         100,
	}); err != nil {
		panic(err)
	}

	// 内部的にTrafficとApplicationのリレーションを保持する
	engine.Traffics[app.Id] = append(engine.Traffics[app.Id], withLatest)
}

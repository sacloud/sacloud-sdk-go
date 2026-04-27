// Copyright 2021-2026 The sacloud/apprun-api-go authors
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

	params := v1.HandlerListTrafficsMeta{}
	return &v1.HandlerListTraffics{
		Meta: v1.NewOptHandlerListTrafficsMeta(&params),
		Data: engine.Traffics[appId],
	}, nil
}

func (engine *Engine) UpdateTraffic(appId string, body *v1.PutTrafficsBody) (*v1.HandlerPutTraffics, error) {
	if _, ok := engine.Traffics[appId]; !ok {
		return nil, newError(
			ErrorTypeNotFound, "traffic", nil,
			"アプリケーションが見つかりませんでした。")
	}

	var listData []v1.HandlerListTrafficsDataItem
	var putData []v1.HandlerPutTrafficsDataItem
	total := 0
	for _, v := range *body {
		if item, ok := v.GetPutTrafficsBodyItem0(); ok {
			total += item.Percent
			listData = append(listData, v1.HandlerListTrafficsDataItem{
				IsLatestVersion: item.IsLatestVersion,
				Percent:         item.Percent,
			})
			putData = append(putData, v1.HandlerPutTrafficsDataItem{
				IsLatestVersion: item.IsLatestVersion,
				Percent:         item.Percent,
			})
			continue
		}

		if item, ok := v.GetPutTrafficsBodyItem1(); ok {
			total += item.Percent
			listData = append(listData, v1.HandlerListTrafficsDataItem{
				VersionName: item.VersionName,
				Percent:     item.Percent,
			})
			putData = append(putData, v1.HandlerPutTrafficsDataItem{
				VersionName: item.VersionName,
				Percent:     item.Percent,
			})
		}
	}

	if total != 100 {
		return nil, newError(
			ErrorTypeInvalidRequest, "traffic", nil,
			"トラフィック分散の割合が合計100になりません")
	}

	engine.Traffics[appId] = listData
	params := v1.HandlerPutTrafficsMeta{}
	return &v1.HandlerPutTraffics{
		Data: putData,
		Meta: &params,
	}, nil
}

func (engine *Engine) initTraffic(app *v1.HandlerGetApplication) {
	// 内部的にTrafficとApplicationのリレーションを保持する
	engine.Traffics[app.ID] = append(engine.Traffics[app.ID], v1.HandlerListTrafficsDataItem{
		IsLatestVersion: true,
		Percent:         100,
	})
}

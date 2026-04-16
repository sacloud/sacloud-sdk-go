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

func (engine *Engine) ReadPacketFilter(appId string) (*v1.HandlerGetPacketFilter, error) {
	defer engine.rLock()()

	v, ok := engine.appPacketFilterRelations[appId]
	if !ok {
		return nil, newError(
			ErrorTypeNotFound, "packet_filter", nil,
			"アプリケーション、またはパケットフィルタが見つかりませんでした。")
	}

	return &v1.HandlerGetPacketFilter{
		IsEnabled: v.IsEnabled,
		Settings:  v.Settings,
	}, nil
}

func (engine *Engine) UpdatePacketFilter(appId string, body *v1.PatchPacketFilter) (*v1.HandlerPatchPacketFilter, error) {
	if _, err := engine.ReadApplication(appId); err != nil {
		return nil, newError(
			ErrorTypeNotFound, "application", nil,
			"アプリケーションが見つかりませんでした。")
	}

	v := &v1.HandlerGetPacketFilter{}
	if body.IsEnabled != nil {
		v.IsEnabled = *body.IsEnabled
	}
	if body.Settings != nil {
		v.Settings = *body.Settings
	}
	engine.appPacketFilterRelations[appId] = v

	return &v1.HandlerPatchPacketFilter{
		IsEnabled: v.IsEnabled,
		Settings:  v.Settings,
	}, nil
}

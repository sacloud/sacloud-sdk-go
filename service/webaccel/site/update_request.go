// Copyright 2022-2025 The sacloud/webaccel-service-go authors
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

package site

import (
	"github.com/sacloud/webaccel-api-go"
)

type UpdateRequest struct {
	Id string `service:"-" validate:"required"`

	// Note: 本来serviceのUpdateパラメータはポインタ型で各値を受け取り、Readで取得した現在の値とマージして更新用パラメータを作成する。
	// しかしWebAccelのSiteについてはAPI側がPATCH動作(指定された値だけ更新する)なため、マージせずそのままパラメータ指定するだけで良い。
	webaccel.UpdateSiteRequest
}

func (req *UpdateRequest) ToRequestParameter() *webaccel.UpdateSiteRequest {
	return &req.UpdateSiteRequest
}

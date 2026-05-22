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

package service

import (
	"github.com/sacloud/sacloud-sdk-go/api/webaccel"
	"github.com/sacloud/sacloud-sdk-go/internal/services"
	"github.com/sacloud/sacloud-sdk-go/service/webaccel/cache"
	"github.com/sacloud/sacloud-sdk-go/service/webaccel/site"
	"github.com/sacloud/sacloud-sdk-go/service/webaccel/site/certificate"
	"github.com/sacloud/sacloud-sdk-go/service/webaccel/usage"
)

// Services サービス一覧を返す
func Services(client *webaccel.Client) []services.Service {
	return []services.Service{
		cache.New(client),
		site.New(client),
		certificate.New(client),
		usage.New(client),
	}
}

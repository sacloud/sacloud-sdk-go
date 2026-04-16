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
	"github.com/sacloud/services"
	"github.com/sacloud/webaccel-api-go"
	"github.com/sacloud/webaccel-service-go/cache"
	"github.com/sacloud/webaccel-service-go/site"
	"github.com/sacloud/webaccel-service-go/site/certificate"
	"github.com/sacloud/webaccel-service-go/usage"
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

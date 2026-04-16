// Copyright 2022-2025 The webaccel-api-go authors
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

package webaccel

import "encoding/json"

// MonthlyUsageResults 月別使用量
type MonthlyUsageResults struct {
	Year          int
	Month         int
	MonthlyUsages []*MonthlyUsage
}

type MonthlyUsage struct {
	SiteID             json.Number
	Domain             string
	ASCIIDomain        string
	Subdomain          string
	AccessCount        int64
	BytesSent          int64
	CacheMissBytesSent int64
	CacheHitRatio      float64
	BytesCacheHitRatio float64
	Price              int64
}

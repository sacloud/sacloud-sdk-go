// Copyright 2025- The sacloud/saclient-go Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

// Terraform Provider model shall provide these methods in addition.
type TerraformProviderInterface interface {
	LookupClientConfigProfileName() (string, bool)

	// Returns if a Service Principal ID is set, and its value
	// (this ID is unformatted)
	LookupClientConfigServicePrincipalID() (string, bool)

	// Returns if a Service Principal Key ID is set, and its value
	// (this ID is unformatted)
	LookupClientConfigServicePrincipalKeyID() (string, bool)

	// Returns if a Service Principal BYOK Private Key's path is set, and its value
	LookupClientConfigPrivateKeyPath() (string, bool)

	// Returns if an API key is set, and its value
	LookupClientConfigAccessToken() (string, bool)

	// Returns if an API secret is set, and its value
	LookupClientConfigAccessTokenSecret() (string, bool)

	// Returns if the IaaS target zone is set, and its value
	LookupClientConfigZone() (string, bool)

	// Returns if the "Default" zone is set, and its value.
	LookupClientConfigDefaultZone() (string, bool)

	// Returns if the available zones are set, and their values
	LookupClientConfigZones() ([]string, bool)

	// Returns if the max number of retries is set, and its value
	LookupClientConfigRetryMax() (int64, bool)

	// Returns if the seconds to wait between retries is set, and its value
	LookupClientConfigRetryWaitMax() (int64, bool)

	// Returns if the seconds to wait between retries is set, and its value
	LookupClientConfigRetryWaitMin() (int64, bool)

	// Returns if the API root URL is set, and its value
	LookupClientConfigAPIRootURL() (string, bool)

	// Returns if the API request timeout (in seconds) is set, and its value
	LookupClientConfigAPIRequestTimeout() (int64, bool)

	// Returns if the API request rate limit (requests per second) is set, and its value
	LookupClientConfigAPIRequestRateLimit() (int64, bool)

	// Returns the mysterious "trace" mode value, if set
	LookupClientConfigTraceMode() (string, bool)
}

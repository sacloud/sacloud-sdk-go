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

package saclient

// EndpointConfig contains the resolved input values needed for SDKs to assemble final API endpoint URLs.
type EndpointConfig struct {
	// Endpoints is the primary information for API routing.
	// The map key identifies a service (e.g. "iaas", "iam") and the value is its base endpoint.
	//
	// saclient-go treats service keys as opaque identifiers:
	// keys are normalized to lowercase and are not interpreted or validated here.
	// Their meaning is defined by each SDK.
	Endpoints map[string]string

	// Zone is the currently selected zone (for compatibility only).
	// Interpretation is SDK-specific.
	Zone string

	// Zones lists all permitted zones (for compatibility only).
	// Interpretation is SDK-specific.
	Zones []string

	// APIRootURL is deprecated and should not be used in new code.
	// It is kept for backward compatibility with IaaS consumers.
	APIRootURL string
}

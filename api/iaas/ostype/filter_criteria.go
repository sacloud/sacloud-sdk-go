// Copyright 2022-2025 The sacloud/iaas-api-go Authors
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

package ostype

import (
	"github.com/sacloud/iaas-api-go/search"
	"github.com/sacloud/iaas-api-go/search/keys"
	"github.com/sacloud/iaas-api-go/types"
)

// ArchiveCriteria OSTypeごとのアーカイブ検索条件
var ArchiveCriteria = map[ArchiveOSType]search.Filter{
	AlmaLinux: {
		search.Key(keys.Tags):  search.TagsAndEqual("current-stable", "distro-alma"),
		search.Key(keys.Scope): search.ExactMatch(types.Scopes.Shared.String()),
	},
	AlmaLinux10: {
		search.Key(keys.Tags):  search.TagsAndEqual("alma-10-latest"),
		search.Key(keys.Scope): search.ExactMatch(types.Scopes.Shared.String()),
	},
	AlmaLinux9: {
		search.Key(keys.Tags):  search.TagsAndEqual("alma-9-latest"),
		search.Key(keys.Scope): search.ExactMatch(types.Scopes.Shared.String()),
	},
	AlmaLinux8: {
		search.Key(keys.Tags):  search.TagsAndEqual("alma-8-latest"),
		search.Key(keys.Scope): search.ExactMatch(types.Scopes.Shared.String()),
	},
	RockyLinux: {
		search.Key(keys.Tags):  search.TagsAndEqual("current-stable", "distro-rocky"),
		search.Key(keys.Scope): search.ExactMatch(types.Scopes.Shared.String()),
	},
	RockyLinux10: {
		search.Key(keys.Tags):  search.TagsAndEqual("rocky-10-latest"),
		search.Key(keys.Scope): search.ExactMatch(types.Scopes.Shared.String()),
	},
	RockyLinux9: {
		search.Key(keys.Tags):  search.TagsAndEqual("rocky-9-latest"),
		search.Key(keys.Scope): search.ExactMatch(types.Scopes.Shared.String()),
	},
	RockyLinux8: {
		search.Key(keys.Tags):  search.TagsAndEqual("rocky-8-latest"),
		search.Key(keys.Scope): search.ExactMatch(types.Scopes.Shared.String()),
	},
	MiracleLinux: {
		search.Key(keys.Tags):  search.TagsAndEqual("current-stable", "distro-miracle"),
		search.Key(keys.Scope): search.ExactMatch(types.Scopes.Shared.String()),
	},
	MiracleLinux9: {
		search.Key(keys.Tags):  search.TagsAndEqual("miracle-9-latest"),
		search.Key(keys.Scope): search.ExactMatch(types.Scopes.Shared.String()),
	},
	MiracleLinux8: {
		search.Key(keys.Tags):  search.TagsAndEqual("miracle-8-latest"),
		search.Key(keys.Scope): search.ExactMatch(types.Scopes.Shared.String()),
	},
	Ubuntu: {
		search.Key(keys.Tags):  search.TagsAndEqual("current-stable", "distro-ubuntu"),
		search.Key(keys.Scope): search.ExactMatch(types.Scopes.Shared.String()),
	},
	Ubuntu2404: {
		search.Key(keys.Tags):  search.TagsAndEqual("ubuntu-24.04-latest"),
		search.Key(keys.Scope): search.ExactMatch(types.Scopes.Shared.String()),
	},
	Ubuntu2204: {
		search.Key(keys.Tags):  search.TagsAndEqual("ubuntu-22.04-latest"),
		search.Key(keys.Scope): search.ExactMatch(types.Scopes.Shared.String()),
	},
	Debian: {
		search.Key(keys.Tags):  search.TagsAndEqual("current-stable", "distro-debian"),
		search.Key(keys.Scope): search.ExactMatch(types.Scopes.Shared.String()),
	},
	Debian12: {
		search.Key(keys.Tags):  search.TagsAndEqual("debian-12-latest"),
		search.Key(keys.Scope): search.ExactMatch(types.Scopes.Shared.String()),
	},
	Debian11: {
		search.Key(keys.Tags):  search.TagsAndEqual("debian-11-latest"),
		search.Key(keys.Scope): search.ExactMatch(types.Scopes.Shared.String()),
	},
	Kusanagi: {
		search.Key(keys.Tags):  search.TagsAndEqual("current-stable", "pkg-kusanagi"),
		search.Key(keys.Scope): search.ExactMatch(types.Scopes.Shared.String()),
	},
}

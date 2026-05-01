// Copyright 2025- The sacloud/iam-api-go Authors
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

package iam

import (
	"github.com/sacloud/sacloud-sdk-go/api/iam/apis/auth"
	"github.com/sacloud/sacloud-sdk-go/api/iam/apis/folder"
	"github.com/sacloud/sacloud-sdk-go/api/iam/apis/group"
	"github.com/sacloud/sacloud-sdk-go/api/iam/apis/iampolicy"
	"github.com/sacloud/sacloud-sdk-go/api/iam/apis/iamrole"
	"github.com/sacloud/sacloud-sdk-go/api/iam/apis/idpolicy"
	"github.com/sacloud/sacloud-sdk-go/api/iam/apis/idrole"
	"github.com/sacloud/sacloud-sdk-go/api/iam/apis/organization"
	"github.com/sacloud/sacloud-sdk-go/api/iam/apis/project"
	"github.com/sacloud/sacloud-sdk-go/api/iam/apis/projectapikey"
	"github.com/sacloud/sacloud-sdk-go/api/iam/apis/scim"
	"github.com/sacloud/sacloud-sdk-go/api/iam/apis/servicepolicy"
	"github.com/sacloud/sacloud-sdk-go/api/iam/apis/serviceprincipal"
	"github.com/sacloud/sacloud-sdk-go/api/iam/apis/sso"
	"github.com/sacloud/sacloud-sdk-go/api/iam/apis/user"
	"github.com/sacloud/sacloud-sdk-go/api/iam/apis/user2fa"
)

type AuthAPI = auth.AuthAPI
type FolderAPI = folder.FolderAPI
type GroupAPI = group.GroupAPI
type IAMPolicyAPI = iampolicy.IAMPolicyAPI
type IAMRoleAPI = iamrole.IAMRoleAPI
type IDPolicyAPI = idpolicy.IDPolicyAPI
type IDRoleAPI = idrole.IDRoleAPI
type OrganizationAPI = organization.OrganizationAPI
type ProjectAPI = project.ProjectAPI
type ProjectApiKeyAPI = projectapikey.ProjectAPIKeyAPI
type ScimAPI = scim.ScimAPI
type ServicePolicyAPI = servicepolicy.ServicePolicyAPI
type ServicePrincipalAPI = serviceprincipal.ServicePrincipalAPI
type SSOAPI = sso.SSOAPI
type UserAPI = user.UserAPI
type User2FAAPI = user2fa.User2FAAPI

var NewAuthOp = auth.NewAuthOp
var NewFolderOp = folder.NewFolderOp
var NewGroupOp = group.NewGroupOp
var NewIAMPolicyOp = iampolicy.NewIAMPolicyOp
var NewIAMRoleOp = iamrole.NewIAMRoleOp
var NewIDPolicyOp = idpolicy.NewIDPolicyOp
var NewIDRoleOp = idrole.NewIdRoleOp
var NewOrganizationOp = organization.NewOrganizationOp
var NewProjectOp = project.NewProjectOp
var NewProjectAPIKeyOp = projectapikey.NewProjectAPIKeyOp
var NewScimOp = scim.NewScimOp
var NewServicePolicyOp = servicepolicy.NewServicePolicyOp
var NewServicePrincipalOp = serviceprincipal.NewServicePrincipalOp
var NewSSOOp = sso.NewSSOOp
var NewUserOp = user.NewUserOp
var NewUser2FAOp = user2fa.NewUser2FAOp

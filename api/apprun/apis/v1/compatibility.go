// Copyright 2026- The sacloud/sacloud-sdk-go Authors
// SPDX-License-Identifier: Apache-2.0

package v1

import "context"

// Deprecated: use CreateApplicationBody.
type PostApplicationBody = CreateApplicationBody

// Deprecated: use CreateApplicationBodyComponentsItem.
type PostApplicationBodyComponentsItem = CreateApplicationBodyComponentsItem

// Deprecated: use CreateApplicationBodyComponentsItemDeploySource.
type PostApplicationBodyComponentsItemDeploySource = CreateApplicationBodyComponentsItemDeploySource

// Deprecated: use CreateApplicationBodyComponentsItemDeploySourceContainerRegistry.
type PostApplicationBodyComponentsItemDeploySourceContainerRegistry = CreateApplicationBodyComponentsItemDeploySourceContainerRegistry

// Deprecated: use CreateApplicationBodyComponentsItemMaxCPU.
type PostApplicationBodyComponentsItemMaxCPU = CreateApplicationBodyComponentsItemMaxCPU

const (
	PostApplicationBodyComponentsItemMaxCPU05 = CreateApplicationBodyComponentsItemMaxCPU05
	PostApplicationBodyComponentsItemMaxCPU1  = CreateApplicationBodyComponentsItemMaxCPU1
	PostApplicationBodyComponentsItemMaxCPU2  = CreateApplicationBodyComponentsItemMaxCPU2
)

// Deprecated: use CreateApplicationBodyComponentsItemMaxMemory.
type PostApplicationBodyComponentsItemMaxMemory = CreateApplicationBodyComponentsItemMaxMemory

const (
	PostApplicationBodyComponentsItemMaxMemory1Gi = CreateApplicationBodyComponentsItemMaxMemory1Gi
	PostApplicationBodyComponentsItemMaxMemory2Gi = CreateApplicationBodyComponentsItemMaxMemory2Gi
	PostApplicationBodyComponentsItemMaxMemory4Gi = CreateApplicationBodyComponentsItemMaxMemory4Gi
)

// Deprecated: use CreateApplicationBodyComponentsItemProbe.
type PostApplicationBodyComponentsItemProbe = CreateApplicationBodyComponentsItemProbe

// Deprecated: use CreateApplicationBodyComponentsItemProbeHTTPGet.
type PostApplicationBodyComponentsItemProbeHTTPGet = CreateApplicationBodyComponentsItemProbeHTTPGet

// Deprecated: use CreateApplicationBodyComponentsItemProbeHTTPGetHeadersItem.
type PostApplicationBodyComponentsItemProbeHTTPGetHeadersItem = CreateApplicationBodyComponentsItemProbeHTTPGetHeadersItem

// Deprecated: use RequestEnvItem.
type PostApplicationBodyComponentsItemEnvItem = RequestEnvItem

// Deprecated: use PatchPacketFilterBodySettingsItem.
type PatchPacketFilterSettingsItem = PatchPacketFilterBodySettingsItem

// Deprecated: use HandlerReadUserLimit.
type HandlerGetUserLimit = HandlerReadUserLimit

// Deprecated: use HandlerCreateUserLimit.
type HandlerPostUserLimit = HandlerCreateUserLimit

// Deprecated: use NewOptNilCreateApplicationBodyComponentsItemProbeHTTPGet.
func NewOptNilPostApplicationBodyComponentsItemProbeHTTPGet(v PostApplicationBodyComponentsItemProbeHTTPGet) OptNilCreateApplicationBodyComponentsItemProbeHTTPGet {
	return NewOptNilCreateApplicationBodyComponentsItemProbeHTTPGet(v)
}

// Deprecated: use NewOptCreateApplicationBodyComponentsItemDeploySourceContainerRegistry.
func NewOptPostApplicationBodyComponentsItemDeploySourceContainerRegistry(v PostApplicationBodyComponentsItemDeploySourceContainerRegistry) OptCreateApplicationBodyComponentsItemDeploySourceContainerRegistry {
	return NewOptCreateApplicationBodyComponentsItemDeploySourceContainerRegistry(v)
}

// Deprecated: use NewOptNilRequestEnv.
func NewOptNilPostApplicationBodyComponentsItemEnvItemArray(v []PostApplicationBodyComponentsItemEnvItem) OptNilRequestEnv {
	return NewOptNilRequestEnv(v)
}

// Deprecated: use NewOptNilCreateApplicationBodyComponentsItemProbe.
func NewOptNilPostApplicationBodyComponentsItemProbe(v PostApplicationBodyComponentsItemProbe) OptNilCreateApplicationBodyComponentsItemProbe {
	return NewOptNilCreateApplicationBodyComponentsItemProbe(v)
}

// Deprecated: use PatchPacketFilterBody.
type PatchPacketFilter = PatchPacketFilterBody

// Deprecated: use UpdateTrafficBody.
type PutTrafficsBody = UpdateTrafficBody

// Deprecated: use UpdateTrafficBodyItem.
type PutTrafficsBodyItem = UpdateTrafficBodyItem

// Deprecated: use UpdateTrafficBodyItem0.
type PutTrafficsBodyItem0 = UpdateTrafficBodyItem0

// Deprecated: use UpdateTrafficBodyItem1.
type PutTrafficsBodyItem1 = UpdateTrafficBodyItem1

// Deprecated: use HandlerCreateApplication.
type HandlerPostApplication = HandlerCreateApplication

// Deprecated: use HandlerCreateApplicationStatus.
type HandlerPostApplicationStatus = HandlerCreateApplicationStatus

// Deprecated: use HandlerReadApplication.
type HandlerGetApplication = HandlerReadApplication

// Deprecated: use HandlerReadApplicationStatus.
type HandlerGetApplicationStatus = HandlerReadApplicationStatus

const (
	HandlerGetApplicationStatusHealthy   = HandlerReadApplicationStatusHealthy
	HandlerGetApplicationStatusUnHealthy = HandlerReadApplicationStatusUnHealthy
	HandlerGetApplicationStatusDeploying = HandlerReadApplicationStatusDeploying
)

// Deprecated: use HandlerReadApplicationComponentsItem.
type HandlerGetApplicationComponentsItem = HandlerReadApplicationComponentsItem

// Deprecated: use HandlerReadApplicationOnlyStatus.
type HandlerGetApplicationOnlyStatus = HandlerReadApplicationOnlyStatus

// Deprecated: use HandlerReadApplicationOnlyStatusStatus.
type HandlerGetApplicationOnlyStatusStatus = HandlerReadApplicationOnlyStatusStatus

const (
	HandlerGetApplicationOnlyStatusStatusHealthy   = HandlerReadApplicationOnlyStatusStatusHealthy
	HandlerGetApplicationOnlyStatusStatusUnHealthy = HandlerReadApplicationOnlyStatusStatusUnHealthy
	HandlerGetApplicationOnlyStatusStatusDeploying = HandlerReadApplicationOnlyStatusStatusDeploying
)

// Deprecated: use HandlerReadVersion.
type HandlerGetVersion = HandlerReadVersion

// Deprecated: use HandlerReadVersionComponentsItem.
type HandlerGetVersionComponentsItem = HandlerReadVersionComponentsItem

// Deprecated: use HandlerReadVersionStatus.
type HandlerGetVersionStatus = HandlerReadVersionStatus

// Deprecated: use HandlerReadApplicationVersionOnlyStatus.
type HandlerGetApplicationVersionOnlyStatus = HandlerReadApplicationVersionOnlyStatus

// Deprecated: use HandlerReadApplicationVersionOnlyStatusStatus.
type HandlerGetApplicationVersionOnlyStatusStatus = HandlerReadApplicationVersionOnlyStatusStatus

const (
	HandlerGetApplicationVersionOnlyStatusStatusHealthy   = HandlerReadApplicationVersionOnlyStatusStatusHealthy
	HandlerGetApplicationVersionOnlyStatusStatusUnHealthy = HandlerReadApplicationVersionOnlyStatusStatusUnHealthy
	HandlerGetApplicationVersionOnlyStatusStatusDeploying = HandlerReadApplicationVersionOnlyStatusStatusDeploying
)

// Deprecated: use HandlerReadPacketFilter.
type HandlerGetPacketFilter = HandlerReadPacketFilter

// Deprecated: use HandlerReadPacketFilterSettingsItem.
type HandlerGetPacketFilterSettingsItem = HandlerReadPacketFilterSettingsItem

// Deprecated: use HandlerReadUser.
type HandlerGetUser = HandlerReadUser

// Deprecated: use HandlerCreateUser.
type HandlerPostUser = HandlerCreateUser

// Deprecated: use HandlerListTraffic.
type HandlerListTraffics = HandlerListTraffic

// Deprecated: use HandlerListTrafficDataItem.
type HandlerListTrafficsDataItem = HandlerListTrafficDataItem

// Deprecated: use HandlerListTrafficMeta.
type HandlerListTrafficsMeta = HandlerListTrafficMeta

// Deprecated: use HandlerUpdateTraffic.
type HandlerPutTraffics = HandlerUpdateTraffic

// Deprecated: use HandlerUpdateTrafficDataItem.
type HandlerPutTrafficsDataItem = HandlerUpdateTrafficDataItem

// Deprecated: use HandlerUpdateTrafficMeta.
type HandlerPutTrafficsMeta = HandlerUpdateTrafficMeta

// Deprecated: use NewOptHandlerListTrafficMeta.
func NewOptHandlerListTrafficsMeta(v *HandlerListTrafficsMeta) OptHandlerListTrafficMeta {
	return NewOptHandlerListTrafficMeta(v)
}

// Deprecated: use NewUpdateTrafficBodyItem0UpdateTrafficBodyItem.
func NewPutTrafficsBodyItem0PutTrafficsBodyItem(v PutTrafficsBodyItem0) PutTrafficsBodyItem {
	return NewUpdateTrafficBodyItem0UpdateTrafficBodyItem(v)
}

// Deprecated: use NewUpdateTrafficBodyItem1UpdateTrafficBodyItem.
func NewPutTrafficsBodyItem1PutTrafficsBodyItem(v PutTrafficsBodyItem1) PutTrafficsBodyItem {
	return NewUpdateTrafficBodyItem1UpdateTrafficBodyItem(v)
}

// Deprecated: use CreateApplicationRes.
type PostApplicationRes = CreateApplicationRes

// Deprecated: use ReadApplicationRes.
type GetApplicationRes = ReadApplicationRes

// Deprecated: use ReadApplicationStatusRes.
type GetApplicationStatusRes = ReadApplicationStatusRes

// Deprecated: use ReadApplicationVersionRes.
type GetApplicationVersionRes = ReadApplicationVersionRes

// Deprecated: use ReadApplicationVersionStatusRes.
type GetApplicationVersionStatusRes = ReadApplicationVersionStatusRes

// Deprecated: use ReadApplicationPacketFilterRes.
type GetPacketFilterRes = ReadApplicationPacketFilterRes

// Deprecated: use PatchApplicationPacketFilterRes.
type PatchPacketFilterRes = PatchApplicationPacketFilterRes

// Deprecated: use ListApplicationTrafficRes.
type ListApplicationTrafficsRes = ListApplicationTrafficRes

// Deprecated: use UpdateApplicationTrafficRes.
type PutApplicationTrafficRes = UpdateApplicationTrafficRes

// Deprecated: use ReadUserRes.
type GetUserRes = ReadUserRes

// Deprecated: use CreateUserRes.
type PostUserRes = CreateUserRes

type (
	PostApplicationBadRequest          = CreateApplicationBadRequest
	PostApplicationUnauthorized        = CreateApplicationUnauthorized
	PostApplicationForbidden           = CreateApplicationForbidden
	PostApplicationConflict            = CreateApplicationConflict
	PostApplicationInternalServerError = CreateApplicationInternalServerError

	GetApplicationBadRequest          = ReadApplicationBadRequest
	GetApplicationUnauthorized        = ReadApplicationUnauthorized
	GetApplicationForbidden           = ReadApplicationForbidden
	GetApplicationNotFound            = ReadApplicationNotFound
	GetApplicationInternalServerError = ReadApplicationInternalServerError

	GetApplicationStatusBadRequest          = ReadApplicationStatusBadRequest
	GetApplicationStatusUnauthorized        = ReadApplicationStatusUnauthorized
	GetApplicationStatusForbidden           = ReadApplicationStatusForbidden
	GetApplicationStatusNotFound            = ReadApplicationStatusNotFound
	GetApplicationStatusInternalServerError = ReadApplicationStatusInternalServerError

	GetApplicationVersionBadRequest          = ReadApplicationVersionBadRequest
	GetApplicationVersionUnauthorized        = ReadApplicationVersionUnauthorized
	GetApplicationVersionForbidden           = ReadApplicationVersionForbidden
	GetApplicationVersionNotFound            = ReadApplicationVersionNotFound
	GetApplicationVersionInternalServerError = ReadApplicationVersionInternalServerError

	GetApplicationVersionStatusBadRequest          = ReadApplicationVersionStatusBadRequest
	GetApplicationVersionStatusUnauthorized        = ReadApplicationVersionStatusUnauthorized
	GetApplicationVersionStatusForbidden           = ReadApplicationVersionStatusForbidden
	GetApplicationVersionStatusNotFound            = ReadApplicationVersionStatusNotFound
	GetApplicationVersionStatusInternalServerError = ReadApplicationVersionStatusInternalServerError

	GetPacketFilterBadRequest          = ReadApplicationPacketFilterBadRequest
	GetPacketFilterUnauthorized        = ReadApplicationPacketFilterUnauthorized
	GetPacketFilterForbidden           = ReadApplicationPacketFilterForbidden
	GetPacketFilterNotFound            = ReadApplicationPacketFilterNotFound
	GetPacketFilterInternalServerError = ReadApplicationPacketFilterInternalServerError

	PatchPacketFilterBadRequest          = PatchApplicationPacketFilterBadRequest
	PatchPacketFilterUnauthorized        = PatchApplicationPacketFilterUnauthorized
	PatchPacketFilterForbidden           = PatchApplicationPacketFilterForbidden
	PatchPacketFilterNotFound            = PatchApplicationPacketFilterNotFound
	PatchPacketFilterInternalServerError = PatchApplicationPacketFilterInternalServerError

	ListApplicationTrafficsBadRequest          = ListApplicationTrafficBadRequest
	ListApplicationTrafficsUnauthorized        = ListApplicationTrafficUnauthorized
	ListApplicationTrafficsForbidden           = ListApplicationTrafficForbidden
	ListApplicationTrafficsNotFound            = ListApplicationTrafficNotFound
	ListApplicationTrafficsInternalServerError = ListApplicationTrafficInternalServerError

	PutApplicationTrafficBadRequest          = UpdateApplicationTrafficBadRequest
	PutApplicationTrafficUnauthorized        = UpdateApplicationTrafficUnauthorized
	PutApplicationTrafficForbidden           = UpdateApplicationTrafficForbidden
	PutApplicationTrafficNotFound            = UpdateApplicationTrafficNotFound
	PutApplicationTrafficInternalServerError = UpdateApplicationTrafficInternalServerError

	GetUserUnauthorized        = ReadUserUnauthorized
	GetUserForbidden           = ReadUserForbidden
	GetUserNotFound            = ReadUserNotFound
	GetUserInternalServerError = ReadUserInternalServerError

	PostUserUnauthorized        = CreateUserUnauthorized
	PostUserForbidden           = CreateUserForbidden
	PostUserConflict            = CreateUserConflict
	PostUserInternalServerError = CreateUserInternalServerError
)

// Deprecated: use ReadApplicationParams.
type GetApplicationParams = ReadApplicationParams

// Deprecated: use ReadApplicationStatusParams.
type GetApplicationStatusParams = ReadApplicationStatusParams

// Deprecated: use ReadApplicationVersionParams.
type GetApplicationVersionParams = ReadApplicationVersionParams

// Deprecated: use ReadApplicationVersionStatusParams.
type GetApplicationVersionStatusParams = ReadApplicationVersionStatusParams

// Deprecated: use ReadApplicationPacketFilterParams.
type GetPacketFilterParams = ReadApplicationPacketFilterParams

// Deprecated: use PatchApplicationPacketFilterParams.
type PatchPacketFilterParams = PatchApplicationPacketFilterParams

// Deprecated: use ListApplicationTrafficParams.
type ListApplicationTrafficsParams = ListApplicationTrafficParams

// Deprecated: use UpdateApplicationTrafficParams.
type PutApplicationTrafficParams = UpdateApplicationTrafficParams

// Deprecated: use CreateApplication.
func (c *Client) PostApplication(ctx context.Context, request *PostApplicationBody) (PostApplicationRes, error) {
	return c.CreateApplication(ctx, request)
}

// Deprecated: use ReadApplication.
func (c *Client) GetApplication(ctx context.Context, params GetApplicationParams) (GetApplicationRes, error) {
	return c.ReadApplication(ctx, params)
}

// Deprecated: use ReadApplicationStatus.
func (c *Client) GetApplicationStatus(ctx context.Context, params GetApplicationStatusParams) (GetApplicationStatusRes, error) {
	return c.ReadApplicationStatus(ctx, params)
}

// Deprecated: use ReadApplicationVersion.
func (c *Client) GetApplicationVersion(ctx context.Context, params GetApplicationVersionParams) (GetApplicationVersionRes, error) {
	return c.ReadApplicationVersion(ctx, params)
}

// Deprecated: use ReadApplicationVersionStatus.
func (c *Client) GetApplicationVersionStatus(ctx context.Context, params GetApplicationVersionStatusParams) (GetApplicationVersionStatusRes, error) {
	return c.ReadApplicationVersionStatus(ctx, params)
}

// Deprecated: use ReadApplicationPacketFilter.
func (c *Client) GetPacketFilter(ctx context.Context, params GetPacketFilterParams) (GetPacketFilterRes, error) {
	return c.ReadApplicationPacketFilter(ctx, params)
}

// Deprecated: use PatchApplicationPacketFilter.
func (c *Client) PatchPacketFilter(ctx context.Context, request *PatchPacketFilter, params PatchPacketFilterParams) (PatchPacketFilterRes, error) {
	return c.PatchApplicationPacketFilter(ctx, request, params)
}

// Deprecated: use ListApplicationTraffic.
func (c *Client) ListApplicationTraffics(ctx context.Context, params ListApplicationTrafficsParams) (ListApplicationTrafficsRes, error) {
	return c.ListApplicationTraffic(ctx, params)
}

// Deprecated: use UpdateApplicationTraffic.
func (c *Client) PutApplicationTraffic(ctx context.Context, request PutTrafficsBody, params PutApplicationTrafficParams) (PutApplicationTrafficRes, error) {
	return c.UpdateApplicationTraffic(ctx, request, params)
}

// Deprecated: use ReadUser.
func (c *Client) GetUser(ctx context.Context) (GetUserRes, error) {
	return c.ReadUser(ctx)
}

// Deprecated: use CreateUser.
func (c *Client) PostUser(ctx context.Context) (PostUserRes, error) {
	return c.CreateUser(ctx)
}

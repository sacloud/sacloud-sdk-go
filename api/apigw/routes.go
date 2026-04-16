// Copyright 2025- The sacloud/apigw-api-go authors
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

package apigw

import (
	"context"
	"errors"

	"github.com/google/uuid"
	v1 "github.com/sacloud/apigw-api-go/apis/v1"
)

type RouteAPI interface {
	List(ctx context.Context) ([]v1.Route, error)
	Create(ctx context.Context, request *v1.RouteDetail) (*v1.RouteDetail, error)
	Read(ctx context.Context, id uuid.UUID) (*v1.RouteDetail, error)
	Update(ctx context.Context, request *v1.RouteDetail, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

var _ RouteAPI = (*routeOp)(nil)

type routeOp struct {
	client    *v1.Client
	serviceId uuid.UUID
}

func NewRouteOp(client *v1.Client, serviceId uuid.UUID) RouteAPI {
	return &routeOp{client: client, serviceId: serviceId}
}

func (op *routeOp) List(ctx context.Context) ([]v1.Route, error) {
	res, err := op.client.GetServiceRoutes(ctx, v1.GetServiceRoutesParams{ServiceId: op.serviceId})
	if err != nil {
		return nil, NewAPIError("Route.List", 0, err)
	}

	switch p := res.(type) {
	case *v1.GetServiceRoutesOK:
		return p.Apigw.Routes, nil
	case *v1.GetServiceRoutesBadRequest:
		return nil, NewAPIError("Route.List", 400, errors.New(p.Message.Value))
	case *v1.GetServiceRoutesNotFound:
		return nil, NewAPIError("Route.List", 404, errors.New(p.Message.Value))
	case *v1.GetServiceRoutesInternalServerError:
		return nil, NewAPIError("Route.List", 500, errors.New(p.Message.Value))
	}

	return nil, NewAPIError("Route.List", 0, nil)
}

func (op *routeOp) Create(ctx context.Context, request *v1.RouteDetail) (*v1.RouteDetail, error) {
	// ogenが現状arrayに対するdefaultsをサポートしてないので、代わりに実装する
	if len(request.Methods) == 0 {
		request.Methods = v1.HTTPMethodGET.AllValues()
	}

	res, err := op.client.AddRoute(ctx, request, v1.AddRouteParams{ServiceId: op.serviceId})
	if err != nil {
		return nil, NewAPIError("Route.Create", 0, err)
	}

	switch p := res.(type) {
	case *v1.AddRouteCreated:
		return &p.Apigw.Route.Value, nil
	case *v1.AddRouteBadRequest:
		return nil, NewAPIError("Route.Create", 400, errors.New(p.Message.Value))
	case *v1.AddRouteNotFound:
		return nil, NewAPIError("Route.Create", 404, errors.New(p.Message.Value))
	case *v1.AddRouteConflict:
		return nil, NewAPIError("Route.Create", 409, errors.New(p.Message.Value))
	case *v1.AddRouteInternalServerError:
		return nil, NewAPIError("Route.Create", 500, errors.New(p.Message.Value))
	}

	return nil, NewAPIError("Route.Create", 0, nil)
}

func (op *routeOp) Read(ctx context.Context, id uuid.UUID) (*v1.RouteDetail, error) {
	res, err := op.client.GetRoute(ctx, v1.GetRouteParams{ServiceId: op.serviceId, RouteId: id})
	if err != nil {
		return nil, NewAPIError("Route.Read", 0, err)
	}

	switch p := res.(type) {
	case *v1.GetRouteOK:
		return &p.Apigw.Route.Value, nil
	case *v1.GetRouteBadRequest:
		return nil, NewAPIError("Route.Read", 400, errors.New(p.Message.Value))
	case *v1.GetRouteNotFound:
		return nil, NewAPIError("Route.Read", 404, errors.New(p.Message.Value))
	case *v1.GetRouteInternalServerError:
		return nil, NewAPIError("Route.Read", 500, errors.New(p.Message.Value))
	}

	return nil, NewAPIError("Route.Read", 0, nil)
}

func (op *routeOp) Update(ctx context.Context, request *v1.RouteDetail, id uuid.UUID) error {
	res, err := op.client.UpdateRoute(ctx, request, v1.UpdateRouteParams{ServiceId: op.serviceId, RouteId: id})
	if err != nil {
		return NewAPIError("Route.Update", 0, err)
	}

	switch p := res.(type) {
	case *v1.UpdateRouteNoContent:
		return nil
	case *v1.UpdateRouteBadRequest:
		return NewAPIError("Route.Update", 400, errors.New(p.Message.Value))
	case *v1.UpdateRouteNotFound:
		return NewAPIError("Route.Update", 404, errors.New(p.Message.Value))
	case *v1.UpdateRouteConflict:
		return NewAPIError("Route.Update", 409, errors.New(p.Message.Value))
	case *v1.UpdateRouteInternalServerError:
		return NewAPIError("Route.Update", 500, errors.New(p.Message.Value))
	}

	return NewAPIError("Route.Update", 0, nil)
}

func (op *routeOp) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := op.client.DeleteRoute(ctx, v1.DeleteRouteParams{ServiceId: op.serviceId, RouteId: id})
	if err != nil {
		return NewAPIError("Route.Delete", 0, err)
	}

	switch p := res.(type) {
	case *v1.DeleteRouteNoContent:
		return nil
	case *v1.DeleteRouteBadRequest:
		return NewAPIError("Route.Delete", 400, errors.New(p.Message.Value))
	case *v1.DeleteRouteUnauthorized:
		return NewAPIError("Route.Delete", 401, errors.New(p.Message.Value))
	case *v1.DeleteRouteNotFound:
		return NewAPIError("Route.Delete", 404, errors.New(p.Message.Value))
	case *v1.DeleteRouteInternalServerError:
		return NewAPIError("Route.Delete", 500, errors.New(p.Message.Value))
	}

	return NewAPIError("Route.Delete", 0, nil)
}

type RouteExtraAPI interface {
	ReadAuthorization(ctx context.Context) (*v1.RouteAuthorizationDetailResponse, error)
	DisableAuthorization(ctx context.Context) error
	EnableAuthorization(ctx context.Context, groups []v1.RouteAuthorization) error
	ReadRequestTransformation(ctx context.Context) (*v1.RequestTransformation, error)
	UpdateRequestTransformation(ctx context.Context, request *v1.RequestTransformation) error
	ReadResponseTransformation(ctx context.Context) (*v1.ResponseTransformation, error)
	UpdateResponseTransformation(ctx context.Context, request *v1.ResponseTransformation) error
}

var _ RouteExtraAPI = (*routeExtraOp)(nil)

type routeExtraOp struct {
	client    *v1.Client
	serviceId uuid.UUID
	routeId   uuid.UUID
}

func NewRouteExtraOp(client *v1.Client, serviceId uuid.UUID, routeId uuid.UUID) RouteExtraAPI {
	return &routeExtraOp{client: client, serviceId: serviceId, routeId: routeId}
}

func (op *routeExtraOp) ReadAuthorization(ctx context.Context) (*v1.RouteAuthorizationDetailResponse, error) {
	res, err := op.client.GetRouteAuthorization(ctx, v1.GetRouteAuthorizationParams{
		ServiceId: op.serviceId, RouteId: op.routeId})
	if err != nil {
		return nil, NewAPIError("RouteExtra.ReadAuthorization", 0, err)
	}

	switch p := res.(type) {
	case *v1.GetRouteAuthorizationOK:
		return &p.Apigw.RouteAuthorization.Value, nil
	case *v1.GetRouteAuthorizationBadRequest:
		return nil, NewAPIError("RouteExtra.ReadAuthorization", 400, errors.New(p.Message.Value))
	case *v1.GetRouteAuthorizationNotFound:
		return nil, NewAPIError("RouteExtra.ReadAuthorization", 404, errors.New(p.Message.Value))
	case *v1.GetRouteAuthorizationInternalServerError:
		return nil, NewAPIError("RouteExtra.ReadAuthorization", 500, errors.New(p.Message.Value))
	}

	return nil, NewAPIError("RouteExtra.ReadAuthorization", 0, nil)
}

func (op *routeExtraOp) DisableAuthorization(ctx context.Context) error {
	res, err := op.client.UpsertRouteAuthorization(ctx, v1.NewOptRouteAuthorizationDetail(v1.RouteAuthorizationDetail{
		Type:                      v1.RouteAuthorizationDetail0RouteAuthorizationDetail,
		RouteAuthorizationDetail0: v1.RouteAuthorizationDetail0{IsACLEnabled: v1.RouteAuthorizationDetail0IsACLEnabledFalse},
	}), v1.UpsertRouteAuthorizationParams{ServiceId: op.serviceId, RouteId: op.routeId})
	if err != nil {
		return NewAPIError("RouteExtra.DisableAuthorization", 0, err)
	}

	switch p := res.(type) {
	case *v1.UpsertRouteAuthorizationNoContent:
		return nil
	case *v1.UpsertRouteAuthorizationBadRequest:
		return NewAPIError("RouteExtra.DisableAuthorization", 400, errors.New(p.Message.Value))
	case *v1.UpsertRouteAuthorizationNotFound:
		return NewAPIError("RouteExtra.DisableAuthorization", 404, errors.New(p.Message.Value))
	case *v1.UpsertRouteAuthorizationInternalServerError:
		return NewAPIError("RouteExtra.DisableAuthorization", 500, errors.New(p.Message.Value))
	}

	return NewAPIError("RouteExtra.DisableAuthorization", 0, nil)
}

func (op *routeExtraOp) EnableAuthorization(ctx context.Context, groups []v1.RouteAuthorization) error {
	res, err := op.client.UpsertRouteAuthorization(ctx, v1.NewOptRouteAuthorizationDetail(v1.RouteAuthorizationDetail{
		Type: v1.RouteAuthorizationDetail1RouteAuthorizationDetail,
		RouteAuthorizationDetail1: v1.RouteAuthorizationDetail1{
			IsACLEnabled: v1.RouteAuthorizationDetail1IsACLEnabledTrue,
			Groups:       groups,
		},
	}), v1.UpsertRouteAuthorizationParams{ServiceId: op.serviceId, RouteId: op.routeId})
	if err != nil {
		return NewAPIError("RouteExtra.EnableAuthorization", 0, err)
	}

	switch p := res.(type) {
	case *v1.UpsertRouteAuthorizationNoContent:
		return nil
	case *v1.UpsertRouteAuthorizationBadRequest:
		return NewAPIError("RouteExtra.EnableAuthorization", 400, errors.New(p.Message.Value))
	case *v1.UpsertRouteAuthorizationNotFound:
		return NewAPIError("RouteExtra.EnableAuthorization", 404, errors.New(p.Message.Value))
	case *v1.UpsertRouteAuthorizationInternalServerError:
		return NewAPIError("RouteExtra.EnableAuthorization", 500, errors.New(p.Message.Value))
	}

	return NewAPIError("RouteExtra.EnableAuthorization", 0, nil)
}

func (op *routeExtraOp) ReadRequestTransformation(ctx context.Context) (*v1.RequestTransformation, error) {
	res, err := op.client.GetRequestTransformation(ctx, v1.GetRequestTransformationParams{
		ServiceId: op.serviceId, RouteId: op.routeId})
	if err != nil {
		return nil, NewAPIError("RouteExtra.ReadRequestTransformation", 0, err)
	}

	switch p := res.(type) {
	case *v1.GetRequestTransformationOK:
		return &p.Apigw.RequestTransformation.Value, nil
	case *v1.GetRequestTransformationBadRequest:
		return nil, NewAPIError("RouteExtra.ReadRequestTransformation", 400, errors.New(p.Message.Value))
	case *v1.GetRequestTransformationNotFound:
		return nil, NewAPIError("RouteExtra.ReadRequestTransformation", 404, errors.New(p.Message.Value))
	case *v1.GetRequestTransformationInternalServerError:
		return nil, NewAPIError("RouteExtra.ReadRequestTransformation", 500, errors.New(p.Message.Value))
	}

	return nil, NewAPIError("RouteExtra.ReadRequestTransformation", 0, nil)
}

func (op *routeExtraOp) UpdateRequestTransformation(ctx context.Context, request *v1.RequestTransformation) error {
	res, err := op.client.UpsertRequestTransformation(ctx, v1.NewOptRequestTransformation(*request), v1.UpsertRequestTransformationParams{
		ServiceId: op.serviceId, RouteId: op.routeId})
	if err != nil {
		return NewAPIError("RouteExtra.UpdateRequestTransformation", 0, err)
	}

	switch p := res.(type) {
	case *v1.UpsertRequestTransformationNoContent:
		return nil
	case *v1.UpsertRequestTransformationBadRequest:
		return NewAPIError("RouteExtra.UpdateRequestTransformation", 400, errors.New(p.Message.Value))
	case *v1.UpsertRequestTransformationNotFound:
		return NewAPIError("RouteExtra.UpdateRequestTransformation", 404, errors.New(p.Message.Value))
	case *v1.UpsertRequestTransformationInternalServerError:
		return NewAPIError("RouteExtra.UpdateRequestTransformation", 500, errors.New(p.Message.Value))
	}

	return NewAPIError("RouteExtra.UpdateRequestTransformation", 0, nil)
}

func (op *routeExtraOp) ReadResponseTransformation(ctx context.Context) (*v1.ResponseTransformation, error) {
	res, err := op.client.GetResponseTransformation(ctx, v1.GetResponseTransformationParams{
		ServiceId: op.serviceId, RouteId: op.routeId})
	if err != nil {
		return nil, NewAPIError("RouteExtra.ReadResponseTransformation", 0, err)
	}

	switch p := res.(type) {
	case *v1.GetResponseTransformationOK:
		return &p.Apigw.ResponseTransformation.Value, nil
	case *v1.GetResponseTransformationBadRequest:
		return nil, NewAPIError("RouteExtra.ReadResponseTransformation", 400, errors.New(p.Message.Value))
	case *v1.GetResponseTransformationNotFound:
		return nil, NewAPIError("RouteExtra.ReadResponseTransformation", 404, errors.New(p.Message.Value))
	case *v1.GetResponseTransformationInternalServerError:
		return nil, NewAPIError("RouteExtra.ReadResponseTransformation", 500, errors.New(p.Message.Value))
	}

	return nil, NewAPIError("RouteExtra.ReadResponseTransformation", 0, nil)
}

func (op *routeExtraOp) UpdateResponseTransformation(ctx context.Context, request *v1.ResponseTransformation) error {
	res, err := op.client.UpsertResponseTransformation(ctx, v1.NewOptResponseTransformation(*request), v1.UpsertResponseTransformationParams{
		ServiceId: op.serviceId, RouteId: op.routeId})
	if err != nil {
		return NewAPIError("RouteExtra.UpdateResponseTransformation", 0, err)
	}

	switch p := res.(type) {
	case *v1.UpsertResponseTransformationNoContent:
		return nil
	case *v1.UpsertResponseTransformationBadRequest:
		return NewAPIError("RouteExtra.UpdateResponseTransformation", 400, errors.New(p.Message.Value))
	case *v1.UpsertResponseTransformationNotFound:
		return NewAPIError("RouteExtra.UpdateResponseTransformation", 404, errors.New(p.Message.Value))
	case *v1.UpsertResponseTransformationInternalServerError:
		return NewAPIError("RouteExtra.UpdateResponseTransformation", 500, errors.New(p.Message.Value))
	}

	return NewAPIError("RouteExtra.UpdateResponseTransformation", 0, nil)
}

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
	v1 "github.com/sacloud/sacloud-sdk-go/api/apigw/apis/v1"
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
	res, err := errorFromDecodedResponse[v1.GetServiceRoutesOK]("Route.List", func() (any, error) {
		return op.client.GetServiceRoutes(ctx, v1.GetServiceRoutesParams{ServiceId: op.serviceId})
	})

	if err != nil {
		return nil, err
	}

	return new(res.GetApigw()).GetRoutes(), nil
}

func (op *routeOp) Create(ctx context.Context, request *v1.RouteDetail) (*v1.RouteDetail, error) {
	// ogenが現状arrayに対するdefaultsをサポートしてないので、代わりに実装する
	if len(request.Methods) == 0 {
		request.Methods = v1.HTTPMethodGET.AllValues()
	}

	res, err := errorFromDecodedResponse[v1.AddRouteCreated]("Route.Create", func() (any, error) {
		return op.client.AddRoute(ctx, request, v1.AddRouteParams{ServiceId: op.serviceId})
	})
	if err != nil {
		return nil, err
	}

	r, ok := new(res.GetApigw()).GetRoute().Get()
	if !ok {
		return nil, NewError("Route.Create", errors.New("response does not contain a route"))
	}

	return &r, nil
}

func (op *routeOp) Read(ctx context.Context, id uuid.UUID) (*v1.RouteDetail, error) {
	res, err := errorFromDecodedResponse[v1.GetRouteOK]("Route.Read", func() (any, error) {
		return op.client.GetRoute(ctx, v1.GetRouteParams{ServiceId: op.serviceId, RouteId: id})
	})
	if err != nil {
		return nil, err
	}

	r, ok := new(res.GetApigw()).GetRoute().Get()
	if !ok {
		return nil, NewError("Route.Read", errors.New("response does not contain a route"))
	}

	return &r, nil
}

func (op *routeOp) Update(ctx context.Context, request *v1.RouteDetail, id uuid.UUID) error {
	return errorFromDecodedResponseE[v1.UpdateRouteNoContent]("Route.Update", func() (any, error) {
		return op.client.UpdateRoute(ctx, request, v1.UpdateRouteParams{ServiceId: op.serviceId, RouteId: id})
	})
}

func (op *routeOp) Delete(ctx context.Context, id uuid.UUID) error {
	return errorFromDecodedResponseE[v1.DeleteRouteNoContent]("Route.Delete", func() (any, error) {
		return op.client.DeleteRoute(ctx, v1.DeleteRouteParams{ServiceId: op.serviceId, RouteId: id})
	})
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
	params := v1.GetRouteAuthorizationParams{
		ServiceId: op.serviceId,
		RouteId:   op.routeId,
	}

	res, err := errorFromDecodedResponse[v1.GetRouteAuthorizationOK]("RouteExtra.ReadAuthorization", func() (any, error) { return op.client.GetRouteAuthorization(ctx, params) })
	if err != nil {
		return nil, err
	}

	a, ok := new(res.GetApigw()).GetRouteAuthorization().Get()
	if !ok {
		return nil, NewError("RouteExtra.ReadAuthorization", errors.New("response does not contain a route authorization"))
	}

	return &a, nil
}

func (op *routeExtraOp) DisableAuthorization(ctx context.Context) error {
	return errorFromDecodedResponseE[v1.UpsertRouteAuthorizationNoContent]("RouteExtra.DisableAuthorization", func() (any, error) {
		return op.client.UpsertRouteAuthorization(ctx, v1.NewOptRouteAuthorizationDetail(v1.RouteAuthorizationDetail{
			Type:                      v1.RouteAuthorizationDetail0RouteAuthorizationDetail,
			RouteAuthorizationDetail0: v1.RouteAuthorizationDetail0{IsACLEnabled: v1.RouteAuthorizationDetail0IsACLEnabledFalse},
		}), v1.UpsertRouteAuthorizationParams{ServiceId: op.serviceId, RouteId: op.routeId})
	})
}

func (op *routeExtraOp) EnableAuthorization(ctx context.Context, groups []v1.RouteAuthorization) error {
	return errorFromDecodedResponseE[v1.UpsertRouteAuthorizationNoContent]("RouteExtra.EnableAuthorization", func() (any, error) {
		return op.client.UpsertRouteAuthorization(ctx, v1.NewOptRouteAuthorizationDetail(v1.RouteAuthorizationDetail{
			Type: v1.RouteAuthorizationDetail1RouteAuthorizationDetail,
			RouteAuthorizationDetail1: v1.RouteAuthorizationDetail1{
				IsACLEnabled: v1.RouteAuthorizationDetail1IsACLEnabledTrue,
				Groups:       groups,
			},
		}), v1.UpsertRouteAuthorizationParams{ServiceId: op.serviceId, RouteId: op.routeId})
	})
}

func (op *routeExtraOp) ReadRequestTransformation(ctx context.Context) (*v1.RequestTransformation, error) {
	res, err := errorFromDecodedResponse[v1.GetRequestTransformationOK]("RouteExtra.ReadRequestTransformation", func() (any, error) {
		return op.client.GetRequestTransformation(ctx, v1.GetRequestTransformationParams{ServiceId: op.serviceId, RouteId: op.routeId})
	})
	if err != nil {
		return nil, err
	}

	t, ok := new(res.GetApigw()).GetRequestTransformation().Get()
	if !ok {
		return nil, NewError("RouteExtra.ReadRequestTransformation", errors.New("response does not contain a request transformation"))
	}

	return &t, nil
}

func (op *routeExtraOp) UpdateRequestTransformation(ctx context.Context, request *v1.RequestTransformation) error {
	return errorFromDecodedResponseE[v1.UpsertRequestTransformationNoContent]("RouteExtra.UpdateRequestTransformation", func() (any, error) {
		return op.client.UpsertRequestTransformation(ctx, v1.NewOptRequestTransformation(*request), v1.UpsertRequestTransformationParams{ServiceId: op.serviceId, RouteId: op.routeId})
	})
}

func (op *routeExtraOp) ReadResponseTransformation(ctx context.Context) (*v1.ResponseTransformation, error) {
	res, err := errorFromDecodedResponse[v1.GetResponseTransformationOK]("RouteExtra.ReadResponseTransformation", func() (any, error) {
		return op.client.GetResponseTransformation(ctx, v1.GetResponseTransformationParams{ServiceId: op.serviceId, RouteId: op.routeId})
	})
	if err != nil {
		return nil, err
	}

	t, ok := new(res.GetApigw()).GetResponseTransformation().Get()
	if !ok {
		return nil, NewError("RouteExtra.ReadResponseTransformation", errors.New("response does not contain a response transformation"))
	}

	return &t, nil
}

func (op *routeExtraOp) UpdateResponseTransformation(ctx context.Context, request *v1.ResponseTransformation) error {
	return errorFromDecodedResponseE[v1.UpsertResponseTransformationNoContent]("RouteExtra.UpdateResponseTransformation", func() (any, error) {
		return op.client.UpsertResponseTransformation(ctx,
			v1.NewOptResponseTransformation(*request),
			v1.UpsertResponseTransformationParams{ServiceId: op.serviceId, RouteId: op.routeId},
		)
	})
}

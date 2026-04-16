// routes_test.go
package apigw_test

import (
	"context"
	"os"
	"testing"

	apigw "github.com/sacloud/apigw-api-go"
	v1 "github.com/sacloud/apigw-api-go/apis/v1"
	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRouteAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN",
		"SAKURA_ACCESS_TOKEN_SECRET", "SAKURA_TEST_HOST")(t)

	var theClient saclient.Client
	client, err := apigw.NewClient(&theClient)
	require.Nil(t, err)

	subReq, err := getSvcSubRequest()
	require.Nil(t, err)

	ctx := context.Background()
	serviceOp := apigw.NewServiceOp(client)

	service, err := serviceOp.Create(ctx, &v1.ServiceDetailRequest{
		Name:         "test-service",
		Host:         os.Getenv("SAKURA_TEST_HOST"),
		Port:         v1.NewOptInt(80),
		Protocol:     "http",
		Subscription: subReq,
	})
	require.Nil(t, err)
	defer func() { _ = serviceOp.Delete(ctx, service.ID.Value) }()

	routeOp := apigw.NewRouteOp(client, service.ID.Value)

	routeReq := v1.RouteDetail{
		Name:      v1.NewOptName("test-route"),
		Methods:   []v1.HTTPMethod{v1.HTTPMethodGET, v1.HTTPMethodPOST},
		Hosts:     []string{service.RouteHost.Value},
		Protocols: v1.NewOptRouteDetailProtocols(v1.RouteDetailProtocolsHTTPHTTPS),
		Tags:      []string{"Test"},
	}
	createdRoute, err := routeOp.Create(ctx, &routeReq)
	require.Nil(t, err)

	routes, err := routeOp.List(ctx)
	assert.Nil(t, err)
	assert.Greater(t, len(routes), 0)

	readRoute, err := routeOp.Read(ctx, createdRoute.ID.Value)
	assert.Nil(t, err)
	assert.Equal(t, createdRoute.ID.Value, readRoute.ID.Value)

	routeReq.Path = v1.NewOptString("/test-updated")
	err = routeOp.Update(ctx, &routeReq, createdRoute.ID.Value)
	assert.Nil(t, err)

	err = routeOp.Delete(ctx, createdRoute.ID.Value)
	assert.Nil(t, err)
}

func TestRouteExtraAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN",
		"SAKURA_ACCESS_TOKEN_SECRET", "SAKURA_TEST_HOST")(t)

	var theClient saclient.Client
	client, err := apigw.NewClient(&theClient)
	require.Nil(t, err)

	subReq, err := getSvcSubRequest()
	require.Nil(t, err)

	ctx := context.Background()
	serviceOp := apigw.NewServiceOp(client)
	service, err := serviceOp.Create(ctx, &v1.ServiceDetailRequest{
		Name:         "test-service",
		Host:         os.Getenv("SAKURA_TEST_HOST"),
		Port:         v1.NewOptInt(80),
		Protocol:     "http",
		Subscription: subReq,
	})
	require.Nil(t, err)
	defer func() { _ = serviceOp.Delete(ctx, service.ID.Value) }()

	routeOp := apigw.NewRouteOp(client, service.ID.Value)
	route, err := routeOp.Create(ctx, &v1.RouteDetail{
		Name:      v1.NewOptName("test-route"),
		Methods:   []v1.HTTPMethod{v1.HTTPMethodGET, v1.HTTPMethodPOST},
		Hosts:     []string{service.RouteHost.Value},
		Protocols: v1.NewOptRouteDetailProtocols(v1.RouteDetailProtocolsHTTPHTTPS),
		Tags:      []string{"Test"},
	})
	require.Nil(t, err)
	defer func() { _ = routeOp.Delete(ctx, route.ID.Value) }()

	groupOp := apigw.NewGroupOp(client)
	group, err := groupOp.Create(ctx, &v1.Group{Name: v1.NewOptName("test-group"), Tags: []string{"Test"}})
	require.Nil(t, err)
	defer func() { _ = groupOp.Delete(ctx, group.ID.Value) }()

	routeExtraOp := apigw.NewRouteExtraOp(client, service.ID.Value, route.ID.Value)

	err = routeExtraOp.EnableAuthorization(ctx, []v1.RouteAuthorization{{
		ID:      group.ID,
		Name:    group.Name,
		Enabled: v1.NewOptBool(true),
	}})
	assert.Nil(t, err)

	auth, err := routeExtraOp.ReadAuthorization(ctx)
	assert.Nil(t, err)
	assert.True(t, auth.IsACLEnabled)
	assert.Greater(t, len(auth.Groups), 0)

	err = routeExtraOp.DisableAuthorization(ctx)
	assert.Nil(t, err)

	auth2, err := routeExtraOp.ReadAuthorization(ctx)
	if err != nil {
		t.Error(err.Error())
	}
	assert.Nil(t, err)
	assert.False(t, auth2.IsACLEnabled)
	assert.Greater(t, len(auth2.Groups), 0)

	reqTrans := v1.RequestTransformation{
		HttpMethod: v1.NewOptHTTPMethod("GET"),
		Add: v1.NewOptRequestModificationDetail(v1.RequestModificationDetail{
			Headers: []v1.RequestModificationDetailHeadersItem{{
				Key:   v1.NewOptRequestHeaderKey("SDK-Test-Header"),
				Value: v1.NewOptRequestHeaderValue("APIGW-foo"),
			}},
		}),
	}
	err = routeExtraOp.UpdateRequestTransformation(ctx, &reqTrans)
	assert.Nil(t, err)

	resReqTrans, err := routeExtraOp.ReadRequestTransformation(ctx)
	assert.Nil(t, err)
	assert.Equal(t, reqTrans.Add.Value.Headers, resReqTrans.Add.Value.Headers)

	resTrans := v1.ResponseTransformation{
		Allow: v1.NewOptResponseAllowDetail(v1.ResponseAllowDetail{JsonKeys: []v1.JSONKey{"test", "num"}}),
		Rename: v1.NewOptResponseRenameDetail(v1.ResponseRenameDetail{
			IfStatusCode: []int{200},
			Headers: []v1.ResponseRenameDetailHeadersItem{{
				From: v1.NewOptResponseHeaderKey("SDK-From-Key"),
				To:   v1.NewOptResponseHeaderKey("SDK-To-Key"),
			}},
		}),
		Add: v1.NewOptResponseModificationDetail(v1.ResponseModificationDetail{
			IfStatusCode: []int{200, 201, 403},
			JSON: []v1.ResponseModificationDetailJSONItem{{
				Key:   v1.NewOptJSONKey("newKey"),
				Value: v1.NewOptString("newValue"),
			}},
		}),
	}
	err = routeExtraOp.UpdateResponseTransformation(ctx, &resTrans)
	assert.Nil(t, err)

	resResTrans, err := routeExtraOp.ReadResponseTransformation(ctx)
	assert.Nil(t, err)
	assert.Equal(t, resTrans.Allow, resResTrans.Allow)
	assert.Equal(t, resTrans.Rename.Value.Headers, resResTrans.Rename.Value.Headers)
	assert.Equal(t, resTrans.Add.Value.JSON, resResTrans.Add.Value.JSON)
}

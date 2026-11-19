// Copyright 2026- The sacloud/service-endpoint-gateway-api-go authors
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

package seg

import (
	"context"
	"errors"
	"net/http"
	"time"

	v1 "github.com/sacloud/sacloud-sdk-go/api/service-endpoint-gateway/apis/v1"
)

const (
	defaultRetryInterval = 15 * time.Second
	defaultRetryTimeout  = 20 * time.Minute
)

type ServiceEndpointGatewayAPI interface {
	List(ctx context.Context) (*v1.ModelsApplianceListServiceEndpointGatewaysResponseBody, error)
	Create(ctx context.Context, request v1.ModelsApplianceCreateServiceEndpointGatewayRequest) (*v1.ModelsApplianceCreateServiceEndpointGatewayResponseBody, error)
	Read(ctx context.Context, id string) (*v1.ModelsApplianceReadServiceEndpointGatewayResponseBody, error)
	Update(ctx context.Context, id string, request v1.ModelsApplianceUpdateServiceEndpointGatewayRequest) (*v1.ModelsApplianceReadServiceEndpointGatewayResponseBody, error)
	Apply(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
	ReadInterface(ctx context.Context, applianceID string, interfaceID string) (*v1.ModelsApplianceReadServiceEndpointGatewayInterfaceResponseBody, error)
	ReadPowerStatus(ctx context.Context, id string) (*v1.ModelsPowerReadServiceEndpointGatewayPowerStatusResponseBody, error)
	PowerOn(ctx context.Context, id string) (*v1.ModelsPowerUpdateServiceEndpointGatewayPowerStatusResponseBody, error)
	Shutdown(ctx context.Context, id string) error
	Reset(ctx context.Context, id string) error
}

var _ ServiceEndpointGatewayAPI = (*ServiceEndpointGatewayOp)(nil)

type ServiceEndpointGatewayOp struct {
	client        *v1.Client
	powerRetryCfg PowerRetryConfig
}

// PowerRetryConfig is used for power control operations (PowerOn, Shutdown, Reset).
type PowerRetryConfig struct {
	Interval time.Duration
	Timeout  time.Duration
}

func NewServiceEndpointGatewayOp(client *v1.Client) ServiceEndpointGatewayAPI {
	return &ServiceEndpointGatewayOp{client: client}
}

func NewServiceEndpointGatewayOpWithPowerRetryConfig(client *v1.Client, powerRetryCfg PowerRetryConfig) ServiceEndpointGatewayAPI {
	return &ServiceEndpointGatewayOp{client: client, powerRetryCfg: powerRetryCfg}
}

func (op *ServiceEndpointGatewayOp) List(ctx context.Context) (*v1.ModelsApplianceListServiceEndpointGatewaysResponseBody, error) {
	res, err := op.client.ListServiceEndpointGateways(ctx)
	if err != nil {
		if e, ok := errors.AsType[*v1.ModelsCommonDefaultErrorResponseBodyStatusCode](err); ok {
			return nil, NewAPIError("ServiceEndpointGateway.List", e.StatusCode, err)
		}
		return nil, err
	}
	return res, nil
}

func (op *ServiceEndpointGatewayOp) Create(ctx context.Context, request v1.ModelsApplianceCreateServiceEndpointGatewayRequest) (*v1.ModelsApplianceCreateServiceEndpointGatewayResponseBody, error) {
	request.Appliance.Class = v1.ModelsApplianceApplianceCreateBodyClassServiceendpointgateway
	request.Appliance.Plan.ID = v1.ModelsPlanPlanID1

	res, err := op.client.CreateServiceEndpointGateway(ctx, &request)
	if err != nil {
		if e, ok := errors.AsType[*v1.ModelsCommonDefaultErrorResponseBodyStatusCode](err); ok {
			return nil, NewAPIError("ServiceEndpointGateway.Create", e.StatusCode, err)
		}
		return nil, err
	}
	return res, nil
}

func (op *ServiceEndpointGatewayOp) Read(ctx context.Context, id string) (*v1.ModelsApplianceReadServiceEndpointGatewayResponseBody, error) {
	res, err := op.client.ReadServiceEndpointGateway(ctx, v1.ReadServiceEndpointGatewayParams{ApplianceID: id})
	if err != nil {
		if e, ok := errors.AsType[*v1.ModelsCommonDefaultErrorResponseBodyStatusCode](err); ok {
			return nil, NewAPIError("ServiceEndpointGateway.Read", e.StatusCode, err)
		}
		return nil, err
	}
	return res, nil
}

func (op *ServiceEndpointGatewayOp) Update(ctx context.Context, id string, request v1.ModelsApplianceUpdateServiceEndpointGatewayRequest) (*v1.ModelsApplianceReadServiceEndpointGatewayResponseBody, error) {
	res, err := op.client.UpdateServiceEndpointGateway(ctx, &request, v1.UpdateServiceEndpointGatewayParams{ApplianceID: id})
	if err != nil {
		if e, ok := errors.AsType[*v1.ModelsCommonDefaultErrorResponseBodyStatusCode](err); ok {
			return nil, NewAPIError("ServiceEndpointGateway.Update", e.StatusCode, err)
		}
		return nil, err
	}
	return res, nil
}

func (op *ServiceEndpointGatewayOp) Apply(ctx context.Context, id string) error {
	_, err := op.client.ApplyServiceEndpointGatewayConfig(ctx, v1.OptModelsApplianceUpdateServiceEndpointGatewayRequest{Set: false}, v1.ApplyServiceEndpointGatewayConfigParams{ApplianceID: id})
	if err != nil {
		if e, ok := errors.AsType[*v1.ModelsCommonDefaultErrorResponseBodyStatusCode](err); ok {
			return NewAPIError("ServiceEndpointGateway.Apply", e.StatusCode, err)
		}
		return err
	}
	return nil
}

func (op *ServiceEndpointGatewayOp) Delete(ctx context.Context, id string) error {
	_, err := op.client.DeleteServiceEndpointGateway(ctx, v1.DeleteServiceEndpointGatewayParams{ApplianceID: id})
	if err != nil {
		if e, ok := errors.AsType[*v1.ModelsCommonDefaultErrorResponseBodyStatusCode](err); ok {
			return NewAPIError("ServiceEndpointGateway.Delete", e.StatusCode, err)
		}
		return err
	}
	return nil
}

func (op *ServiceEndpointGatewayOp) ReadInterface(ctx context.Context, applianceID string, interfaceID string) (*v1.ModelsApplianceReadServiceEndpointGatewayInterfaceResponseBody, error) {
	res, err := op.client.ReadServiceEndpointGatewayInterface(ctx, v1.ReadServiceEndpointGatewayInterfaceParams{ApplianceID: applianceID, InterfaceID: interfaceID})
	if err != nil {
		if e, ok := errors.AsType[*v1.ModelsCommonDefaultErrorResponseBodyStatusCode](err); ok {
			return nil, NewAPIError("ServiceEndpointGateway.ReadInterface", e.StatusCode, err)
		}
		return nil, err
	}
	return res, nil
}

func (op *ServiceEndpointGatewayOp) ReadPowerStatus(ctx context.Context, id string) (*v1.ModelsPowerReadServiceEndpointGatewayPowerStatusResponseBody, error) {
	res, err := op.client.ReadServiceEndpointGatewayPowerStatus(ctx, v1.ReadServiceEndpointGatewayPowerStatusParams{ApplianceID: id})
	if err != nil {
		if e, ok := errors.AsType[*v1.ModelsCommonDefaultErrorResponseBodyStatusCode](err); ok {
			return nil, NewAPIError("ServiceEndpointGateway.ReadPowerStatus", e.StatusCode, err)
		}
		return nil, err
	}
	return res, nil
}

func (op *ServiceEndpointGatewayOp) PowerOn(ctx context.Context, id string) (*v1.ModelsPowerUpdateServiceEndpointGatewayPowerStatusResponseBody, error) {
	var result *v1.ModelsPowerUpdateServiceEndpointGatewayPowerStatusResponseBody
	err := retryWithTimeout(ctx, op.powerRetryCfg.Interval, op.powerRetryCfg.Timeout, func() (bool, error) {
		res, err := op.client.PowerOnServiceEndpointGateway(ctx, v1.PowerOnServiceEndpointGatewayParams{ApplianceID: id})
		if err != nil {
			var e *v1.ModelsCommonDefaultErrorResponseBodyStatusCode
			if errors.As(err, &e) && e.StatusCode == http.StatusConflict && e.Response.ErrorCode.Value == "still_creating" {
				return true, err
			}
			if errors.As(err, &e) {
				return false, NewAPIError("ServiceEndpointGateway.PowerOn", e.StatusCode, err)
			}
			return false, err
		}
		result = res
		return false, nil
	})
	return result, err
}

func (op *ServiceEndpointGatewayOp) Shutdown(ctx context.Context, id string) error {
	request := v1.NewOptModelsPowerApplianceDeletePowerStatusRequest(v1.ModelsPowerApplianceDeletePowerStatusRequest{Force: true})
	return retryWithTimeout(ctx, op.powerRetryCfg.Interval, op.powerRetryCfg.Timeout, func() (bool, error) {
		_, err := op.client.PowerOffServiceEndpointGateway(ctx, request, v1.PowerOffServiceEndpointGatewayParams{ApplianceID: id})
		if err != nil {
			var e *v1.ModelsCommonDefaultErrorResponseBodyStatusCode
			if errors.As(err, &e) && e.StatusCode == http.StatusConflict && e.Response.ErrorCode.Value == "still_creating" {
				return true, err
			}
			if errors.As(err, &e) {
				return false, NewAPIError("ServiceEndpointGateway.Shutdown", e.StatusCode, err)
			}
			return false, err
		}
		return false, nil
	})
}

func (op *ServiceEndpointGatewayOp) Reset(ctx context.Context, id string) error {
	return retryWithTimeout(ctx, op.powerRetryCfg.Interval, op.powerRetryCfg.Timeout, func() (bool, error) {
		_, err := op.client.ResetServiceEndpointGateway(ctx, v1.ResetServiceEndpointGatewayParams{ApplianceID: id})
		if err != nil {
			var e *v1.ModelsCommonDefaultErrorResponseBodyStatusCode
			if errors.As(err, &e) && e.StatusCode == http.StatusConflict && e.Response.ErrorCode.Value == "still_creating" {
				return true, err
			}
			if errors.As(err, &e) {
				return false, NewAPIError("ServiceEndpointGateway.Reset", e.StatusCode, err)
			}
			return false, err
		}
		return false, nil
	})
}

// retryWithTimeout is a helper to retry a function with timeout and interval.
func retryWithTimeout(ctx context.Context, interval, timeout time.Duration, retryable func() (bool, error)) error {
	interval, timeout = validateRetryParams(interval, timeout)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for {
		retry, err := retryable()
		if err == nil {
			return nil
		}
		if retry {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(interval):
			}
			continue
		}
		return err
	}
}

// validateRetryParams ensures that the interval and timeout values are positive, applying defaults if necessary.
func validateRetryParams(interval, timeout time.Duration) (time.Duration, time.Duration) {
	// interval and timeout should be positive, if not set to default values
	if interval <= 0 {
		interval = defaultRetryInterval
	}
	if timeout <= 0 {
		timeout = defaultRetryTimeout
	}
	return interval, timeout
}

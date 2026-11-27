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
	return errorFromDecodedResponse("ServiceEndpointGateway.List", func() (*v1.ModelsApplianceListServiceEndpointGatewaysResponseBody, error) {
		return op.client.ListServiceEndpointGateways(ctx)
	})
}

func (op *ServiceEndpointGatewayOp) Create(ctx context.Context, request v1.ModelsApplianceCreateServiceEndpointGatewayRequest) (*v1.ModelsApplianceCreateServiceEndpointGatewayResponseBody, error) {
	request.Appliance.Class = v1.ModelsApplianceApplianceCreateBodyClassServiceendpointgateway
	request.Appliance.Plan.ID = v1.ModelsPlanPlanID1

	return errorFromDecodedResponse("ServiceEndpointGateway.Create", func() (*v1.ModelsApplianceCreateServiceEndpointGatewayResponseBody, error) {
		return op.client.CreateServiceEndpointGateway(ctx, &request)
	})
}

func (op *ServiceEndpointGatewayOp) Read(ctx context.Context, id string) (*v1.ModelsApplianceReadServiceEndpointGatewayResponseBody, error) {
	return errorFromDecodedResponse("ServiceEndpointGateway.Read", func() (*v1.ModelsApplianceReadServiceEndpointGatewayResponseBody, error) {
		return op.client.ReadServiceEndpointGateway(ctx, v1.ReadServiceEndpointGatewayParams{ApplianceID: id})
	})
}

func (op *ServiceEndpointGatewayOp) Update(ctx context.Context, id string, request v1.ModelsApplianceUpdateServiceEndpointGatewayRequest) (*v1.ModelsApplianceReadServiceEndpointGatewayResponseBody, error) {
	return errorFromDecodedResponse("ServiceEndpointGateway.Update", func() (*v1.ModelsApplianceReadServiceEndpointGatewayResponseBody, error) {
		return op.client.UpdateServiceEndpointGateway(ctx, &request, v1.UpdateServiceEndpointGatewayParams{ApplianceID: id})
	})
}

func (op *ServiceEndpointGatewayOp) Apply(ctx context.Context, id string) error {
	return errorFromDecodedResponseE("ServiceEndpointGateway.Apply", func() (*v1.ModelsApplianceApplianceApplyResponseBody, error) {
		return op.client.ApplyServiceEndpointGatewayConfig(ctx, v1.OptModelsApplianceUpdateServiceEndpointGatewayRequest{Set: false}, v1.ApplyServiceEndpointGatewayConfigParams{ApplianceID: id})
	})
}

func (op *ServiceEndpointGatewayOp) Delete(ctx context.Context, id string) error {
	return errorFromDecodedResponseE("ServiceEndpointGateway.Delete", func() (*v1.ModelsApplianceReadServiceEndpointGatewayResponseBody, error) {
		return op.client.DeleteServiceEndpointGateway(ctx, v1.DeleteServiceEndpointGatewayParams{ApplianceID: id})
	})
}

func (op *ServiceEndpointGatewayOp) ReadInterface(ctx context.Context, applianceID string, interfaceID string) (*v1.ModelsApplianceReadServiceEndpointGatewayInterfaceResponseBody, error) {
	return errorFromDecodedResponse("ServiceEndpointGateway.ReadInterface", func() (*v1.ModelsApplianceReadServiceEndpointGatewayInterfaceResponseBody, error) {
		return op.client.ReadServiceEndpointGatewayInterface(ctx, v1.ReadServiceEndpointGatewayInterfaceParams{ApplianceID: applianceID, InterfaceID: interfaceID})
	})
}

func (op *ServiceEndpointGatewayOp) ReadPowerStatus(ctx context.Context, id string) (*v1.ModelsPowerReadServiceEndpointGatewayPowerStatusResponseBody, error) {
	return errorFromDecodedResponse("ServiceEndpointGateway.ReadPowerStatus", func() (*v1.ModelsPowerReadServiceEndpointGatewayPowerStatusResponseBody, error) {
		return op.client.ReadServiceEndpointGatewayPowerStatus(ctx, v1.ReadServiceEndpointGatewayPowerStatusParams{ApplianceID: id})
	})
}

func (op *ServiceEndpointGatewayOp) PowerOn(ctx context.Context, id string) (*v1.ModelsPowerUpdateServiceEndpointGatewayPowerStatusResponseBody, error) {
	var result *v1.ModelsPowerUpdateServiceEndpointGatewayPowerStatusResponseBody
	err := retryWithTimeout(ctx, op.powerRetryCfg.Interval, op.powerRetryCfg.Timeout, func() (ok bool, err error) {
		result, err = errorFromDecodedResponse("ServiceEndpointGateway.PowerOn", func() (*v1.ModelsPowerUpdateServiceEndpointGatewayPowerStatusResponseBody, error) {
			return op.client.PowerOnServiceEndpointGateway(ctx, v1.PowerOnServiceEndpointGatewayParams{ApplianceID: id})
		})
		return isStillCreatingError(err), err
	})
	return result, err
}

func (op *ServiceEndpointGatewayOp) Shutdown(ctx context.Context, id string) error {
	request := v1.NewOptModelsPowerApplianceDeletePowerStatusRequest(v1.ModelsPowerApplianceDeletePowerStatusRequest{Force: true})
	return retryWithTimeout(ctx, op.powerRetryCfg.Interval, op.powerRetryCfg.Timeout, func() (ok bool, err error) {
		err = errorFromDecodedResponseE("ServiceEndpointGateway.Shutdown", func() (*v1.ModelsPowerUpdateServiceEndpointGatewayPowerStatusResponseBody, error) {
			return op.client.PowerOffServiceEndpointGateway(ctx, request, v1.PowerOffServiceEndpointGatewayParams{ApplianceID: id})
		})
		return isStillCreatingError(err), err
	})
}

func (op *ServiceEndpointGatewayOp) Reset(ctx context.Context, id string) error {
	return retryWithTimeout(ctx, op.powerRetryCfg.Interval, op.powerRetryCfg.Timeout, func() (ok bool, err error) {
		err = errorFromDecodedResponseE("ServiceEndpointGateway.Reset", func() (*v1.ModelsPowerUpdateServiceEndpointGatewayPowerStatusResponseBody, error) {
			return op.client.ResetServiceEndpointGateway(ctx, v1.ResetServiceEndpointGatewayParams{ApplianceID: id})
		})
		return isStillCreatingError(err), err
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

func isStillCreatingError(err error) bool {
	e, ok := errors.AsType[*v1.ModelsCommonDefaultErrorResponseBodyStatusCode](err)
	if !ok {
		return false
	}

	if e.GetStatusCode() != http.StatusConflict {
		return false
	}

	c, ok := new(e.GetResponse()).GetErrorCode().Get()
	if !ok {
		return false
	}

	if c != "still_creating" {
		return false
	}

	return true
}

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

package service_endpoint_gateway

import (
	"context"
	"errors"
	"net/http"
	"time"

	v1 "github.com/sacloud/service-endpoint-gateway-api-go/apis/v1"
)

const (
	defRetryInterval = 15 * time.Second
	defRetryTimeout  = 20 * time.Minute
)

type PowerAPI interface {
	Read(ctx context.Context, id string) (*v1.ModelsPowerApplianceGetPowerStatusResponseBody, error)
	Update(ctx context.Context, id string) (*v1.ModelsPowerApplianceUpdatePowerStatusResponseBody, error)
	Delete(ctx context.Context, id string) error
	Reset(ctx context.Context, id string) error
}

type PowerRetryConfig struct {
	Interval time.Duration
	Timeout  time.Duration
}

func defaultPowerRetryConfig() PowerRetryConfig {
	return PowerRetryConfig{
		Interval: defRetryInterval,
		Timeout:  defRetryTimeout,
	}
}

var _ PowerAPI = (*PowerOp)(nil)

type PowerOp struct {
	client   *v1.Client
	retryCfg PowerRetryConfig
}

func NewPowerOp(client *v1.Client) PowerAPI {
	cfg := defaultPowerRetryConfig()
	return &PowerOp{client: client, retryCfg: cfg}
}

func NewPowerOpWithRetryConfig(client *v1.Client, retryCfg PowerRetryConfig) PowerAPI {
	return &PowerOp{client: client, retryCfg: retryCfg}
}

func (op *PowerOp) Read(ctx context.Context, id string) (*v1.ModelsPowerApplianceGetPowerStatusResponseBody, error) {
	const methodName = "Power.Read"

	res, err := op.client.SegStatusGetPowerStatus(ctx, v1.SegStatusGetPowerStatusParams{ApplianceID: id})
	if err != nil {
		var e *v1.ModelsCommonDefaultErrorResponseBodyStatusCode
		if errors.As(err, &e) {
			return nil, NewAPIError(methodName, e.StatusCode, err)
		}
		return nil, err
	}
	return res, nil
}

func (op *PowerOp) Update(ctx context.Context, id string) (*v1.ModelsPowerApplianceUpdatePowerStatusResponseBody, error) {
	const methodName = "Power.Update"

	var result *v1.ModelsPowerApplianceUpdatePowerStatusResponseBody
	err := retryWithTimeout(ctx, op.retryCfg.Interval, op.retryCfg.Timeout, func() (bool, error) {
		res, err := op.client.SegStatusUpdatePowerStatus(ctx, v1.SegStatusUpdatePowerStatusParams{ApplianceID: id})
		if err != nil {
			var e *v1.ModelsCommonDefaultErrorResponseBodyStatusCode
			if errors.As(err, &e) && e.StatusCode == http.StatusConflict && e.Response.ErrorCode.Value == "still_creating" {
				return true, err
			}
			if errors.As(err, &e) {
				return false, NewAPIError(methodName, e.StatusCode, err)
			}
			return false, err
		}
		result = res
		return false, nil
	})
	return result, err
}

func (op *PowerOp) Delete(ctx context.Context, id string) error {
	const methodName = "Power.Delete"

	request := v1.NewOptModelsPowerApplianceDeletePowerStatusRequest(v1.ModelsPowerApplianceDeletePowerStatusRequest{Force: true})
	return retryWithTimeout(ctx, op.retryCfg.Interval, op.retryCfg.Timeout, func() (bool, error) {
		_, err := op.client.SegStatusDeletePowerStatus(ctx, request, v1.SegStatusDeletePowerStatusParams{ApplianceID: id})
		if err != nil {
			var e *v1.ModelsCommonDefaultErrorResponseBodyStatusCode
			if errors.As(err, &e) && e.StatusCode == http.StatusConflict && e.Response.ErrorCode.Value == "still_creating" {
				return true, err
			}
			if errors.As(err, &e) {
				return false, NewAPIError(methodName, e.StatusCode, err)
			}
			return false, err
		}
		return false, nil
	})
}

func (op *PowerOp) Reset(ctx context.Context, id string) error {
	const methodName = "Power.Reset"

	return retryWithTimeout(ctx, op.retryCfg.Interval, op.retryCfg.Timeout, func() (bool, error) {
		_, err := op.client.SegStatusResetPowerStatus(ctx, v1.SegStatusResetPowerStatusParams{ApplianceID: id})
		if err != nil {
			var e *v1.ModelsCommonDefaultErrorResponseBodyStatusCode
			if errors.As(err, &e) && e.StatusCode == http.StatusConflict && e.Response.ErrorCode.Value == "still_creating" {
				return true, err
			}
			if errors.As(err, &e) {
				return false, NewAPIError(methodName, e.StatusCode, err)
			}
			return false, err
		}
		return false, nil
	})
}

// retryWithTimeout is a helper to retry a function with timeout and interval.
func retryWithTimeout(ctx context.Context, interval, timeout time.Duration, retryable func() (bool, error)) error {
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

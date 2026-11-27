// Copyright 2025- The sacloud/cloudhsm-api-go Authors
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

package cloudhsm

import (
	"context"
	"errors"
	"time"

	v1 "github.com/sacloud/sacloud-sdk-go/api/cloudhsm/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/packages/into"
)

type ClientAPI interface {
	List(ctx context.Context, count, from *int) ([]v1.CloudHSMClient, error)
	Create(ctx context.Context, request CloudHSMClientCreateParams) (*v1.CloudHSMClient, error)
	Read(ctx context.Context, id string) (*v1.CloudHSMClient, error)
	Update(ctx context.Context, id string, name string) (*v1.CloudHSMClient, error)
	Delete(ctx context.Context, id string) error
}

var _ ClientAPI = (*ClientOp)(nil)

type ClientOp struct {
	client *v1.Client
	hsm    *v1.CloudHSM
}

func NewClientOp(client *v1.Client, hsm *v1.CloudHSM) (ClientAPI, error) {
	if hsm.GetAvailability() == v1.CloudHSMAvailabilityAvailable {
		return &ClientOp{
			client: client,
			hsm:    hsm,
		}, nil
	}
	return nil, errors.New("CloudHSM unavailable")
}

func (op *ClientOp) List(ctx context.Context, count, from *int) ([]v1.CloudHSMClient, error) {
	res, err := errorFromDecodedResponse("Client.List", func() (*v1.PaginatedCloudHSMClientList, error) {
		return op.client.ListCloudHSMClients(ctx, v1.ListCloudHSMClientsParams{
			CloudhsmResourceID: op.hsm.GetID(),
			Count:              into.Opt[v1.OptInt](count),
			From:               into.Opt[v1.OptInt](from),
		})
	})
	if err != nil {
		return nil, err
	}

	return res.GetClients(), nil
}

type CloudHSMClientCreateParams struct {
	Name        string
	Certificate string
}

var epoc = v1.DateTime(time.Unix(0, 0).Format(time.RFC3339))

func (op *ClientOp) Create(ctx context.Context, p CloudHSMClientCreateParams) (*v1.CloudHSMClient, error) {
	resp, err := errorFromDecodedResponse("Client.Create", func() (*v1.WrappedCreateCloudHSMClient, error) {
		return op.client.CreateCloudHSMClient(ctx, &v1.WrappedCreateCloudHSMClientRequest{
			Client: v1.CreateCloudHSMClientRequest{
				Name:        p.Name,
				Certificate: p.Certificate,
			},
		}, v1.CreateCloudHSMClientParams{
			CloudhsmResourceID: op.hsm.GetID(),
		})
	})
	if err != nil {
		return nil, err
	}

	c := resp.GetClient()
	// Convert v1.CreateCloudHSMClient to v1.CloudHSMClient
	client := v1.CloudHSMClient{
		ID:           c.GetID().Or(""),
		CreatedAt:    c.GetCreatedAt().Or(epoc),
		ModifiedAt:   c.GetModifiedAt().Or(epoc),
		Availability: v1.CloudHSMClientAvailability(c.GetAvailability().Or(v1.CreateCloudHSMClientAvailabilityPrecreate)),
		Name:         c.Name,
		Certificate:  c.Certificate,
	}
	return &client, nil
}

func (op *ClientOp) Read(ctx context.Context, id string) (*v1.CloudHSMClient, error) {
	resp, err := errorFromDecodedResponse("Client.Read", func() (*v1.WrappedCloudHSMClient, error) {
		return op.client.ReadCloudHSMClient(ctx, v1.ReadCloudHSMClientParams{
			CloudhsmResourceID: op.hsm.GetID(),
			ID:                 id,
		})
	})
	if err != nil {
		return nil, err
	}

	return new(resp.GetClient()), nil
}

func (op *ClientOp) Update(ctx context.Context, id string, p string) (*v1.CloudHSMClient, error) {
	resp, err := errorFromDecodedResponse("Client.Update", func() (*v1.WrappedCloudHSMClient, error) {
		return op.client.UpdateCloudHSMClient(ctx, &v1.WrappedCloudHSMClientRequest{
			Client: v1.CloudHSMClientRequest{
				Name: p,
			},
		}, v1.UpdateCloudHSMClientParams{
			CloudhsmResourceID: op.hsm.GetID(),
			ID:                 id,
		})
	})
	if err != nil {
		return nil, err
	}

	return new(resp.GetClient()), nil
}

func (op *ClientOp) Delete(ctx context.Context, id string) error {
	return errorFromDecodedResponseE("Client.Delete", func() error {
		return op.client.DeleteCloudHSMClient(ctx, v1.DeleteCloudHSMClientParams{
			CloudhsmResourceID: op.hsm.GetID(),
			ID:                 id,
		})
	})
}

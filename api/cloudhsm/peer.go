// Copyright 2025- The sacloud/cloudhsm-api-go authors
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

	v1 "github.com/sacloud/sacloud-sdk-go/api/cloudhsm/apis/v1"
)

type PeerAPI interface {
	List(ctx context.Context) ([]v1.CloudHSMPeer, error)
	Create(ctx context.Context, request CloudHSMPeerCreateParams) error
	Delete(ctx context.Context, id string) error
}

var _ PeerAPI = (*PeerOp)(nil)

type PeerOp struct {
	client *v1.Client
	hsm    *v1.CloudHSM
}

func NewPeerOp(client *v1.Client, hsm *v1.CloudHSM) (PeerAPI, error) {
	// The HSM partition has to be "available" before doing anything with its peers.
	if hsm.GetAvailability() == v1.CloudHSMAvailabilityAvailable {
		return &PeerOp{
			client: client,
			hsm:    hsm,
		}, nil
	}

	return nil, errors.New("CloudHSM unavailable")
}

func (op *PeerOp) List(ctx context.Context) ([]v1.CloudHSMPeer, error) {
	resp, err := errorFromDecodedResponse("Peer.List", func() (*v1.WrappedPeerList, error) {
		return op.client.ListCloudHSMPeers(
			ctx,
			v1.ListCloudHSMPeersParams{
				ResourceID: op.hsm.GetID(),
			},
		)
	})
	if err != nil {
		return nil, err
	}

	return resp.GetPeers(), nil
}

type CloudHSMPeerCreateParams struct {
	RouterID  string
	SecretKey string
}

func (op *PeerOp) Create(ctx context.Context, p CloudHSMPeerCreateParams) error {
	return errorFromDecodedResponseE("Peer.Create", func() error {
		return op.client.CreateCloudHSMPeer(
			ctx,
			&v1.WrappedPeerRequest{
				Peer: v1.PeerRequest{
					ID:        p.RouterID,
					SecretKey: p.SecretKey,
				},
			},
			v1.CreateCloudHSMPeerParams{
				ResourceID: op.hsm.GetID(),
			},
		)
	})
}

func (op *PeerOp) Delete(ctx context.Context, id string) error {
	return errorFromDecodedResponseE("Peer.Delete", func() error {
		return op.client.DeleteCloudHSMPeer(
			ctx,
			v1.DeleteCloudHSMPeerParams{
				ResourceID: op.hsm.GetID(),
				PeerID:     id,
			},
		)
	})
}

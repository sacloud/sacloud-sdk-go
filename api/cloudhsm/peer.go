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
	"net/http"

	"github.com/go-faster/errors"
	ogen "github.com/ogen-go/ogen/validate"
	v1 "github.com/sacloud/cloudhsm-api-go/apis/v1"
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
	if hsm.GetAvailability() == v1.AvailabilityEnumAvailable {
		return &PeerOp{
			client: client,
			hsm:    hsm,
		}, nil
	}

	return nil, errors.New("CloudHSM unavailable")
}

func (op *PeerOp) List(ctx context.Context) ([]v1.CloudHSMPeer, error) {
	resp, err := op.client.CloudhsmCloudhsmsPeersRetrieve(
		ctx,
		v1.CloudhsmCloudhsmsPeersRetrieveParams{
			ResourceID: op.hsm.GetID(),
		},
	)

	if err == nil {
		return resp.GetPeers(), nil
	} else if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); !ok {
		return nil, NewAPIError("Peer.List", 0, err)
	} else if e.StatusCode == http.StatusNotFound {
		return nil, NewAPIError("Peer.List", e.StatusCode, errors.Wrap(err, "not found"))
	} else {
		return nil, NewAPIError("Peer.List", e.StatusCode, errors.Wrap(err, "internal server error"))
	}
}

type CloudHSMPeerCreateParams struct {
	RouterID  string
	SecretKey string
}

func (op *PeerOp) Create(ctx context.Context, p CloudHSMPeerCreateParams) error {
	err := op.client.CloudhsmCloudhsmsPeersCreate(
		ctx,
		&v1.WrappedCreateCloudHSMPeer{
			Peer: v1.CreateCloudHSMPeer{
				ID:        p.RouterID,
				SecretKey: p.SecretKey,
			},
		},
		v1.CloudhsmCloudhsmsPeersCreateParams{
			ResourceID: op.hsm.GetID(),
		},
	)

	if err == nil {
		return nil
	} else if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); !ok {
		return NewAPIError("Peer.Create", 0, err)
	} else if e.StatusCode == http.StatusUnprocessableEntity {
		return NewAPIError("Peer.Create", e.StatusCode, errors.Wrap(err, "invalid parameter"))
	} else {
		return NewAPIError("Peer.Create", e.StatusCode, errors.Wrap(err, "internal server error"))
	}
}

func (op *PeerOp) Delete(ctx context.Context, id string) error {
	err := op.client.CloudhsmCloudhsmsPeersDestroy(
		ctx,
		v1.CloudhsmCloudhsmsPeersDestroyParams{
			ResourceID: op.hsm.GetID(),
			PeerID:     id,
		},
	)

	if err == nil {
		return nil
	} else if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); !ok {
		return NewAPIError("Peer.Delete", 0, err)
	} else if e.StatusCode == http.StatusNotFound {
		return NewAPIError("Peer.Delete", e.StatusCode, errors.Wrap(err, "not found"))
	} else {
		return NewAPIError("Peer.Delete", e.StatusCode, errors.Wrap(err, "internal server error"))
	}
}

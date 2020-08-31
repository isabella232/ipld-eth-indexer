// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package eth

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/statediff"
	"github.com/sirupsen/logrus"
)

const (
	PayloadChanBufferSize = 20000 // the max eth sub buffer size
)

// StreamClient is an interface for subscribing and streaming from geth
type StreamClient interface {
	Subscribe(ctx context.Context, namespace string, payloadChan interface{}, args ...interface{}) (*rpc.ClientSubscription, error)
}

// Streamer interface for substituting mocks in tests
type Streamer interface {
	Stream(payloadChan chan statediff.Payload) (*rpc.ClientSubscription, error)
}

// PayloadStreamer satisfies the PayloadStreamer interface for ethereum
type PayloadStreamer struct {
	Client StreamClient
	params statediff.Params
}

// NewPayloadStreamer creates a pointer to a new PayloadStreamer which satisfies the PayloadStreamer interface for ethereum
func NewPayloadStreamer(client StreamClient) *PayloadStreamer {
	return &PayloadStreamer{
		Client: client,
		params: statediff.Params{
			IncludeBlock:             true,
			IncludeTD:                true,
			IncludeReceipts:          true,
			IntermediateStorageNodes: true,
			IntermediateStateNodes:   true,
		},
	}
}

// Stream is the main loop for subscribing to data from the Geth state diff process
// Satisfies the shared.PayloadStreamer interface
func (ps *PayloadStreamer) Stream(payloadChan chan statediff.Payload) (*rpc.ClientSubscription, error) {
	logrus.Debug("streaming diffs from geth")
	return ps.Client.Subscribe(context.Background(), "statediff", payloadChan, "stream", ps.params)
}

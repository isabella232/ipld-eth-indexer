// Copyright 2019 Vulcanize
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

package eth_test

import (
	"github.com/ethereum/go-ethereum/statediff"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/ipld-eth-indexer/pkg/eth"
	"github.com/vulcanize/ipld-eth-indexer/pkg/eth/mocks"
)

var _ = Describe("StateDiff Streamer", func() {
	It("subscribes to the geth statediff service", func() {
		client := &mocks.StreamClient{}
		streamer := eth.NewPayloadStreamer(client)
		payloadChan := make(chan statediff.Payload)
		_, err := streamer.Stream(payloadChan)
		Expect(err).NotTo(HaveOccurred())
	})
})

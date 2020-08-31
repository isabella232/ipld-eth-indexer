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

package eth_test

import (
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/ipld-eth-indexer/pkg/eth"
	"github.com/vulcanize/ipld-eth-indexer/pkg/eth/mocks"
)

var _ = Describe("Converter", func() {
	Describe("Convert", func() {
		It("Converts mock statediff.Payloads into the expected IPLDPayloads", func() {
			converter := eth.NewPayloadConverter(params.MainnetChainConfig)
			payload, err := converter.Convert(mocks.MockStateDiffPayload)
			Expect(err).ToNot(HaveOccurred())
			Expect(payload.Block.Number().String()).To(Equal(mocks.BlockNumber.String()))
			Expect(payload.Block.Hash().String()).To(Equal(mocks.MockBlock.Hash().String()))
			Expect(payload.StateNodes).To(Equal(mocks.MockStateNodes))
			Expect(payload.StorageNodes).To(Equal(mocks.MockStorageNodes))
			Expect(payload.TotalDifficulty.Int64()).To(Equal(mocks.MockStateDiffPayload.TotalDifficulty.Int64()))
			gotBody, err := rlp.EncodeToBytes(payload.Block.Body())
			Expect(err).ToNot(HaveOccurred())
			expectedBody, err := rlp.EncodeToBytes(mocks.MockBlock.Body())
			Expect(err).ToNot(HaveOccurred())
			Expect(gotBody).To(Equal(expectedBody))
			gotHeader, err := rlp.EncodeToBytes(payload.Block.Header())
			Expect(err).ToNot(HaveOccurred())
			Expect(gotHeader).To(Equal(mocks.MockHeaderRlp))
			Expect(payload.TxMetaData).To(Equal(mocks.MockTrxMeta))
			Expect(payload.ReceiptMetaData).To(Equal(mocks.MockRctMeta))
		})
	})
})

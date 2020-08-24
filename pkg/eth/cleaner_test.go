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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/ipfs-blockchain-watcher/pkg/eth"
	"github.com/vulcanize/ipfs-blockchain-watcher/pkg/postgres"
	"github.com/vulcanize/ipfs-blockchain-watcher/pkg/shared"
)

var (
	// Block 0
	// header variables
	blockHash1      = crypto.Keccak256Hash([]byte{00, 02})
	blocKNumber1    = big.NewInt(0)
	headerCID1      = shared.TestCID([]byte("mockHeader1CID"))
	headerMhKey1    = shared.MultihashKeyFromCID(headerCID1)
	parentHash      = crypto.Keccak256Hash([]byte{00, 01})
	totalDifficulty = "50000000000000000000"
	reward          = "5000000000000000000"
	headerModel     = eth.HeaderModel{
		BlockHash:       blockHash1.String(),
		BlockNumber:     blocKNumber1.String(),
		CID:             headerCID1.String(),
		MhKey:           headerMhKey1,
		ParentHash:      parentHash.String(),
		TotalDifficulty: totalDifficulty,
		Reward:          reward,
	}

	// tx variables
	tx1CID    = shared.TestCID([]byte("mockTx1CID"))
	tx1MhKey  = shared.MultihashKeyFromCID(tx1CID)
	tx2CID    = shared.TestCID([]byte("mockTx2CID"))
	tx2MhKey  = shared.MultihashKeyFromCID(tx2CID)
	tx1Hash   = crypto.Keccak256Hash([]byte{01, 01})
	tx2Hash   = crypto.Keccak256Hash([]byte{01, 02})
	txSrc     = common.HexToAddress("0x010a")
	txDst     = common.HexToAddress("0x020a")
	txModels1 = []eth.TxModel{
		{
			CID:    tx1CID.String(),
			MhKey:  tx1MhKey,
			TxHash: tx1Hash.String(),
			Index:  0,
		},
		{
			CID:    tx2CID.String(),
			MhKey:  tx2MhKey,
			TxHash: tx2Hash.String(),
			Index:  1,
		},
	}

	// uncle variables
	uncleCID        = shared.TestCID([]byte("mockUncle1CID"))
	uncleMhKey      = shared.MultihashKeyFromCID(uncleCID)
	uncleHash       = crypto.Keccak256Hash([]byte{02, 02})
	uncleParentHash = crypto.Keccak256Hash([]byte{02, 01})
	uncleReward     = "1000000000000000000"
	uncleModels1    = []eth.UncleModel{
		{
			CID:        uncleCID.String(),
			MhKey:      uncleMhKey,
			Reward:     uncleReward,
			BlockHash:  uncleHash.String(),
			ParentHash: uncleParentHash.String(),
		},
	}

	// receipt variables
	rct1CID        = shared.TestCID([]byte("mockRct1CID"))
	rct1MhKey      = shared.MultihashKeyFromCID(rct1CID)
	rct2CID        = shared.TestCID([]byte("mockRct2CID"))
	rct2MhKey      = shared.MultihashKeyFromCID(rct2CID)
	rct1Contract   = common.Address{}
	rct2Contract   = common.HexToAddress("0x010c")
	receiptModels1 = map[common.Hash]eth.ReceiptModel{
		tx1Hash: {
			CID:          rct1CID.String(),
			MhKey:        rct1MhKey,
			ContractHash: crypto.Keccak256Hash(rct1Contract.Bytes()).String(),
		},
		tx2Hash: {
			CID:          rct2CID.String(),
			MhKey:        rct2MhKey,
			ContractHash: crypto.Keccak256Hash(rct2Contract.Bytes()).String(),
		},
	}

	// state variables
	state1CID1   = shared.TestCID([]byte("mockState1CID1"))
	state1MhKey1 = shared.MultihashKeyFromCID(state1CID1)
	state1Path   = []byte{'\x01'}
	state1Key    = crypto.Keccak256Hash(txSrc.Bytes())
	state2CID1   = shared.TestCID([]byte("mockState2CID1"))
	state2MhKey1 = shared.MultihashKeyFromCID(state2CID1)
	state2Path   = []byte{'\x02'}
	state2Key    = crypto.Keccak256Hash(txDst.Bytes())
	stateModels1 = []eth.StateNodeModel{
		{
			CID:      state1CID1.String(),
			MhKey:    state1MhKey1,
			Path:     state1Path,
			NodeType: 2,
			StateKey: state1Key.String(),
		},
		{
			CID:      state2CID1.String(),
			MhKey:    state2MhKey1,
			Path:     state2Path,
			NodeType: 2,
			StateKey: state2Key.String(),
		},
	}

	// storage variables
	storageCID     = shared.TestCID([]byte("mockStorageCID1"))
	storageMhKey   = shared.MultihashKeyFromCID(storageCID)
	storagePath    = []byte{'\x01'}
	storageKey     = crypto.Keccak256Hash(common.Hex2Bytes("0x0000000000000000000000000000000000000000000000000000000000000000"))
	storageModels1 = map[string][]eth.StorageNodeModel{
		common.Bytes2Hex(state1Path): {
			{
				CID:        storageCID.String(),
				MhKey:      storageMhKey,
				StorageKey: storageKey.String(),
				Path:       storagePath,
				NodeType:   2,
			},
		},
	}
	mockCIDPayload1 = eth.CIDPayload{
		HeaderCID:       headerModel,
		UncleCIDs:       uncleModels1,
		TransactionCIDs: txModels1,
		ReceiptCIDs:     receiptModels1,
		StateNodeCIDs:   stateModels1,
		StorageNodeCIDs: storageModels1,
	}

	// Block 1
	// header variables
	blockHash2   = crypto.Keccak256Hash([]byte{00, 03})
	blocKNumber2 = big.NewInt(1)
	headerCID2   = shared.TestCID([]byte("mockHeaderCID2"))
	headerMhKey2 = shared.MultihashKeyFromCID(headerCID2)
	headerModel2 = eth.HeaderModel{
		BlockHash:       blockHash2.String(),
		BlockNumber:     blocKNumber2.String(),
		CID:             headerCID2.String(),
		MhKey:           headerMhKey2,
		ParentHash:      blockHash1.String(),
		TotalDifficulty: totalDifficulty,
		Reward:          reward,
	}
	// tx variables
	tx3CID    = shared.TestCID([]byte("mockTx3CID"))
	tx3MhKey  = shared.MultihashKeyFromCID(tx3CID)
	tx3Hash   = crypto.Keccak256Hash([]byte{01, 03})
	txModels2 = []eth.TxModel{
		{
			CID:    tx3CID.String(),
			MhKey:  tx3MhKey,
			TxHash: tx3Hash.String(),
			Index:  0,
		},
	}
	// receipt variables
	rct3CID        = shared.TestCID([]byte("mockRct3CID"))
	rct3MhKey      = shared.MultihashKeyFromCID(rct3CID)
	receiptModels2 = map[common.Hash]eth.ReceiptModel{
		tx3Hash: {
			CID:          rct3CID.String(),
			MhKey:        rct3MhKey,
			ContractHash: crypto.Keccak256Hash(rct1Contract.Bytes()).String(),
		},
	}

	// state variables
	state1CID2   = shared.TestCID([]byte("mockState1CID2"))
	state1MhKey2 = shared.MultihashKeyFromCID(state1CID2)
	stateModels2 = []eth.StateNodeModel{
		{
			CID:      state1CID2.String(),
			MhKey:    state1MhKey2,
			Path:     state1Path,
			NodeType: 2,
			StateKey: state1Key.String(),
		},
	}
	mockCIDPayload2 = eth.CIDPayload{
		HeaderCID:       headerModel2,
		TransactionCIDs: txModels2,
		ReceiptCIDs:     receiptModels2,
		StateNodeCIDs:   stateModels2,
	}
	rngs   = [][2]uint64{{0, 1}}
	mhKeys = []string{
		headerMhKey1,
		headerMhKey2,
		uncleMhKey,
		tx1MhKey,
		tx2MhKey,
		tx3MhKey,
		rct1MhKey,
		rct2MhKey,
		rct3MhKey,
		state1MhKey1,
		state2MhKey1,
		state1MhKey2,
		storageMhKey,
	}
	mockData = []byte{'\x01'}
)

var _ = Describe("Cleaner", func() {
	var (
		db      *postgres.DB
		repo    *eth.CIDIndexer
		cleaner *eth.DBCleaner
	)
	BeforeEach(func() {
		var err error
		db, err = shared.SetupDB()
		Expect(err).ToNot(HaveOccurred())
		repo = eth.NewCIDIndexer(db)
		cleaner = eth.NewCleaner(db)
	})
	Describe("Clean", func() {
		BeforeEach(func() {
			for _, key := range mhKeys {
				_, err := db.Exec(`INSERT INTO public.blocks (key, data) VALUES ($1, $2)`, key, mockData)
				Expect(err).ToNot(HaveOccurred())
			}

			err := repo.Index(mockCIDPayload1)
			Expect(err).ToNot(HaveOccurred())
			err = repo.Index(mockCIDPayload2)
			Expect(err).ToNot(HaveOccurred())

			tx, err := db.Beginx()
			Expect(err).ToNot(HaveOccurred())

			var startingIPFSBlocksCount int
			pgStr := `SELECT COUNT(*) FROM public.blocks`
			err = tx.Get(&startingIPFSBlocksCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var startingStorageCount int
			pgStr = `SELECT COUNT(*) FROM eth.storage_cids`
			err = tx.Get(&startingStorageCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var startingStateCount int
			pgStr = `SELECT COUNT(*) FROM eth.state_cids`
			err = tx.Get(&startingStateCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var startingReceiptCount int
			pgStr = `SELECT COUNT(*) FROM eth.receipt_cids`
			err = tx.Get(&startingReceiptCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var startingTxCount int
			pgStr = `SELECT COUNT(*) FROM eth.transaction_cids`
			err = tx.Get(&startingTxCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var startingUncleCount int
			pgStr = `SELECT COUNT(*) FROM eth.uncle_cids`
			err = tx.Get(&startingUncleCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var startingHeaderCount int
			pgStr = `SELECT COUNT(*) FROM eth.header_cids`
			err = tx.Get(&startingHeaderCount, pgStr)
			Expect(err).ToNot(HaveOccurred())

			err = tx.Commit()
			Expect(err).ToNot(HaveOccurred())

			Expect(startingIPFSBlocksCount).To(Equal(13))
			Expect(startingStorageCount).To(Equal(1))
			Expect(startingStateCount).To(Equal(3))
			Expect(startingReceiptCount).To(Equal(3))
			Expect(startingTxCount).To(Equal(3))
			Expect(startingUncleCount).To(Equal(1))
			Expect(startingHeaderCount).To(Equal(2))
		})
		AfterEach(func() {
			eth.TearDownDB(db)
		})
		It("Cleans everything", func() {
			err := cleaner.Clean(rngs, shared.Full)
			Expect(err).ToNot(HaveOccurred())

			tx, err := db.Beginx()
			Expect(err).ToNot(HaveOccurred())

			pgStr := `SELECT COUNT(*) FROM eth.header_cids`
			var headerCount int
			err = tx.Get(&headerCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var uncleCount int
			pgStr = `SELECT COUNT(*) FROM eth.uncle_cids`
			err = tx.Get(&uncleCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var txCount int
			pgStr = `SELECT COUNT(*) FROM eth.transaction_cids`
			err = tx.Get(&txCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var rctCount int
			pgStr = `SELECT COUNT(*) FROM eth.receipt_cids`
			err = tx.Get(&rctCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var stateCount int
			pgStr = `SELECT COUNT(*) FROM eth.state_cids`
			err = tx.Get(&stateCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var storageCount int
			pgStr = `SELECT COUNT(*) FROM eth.storage_cids`
			err = tx.Get(&storageCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var blocksCount int
			pgStr = `SELECT COUNT(*) FROM public.blocks`
			err = tx.Get(&blocksCount, pgStr)
			Expect(err).ToNot(HaveOccurred())

			err = tx.Commit()
			Expect(err).ToNot(HaveOccurred())

			Expect(headerCount).To(Equal(0))
			Expect(uncleCount).To(Equal(0))
			Expect(txCount).To(Equal(0))
			Expect(rctCount).To(Equal(0))
			Expect(stateCount).To(Equal(0))
			Expect(storageCount).To(Equal(0))
			Expect(blocksCount).To(Equal(0))
		})
		It("Cleans headers and all linked data (same as full)", func() {
			err := cleaner.Clean(rngs, shared.Headers)
			Expect(err).ToNot(HaveOccurred())

			tx, err := db.Beginx()
			Expect(err).ToNot(HaveOccurred())

			var headerCount int
			pgStr := `SELECT COUNT(*) FROM eth.header_cids`
			err = tx.Get(&headerCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var uncleCount int
			pgStr = `SELECT COUNT(*) FROM eth.uncle_cids`
			err = tx.Get(&uncleCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var txCount int
			pgStr = `SELECT COUNT(*) FROM eth.transaction_cids`
			err = tx.Get(&txCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var rctCount int
			pgStr = `SELECT COUNT(*) FROM eth.receipt_cids`
			err = tx.Get(&rctCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var stateCount int
			pgStr = `SELECT COUNT(*) FROM eth.state_cids`
			err = tx.Get(&stateCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var storageCount int
			pgStr = `SELECT COUNT(*) FROM eth.storage_cids`
			err = tx.Get(&storageCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var blocksCount int
			pgStr = `SELECT COUNT(*) FROM public.blocks`
			err = tx.Get(&blocksCount, pgStr)
			Expect(err).ToNot(HaveOccurred())

			err = tx.Commit()
			Expect(err).ToNot(HaveOccurred())

			Expect(headerCount).To(Equal(0))
			Expect(uncleCount).To(Equal(0))
			Expect(txCount).To(Equal(0))
			Expect(rctCount).To(Equal(0))
			Expect(stateCount).To(Equal(0))
			Expect(storageCount).To(Equal(0))
			Expect(blocksCount).To(Equal(0))
		})
		It("Cleans uncles", func() {
			err := cleaner.Clean(rngs, shared.Uncles)
			Expect(err).ToNot(HaveOccurred())

			tx, err := db.Beginx()
			Expect(err).ToNot(HaveOccurred())

			var headerCount int
			pgStr := `SELECT COUNT(*) FROM eth.header_cids`
			err = tx.Get(&headerCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var uncleCount int
			pgStr = `SELECT COUNT(*) FROM eth.uncle_cids`
			err = tx.Get(&uncleCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var txCount int
			pgStr = `SELECT COUNT(*) FROM eth.transaction_cids`
			err = tx.Get(&txCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var rctCount int
			pgStr = `SELECT COUNT(*) FROM eth.receipt_cids`
			err = tx.Get(&rctCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var stateCount int
			pgStr = `SELECT COUNT(*) FROM eth.state_cids`
			err = tx.Get(&stateCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var storageCount int
			pgStr = `SELECT COUNT(*) FROM eth.storage_cids`
			err = tx.Get(&storageCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var blocksCount int
			pgStr = `SELECT COUNT(*) FROM public.blocks`
			err = tx.Get(&blocksCount, pgStr)
			Expect(err).ToNot(HaveOccurred())

			err = tx.Commit()
			Expect(err).ToNot(HaveOccurred())

			Expect(headerCount).To(Equal(2))
			Expect(uncleCount).To(Equal(0))
			Expect(txCount).To(Equal(3))
			Expect(rctCount).To(Equal(3))
			Expect(stateCount).To(Equal(3))
			Expect(storageCount).To(Equal(1))
			Expect(blocksCount).To(Equal(12))
		})
		It("Cleans transactions and linked receipts", func() {
			err := cleaner.Clean(rngs, shared.Transactions)
			Expect(err).ToNot(HaveOccurred())

			tx, err := db.Beginx()
			Expect(err).ToNot(HaveOccurred())

			var headerCount int
			pgStr := `SELECT COUNT(*) FROM eth.header_cids`
			err = tx.Get(&headerCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var uncleCount int
			pgStr = `SELECT COUNT(*) FROM eth.uncle_cids`
			err = tx.Get(&uncleCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var txCount int
			pgStr = `SELECT COUNT(*) FROM eth.transaction_cids`
			err = tx.Get(&txCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var rctCount int
			pgStr = `SELECT COUNT(*) FROM eth.receipt_cids`
			err = tx.Get(&rctCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var stateCount int
			pgStr = `SELECT COUNT(*) FROM eth.state_cids`
			err = tx.Get(&stateCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var storageCount int
			pgStr = `SELECT COUNT(*) FROM eth.storage_cids`
			err = tx.Get(&storageCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var blocksCount int
			pgStr = `SELECT COUNT(*) FROM public.blocks`
			err = tx.Get(&blocksCount, pgStr)
			Expect(err).ToNot(HaveOccurred())

			err = tx.Commit()
			Expect(err).ToNot(HaveOccurred())

			Expect(headerCount).To(Equal(2))
			Expect(uncleCount).To(Equal(1))
			Expect(txCount).To(Equal(0))
			Expect(rctCount).To(Equal(0))
			Expect(stateCount).To(Equal(3))
			Expect(storageCount).To(Equal(1))
			Expect(blocksCount).To(Equal(7))
		})
		It("Cleans receipts", func() {
			err := cleaner.Clean(rngs, shared.Receipts)
			Expect(err).ToNot(HaveOccurred())

			tx, err := db.Beginx()
			Expect(err).ToNot(HaveOccurred())

			var headerCount int
			pgStr := `SELECT COUNT(*) FROM eth.header_cids`
			err = tx.Get(&headerCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var uncleCount int
			pgStr = `SELECT COUNT(*) FROM eth.uncle_cids`
			err = tx.Get(&uncleCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var txCount int
			pgStr = `SELECT COUNT(*) FROM eth.transaction_cids`
			err = tx.Get(&txCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var rctCount int
			pgStr = `SELECT COUNT(*) FROM eth.receipt_cids`
			err = tx.Get(&rctCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var stateCount int
			pgStr = `SELECT COUNT(*) FROM eth.state_cids`
			err = tx.Get(&stateCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var storageCount int
			pgStr = `SELECT COUNT(*) FROM eth.storage_cids`
			err = tx.Get(&storageCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var blocksCount int
			pgStr = `SELECT COUNT(*) FROM public.blocks`
			err = tx.Get(&blocksCount, pgStr)
			Expect(err).ToNot(HaveOccurred())

			err = tx.Commit()
			Expect(err).ToNot(HaveOccurred())

			Expect(headerCount).To(Equal(2))
			Expect(uncleCount).To(Equal(1))
			Expect(txCount).To(Equal(3))
			Expect(rctCount).To(Equal(0))
			Expect(stateCount).To(Equal(3))
			Expect(storageCount).To(Equal(1))
			Expect(blocksCount).To(Equal(10))
		})
		It("Cleans state and linked storage", func() {
			err := cleaner.Clean(rngs, shared.State)
			Expect(err).ToNot(HaveOccurred())

			tx, err := db.Beginx()
			Expect(err).ToNot(HaveOccurred())

			var headerCount int
			pgStr := `SELECT COUNT(*) FROM eth.header_cids`
			err = tx.Get(&headerCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var uncleCount int
			pgStr = `SELECT COUNT(*) FROM eth.uncle_cids`
			err = tx.Get(&uncleCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var txCount int
			pgStr = `SELECT COUNT(*) FROM eth.transaction_cids`
			err = tx.Get(&txCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var rctCount int
			pgStr = `SELECT COUNT(*) FROM eth.receipt_cids`
			err = tx.Get(&rctCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var stateCount int
			pgStr = `SELECT COUNT(*) FROM eth.state_cids`
			err = tx.Get(&stateCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var storageCount int
			pgStr = `SELECT COUNT(*) FROM eth.storage_cids`
			err = tx.Get(&storageCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var blocksCount int
			pgStr = `SELECT COUNT(*) FROM public.blocks`
			err = tx.Get(&blocksCount, pgStr)
			Expect(err).ToNot(HaveOccurred())

			err = tx.Commit()
			Expect(err).ToNot(HaveOccurred())

			Expect(headerCount).To(Equal(2))
			Expect(uncleCount).To(Equal(1))
			Expect(txCount).To(Equal(3))
			Expect(rctCount).To(Equal(3))
			Expect(stateCount).To(Equal(0))
			Expect(storageCount).To(Equal(0))
			Expect(blocksCount).To(Equal(9))
		})
		It("Cleans storage", func() {
			err := cleaner.Clean(rngs, shared.Storage)
			Expect(err).ToNot(HaveOccurred())

			tx, err := db.Beginx()
			Expect(err).ToNot(HaveOccurred())

			var headerCount int
			pgStr := `SELECT COUNT(*) FROM eth.header_cids`
			err = tx.Get(&headerCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var uncleCount int
			pgStr = `SELECT COUNT(*) FROM eth.uncle_cids`
			err = tx.Get(&uncleCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var txCount int
			pgStr = `SELECT COUNT(*) FROM eth.transaction_cids`
			err = tx.Get(&txCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var rctCount int
			pgStr = `SELECT COUNT(*) FROM eth.receipt_cids`
			err = tx.Get(&rctCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var stateCount int
			pgStr = `SELECT COUNT(*) FROM eth.state_cids`
			err = tx.Get(&stateCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var storageCount int
			pgStr = `SELECT COUNT(*) FROM eth.storage_cids`
			err = tx.Get(&storageCount, pgStr)
			Expect(err).ToNot(HaveOccurred())
			var blocksCount int
			pgStr = `SELECT COUNT(*) FROM public.blocks`
			err = tx.Get(&blocksCount, pgStr)
			Expect(err).ToNot(HaveOccurred())

			err = tx.Commit()
			Expect(err).ToNot(HaveOccurred())

			Expect(headerCount).To(Equal(2))
			Expect(uncleCount).To(Equal(1))
			Expect(txCount).To(Equal(3))
			Expect(rctCount).To(Equal(3))
			Expect(stateCount).To(Equal(3))
			Expect(storageCount).To(Equal(0))
			Expect(blocksCount).To(Equal(12))
		})
	})

	Describe("ResetValidation", func() {
		BeforeEach(func() {
			for _, key := range mhKeys {
				_, err := db.Exec(`INSERT INTO public.blocks (key, data) VALUES ($1, $2)`, key, mockData)
				Expect(err).ToNot(HaveOccurred())
			}

			err := repo.Index(mockCIDPayload1)
			Expect(err).ToNot(HaveOccurred())
			err = repo.Index(mockCIDPayload2)
			Expect(err).ToNot(HaveOccurred())

			var validationTimes []int
			pgStr := `SELECT times_validated FROM eth.header_cids`
			err = db.Select(&validationTimes, pgStr)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(validationTimes)).To(Equal(2))
			Expect(validationTimes[0]).To(Equal(1))
			Expect(validationTimes[1]).To(Equal(1))

			err = repo.Index(mockCIDPayload1)
			Expect(err).ToNot(HaveOccurred())

			validationTimes = []int{}
			pgStr = `SELECT times_validated FROM eth.header_cids ORDER BY block_number`
			err = db.Select(&validationTimes, pgStr)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(validationTimes)).To(Equal(2))
			Expect(validationTimes[0]).To(Equal(2))
			Expect(validationTimes[1]).To(Equal(1))
		})
		AfterEach(func() {
			eth.TearDownDB(db)
		})
		It("Resets the validation level", func() {
			err := cleaner.ResetValidation(rngs)
			Expect(err).ToNot(HaveOccurred())

			var validationTimes []int
			pgStr := `SELECT times_validated FROM eth.header_cids`
			err = db.Select(&validationTimes, pgStr)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(validationTimes)).To(Equal(2))
			Expect(validationTimes[0]).To(Equal(0))
			Expect(validationTimes[1]).To(Equal(0))

			err = repo.Index(mockCIDPayload2)
			Expect(err).ToNot(HaveOccurred())

			validationTimes = []int{}
			pgStr = `SELECT times_validated FROM eth.header_cids ORDER BY block_number`
			err = db.Select(&validationTimes, pgStr)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(validationTimes)).To(Equal(2))
			Expect(validationTimes[0]).To(Equal(0))
			Expect(validationTimes[1]).To(Equal(1))
		})
	})
})

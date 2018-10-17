// Copyright 2018 Vulcanize
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

package vat_move_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	vat_move_mocks "github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks/vat_move"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/vat_move"
)

var _ = Describe("Vat move transformer", func() {
	var fetcher mocks.MockLogFetcher
	var converter vat_move_mocks.MockVatMoveConverter
	var repository vat_move_mocks.MockVatMoveRepository
	var config = vat_move.VatMoveConfig
	var headerOne core.Header
	var headerTwo core.Header

	BeforeEach(func() {
		fetcher = mocks.MockLogFetcher{}
		converter = vat_move_mocks.MockVatMoveConverter{}
		repository = vat_move_mocks.MockVatMoveRepository{}
		headerOne = core.Header{Id: GinkgoRandomSeed(), BlockNumber: GinkgoRandomSeed()}
		headerTwo = core.Header{Id: GinkgoRandomSeed(), BlockNumber: GinkgoRandomSeed()}
	})

	It("gets missing headers for block numbers specified in config", func() {
		transformer := vat_move.VatMoveTransformer{
			Config:     config,
			Converter:  &converter,
			Fetcher:    &fetcher,
			Repository: &repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(vat_move.VatMoveConfig.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(vat_move.VatMoveConfig.EndingBlockNumber))
	})

	It("returns error if repository returns error for missing headers", func() {
		repository.SetMissingHeadersError(fakes.FakeError)
		transformer := vat_move.VatMoveTransformer{
			Config:     config,
			Converter:  &converter,
			Fetcher:    &fetcher,
			Repository: &repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches logs for missing headers", func() {
		repository.SetMissingHeaders([]core.Header{headerOne, headerTwo})
		transformer := vat_move.VatMoveTransformer{
			Config:     config,
			Fetcher:    &fetcher,
			Converter:  &vat_move_mocks.MockVatMoveConverter{},
			Repository: &repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{headerOne.BlockNumber, headerTwo.BlockNumber}))
		Expect(fetcher.FetchedContractAddresses).To(Equal([][]string{
			vat_move.VatMoveConfig.ContractAddresses,
			vat_move.VatMoveConfig.ContractAddresses,
		}))
		Expect(fetcher.FetchedTopics).To(Equal([][]common.Hash{{common.HexToHash(shared.VatMoveSignature)}}))
	})

	It("returns error if fetcher returns error", func() {
		fetcher.SetFetcherError(fakes.FakeError)
		repository.SetMissingHeaders([]core.Header{headerOne})
		transformer := vat_move.VatMoveTransformer{
			Config:     config,
			Fetcher:    &fetcher,
			Converter:  &converter,
			Repository: &repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts matching logs", func() {
		fetcher.SetFetchedLogs([]types.Log{test_data.EthVatMoveLog})
		repository.SetMissingHeaders([]core.Header{headerOne})
		transformer := vat_move.VatMoveTransformer{
			Config:     config,
			Fetcher:    &fetcher,
			Converter:  &converter,
			Repository: &repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedLogs).To(Equal([]types.Log{test_data.EthVatMoveLog}))
	})

	It("returns error if converter returns error", func() {
		converter.SetConverterError(fakes.FakeError)
		fetcher.SetFetchedLogs([]types.Log{test_data.EthVatMoveLog})
		repository.SetMissingHeaders([]core.Header{headerOne})
		transformer := vat_move.VatMoveTransformer{
			Config:     config,
			Fetcher:    &fetcher,
			Converter:  &converter,
			Repository: &repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("marks header as checked even if no logs were returned", func() {
		repository.SetMissingHeaders([]core.Header{headerOne, headerTwo})
		fetcher.SetFetchedLogs([]types.Log{})
		transformer := vat_move.VatMoveTransformer{
			Config:     config,
			Converter:  &converter,
			Fetcher:    &fetcher,
			Repository: &repository,
		}

		err := transformer.Execute()
		Expect(err).NotTo(HaveOccurred())
		Expect(repository.CheckedHeaderIDs).To(ContainElement(headerOne.Id))
		Expect(repository.CheckedHeaderIDs).To(ContainElement(headerTwo.Id))
	})

	It("returns error if marking header checked returns err", func() {
		repository.SetMissingHeaders([]core.Header{headerOne, headerTwo})
		repository.SetCheckedHeaderError(fakes.FakeError)
		fetcher.SetFetchedLogs([]types.Log{})
		transformer := vat_move.VatMoveTransformer{
			Config:     config,
			Converter:  &converter,
			Fetcher:    &fetcher,
			Repository: &repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists vat move model", func() {
		fetcher.SetFetchedLogs([]types.Log{test_data.EthVatMoveLog})
		repository.SetMissingHeaders([]core.Header{headerOne})
		transformer := vat_move.VatMoveTransformer{
			Config:     config,
			Fetcher:    &fetcher,
			Converter:  &converter,
			Repository: &repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(headerOne.Id))
		Expect(repository.PassedModels).To(Equal([]vat_move.VatMoveModel{test_data.VatMoveModel}))
	})

	It("returns error if repository returns error for create", func() {
		fetcher.SetFetchedLogs([]types.Log{test_data.EthVatMoveLog})
		repository.SetMissingHeaders([]core.Header{headerOne})
		repository.SetCreateError(fakes.FakeError)
		transformer := vat_move.VatMoveTransformer{
			Config:     config,
			Fetcher:    &fetcher,
			Converter:  &converter,
			Repository: &repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})
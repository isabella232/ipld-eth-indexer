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

package ilk_test

import (
	"database/sql"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/drip_file/ilk"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Drip file ilk repository", func() {
	var (
		db                    *postgres.DB
		dripFileIlkRepository ilk.DripFileIlkRepository
		err                   error
		headerRepository      datastore.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		dripFileIlkRepository = ilk.DripFileIlkRepository{}
		dripFileIlkRepository.SetDB(db)
	})

	Describe("Create", func() {
		var headerID int64

		BeforeEach(func() {
			headerID, err = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())

		})

		It("adds a drip file ilk event", func() {
			err = dripFileIlkRepository.Create(headerID, []interface{}{test_data.DripFileIlkModel})

			Expect(err).NotTo(HaveOccurred())
			var dbDripFileIlk ilk.DripFileIlkModel
			err = db.Get(&dbDripFileIlk, `SELECT ilk, vow, tax, log_idx, tx_idx, raw_log FROM maker.drip_file_ilk WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbDripFileIlk.Ilk).To(Equal(test_data.DripFileIlkModel.Ilk))
			Expect(dbDripFileIlk.Vow).To(Equal(test_data.DripFileIlkModel.Vow))
			Expect(dbDripFileIlk.Tax).To(Equal(test_data.DripFileIlkModel.Tax))
			Expect(dbDripFileIlk.LogIndex).To(Equal(test_data.DripFileIlkModel.LogIndex))
			Expect(dbDripFileIlk.TransactionIndex).To(Equal(test_data.DripFileIlkModel.TransactionIndex))
			Expect(dbDripFileIlk.Raw).To(MatchJSON(test_data.DripFileIlkModel.Raw))
		})

		It("marks header as checked for logs", func() {
			err = dripFileIlkRepository.Create(headerID, []interface{}{test_data.DripFileIlkModel})

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT drip_file_ilk_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("updates the header to checked if checked headers row already exists", func() {
			_, err = db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerID)
			Expect(err).NotTo(HaveOccurred())

			err = dripFileIlkRepository.Create(headerID, []interface{}{test_data.DripFileIlkModel})

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT drip_file_ilk_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("does not duplicate drip file events", func() {
			err = dripFileIlkRepository.Create(headerID, []interface{}{test_data.DripFileIlkModel})
			Expect(err).NotTo(HaveOccurred())

			err = dripFileIlkRepository.Create(headerID, []interface{}{test_data.DripFileIlkModel})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("pq: duplicate key value violates unique constraint"))
		})

		It("removes drip file if corresponding header is deleted", func() {
			err = dripFileIlkRepository.Create(headerID, []interface{}{test_data.DripFileIlkModel})
			Expect(err).NotTo(HaveOccurred())

			_, err = db.Exec(`DELETE FROM headers WHERE id = $1`, headerID)

			Expect(err).NotTo(HaveOccurred())
			var dbDripFileIlk ilk.DripFileIlkModel
			err = db.Get(&dbDripFileIlk, `SELECT ilk, vow, tax, log_idx, tx_idx, raw_log FROM maker.drip_file_ilk WHERE header_id = $1`, headerID)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(sql.ErrNoRows))
		})

		It("returns an error if model is of wrong type", func() {
			err = dripFileIlkRepository.Create(headerID, []interface{}{test_data.WrongModel{}})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("model of type"))
		})
	})

	Describe("MarkHeaderChecked", func() {
		var headerID int64

		BeforeEach(func() {
			headerID, err = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
		})

		It("creates a row for a new headerID", func() {
			err = dripFileIlkRepository.MarkHeaderChecked(headerID)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT drip_file_ilk_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})

		It("updates row when headerID already exists", func() {
			_, err = db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerID)

			err = dripFileIlkRepository.MarkHeaderChecked(headerID)

			Expect(err).NotTo(HaveOccurred())
			var headerChecked bool
			err = db.Get(&headerChecked, `SELECT drip_file_ilk_checked FROM public.checked_headers WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(headerChecked).To(BeTrue())
		})
	})

	Describe("MissingHeaders", func() {
		var (
			startingBlock, endingBlock, dripFileBlock int64
			blockNumbers, headerIDs                   []int64
		)

		BeforeEach(func() {
			startingBlock = rand.Int63()
			dripFileBlock = startingBlock + 1
			endingBlock = startingBlock + 2

			blockNumbers = []int64{startingBlock, dripFileBlock, endingBlock, endingBlock + 1}

			headerIDs = []int64{}

			for _, n := range blockNumbers {
				headerID, err := headerRepository.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				Expect(err).NotTo(HaveOccurred())
				headerIDs = append(headerIDs, headerID)
			}
		})

		It("returns headers with no associated drip file event", func() {
			err := dripFileIlkRepository.MarkHeaderChecked(headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := dripFileIlkRepository.MissingHeaders(startingBlock, endingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(2))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock)))
		})

		It("only treats headers as checked if drip file ilk logs have been checked", func() {
			_, err := db.Exec(`INSERT INTO public.checked_headers (header_id) VALUES ($1)`, headerIDs[1])
			Expect(err).NotTo(HaveOccurred())

			headers, err := dripFileIlkRepository.MissingHeaders(startingBlock, endingBlock)

			Expect(err).NotTo(HaveOccurred())
			Expect(len(headers)).To(Equal(3))
			Expect(headers[0].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(dripFileBlock)))
			Expect(headers[1].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(dripFileBlock)))
			Expect(headers[2].BlockNumber).To(Or(Equal(startingBlock), Equal(endingBlock), Equal(dripFileBlock)))
		})

		It("only returns headers associated with the current node", func() {
			dbTwo := test_config.NewTestDB(core.Node{ID: "second"})
			headerRepositoryTwo := repositories.NewHeaderRepository(dbTwo)
			for _, n := range blockNumbers {
				_, err = headerRepositoryTwo.CreateOrUpdateHeader(fakes.GetFakeHeader(n))
				Expect(err).NotTo(HaveOccurred())
			}
			dripFileIlkRepositoryTwo := ilk.DripFileIlkRepository{}
			dripFileIlkRepositoryTwo.SetDB(dbTwo)
			err := dripFileIlkRepository.MarkHeaderChecked(headerIDs[0])
			Expect(err).NotTo(HaveOccurred())

			nodeOneMissingHeaders, err := dripFileIlkRepository.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeOneMissingHeaders)).To(Equal(len(blockNumbers) - 1))

			nodeTwoMissingHeaders, err := dripFileIlkRepositoryTwo.MissingHeaders(blockNumbers[0], blockNumbers[len(blockNumbers)-1])
			Expect(err).NotTo(HaveOccurred())
			Expect(len(nodeTwoMissingHeaders)).To(Equal(len(blockNumbers)))
		})
	})
})
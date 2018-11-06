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

package chop_lump

import (
	"fmt"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type CatFileChopLumpRepository struct {
	db *postgres.DB
}

func (repository CatFileChopLumpRepository) Create(headerID int64, models []interface{}) error {
	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}

	for _, model := range models {
		chopLump, ok := model.(CatFileChopLumpModel)
		if !ok {
			tx.Rollback()
			return fmt.Errorf("model of type %T, not %T", model, CatFileChopLumpModel{})
		}

		_, err := tx.Exec(
			`INSERT into maker.cat_file_chop_lump (header_id, ilk, what, data, tx_idx, log_idx, raw_log)
			VALUES($1, $2, $3, $4::NUMERIC, $5, $6, $7)`,
			headerID, chopLump.Ilk, chopLump.What, chopLump.Data, chopLump.TransactionIndex, chopLump.LogIndex, chopLump.Raw,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	_, err = tx.Exec(`INSERT INTO public.checked_headers (header_id, cat_file_chop_lump_checked)
			VALUES ($1, $2)
		ON CONFLICT (header_id) DO
			UPDATE SET cat_file_chop_lump_checked = $2`, headerID, true)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (repository CatFileChopLumpRepository) MarkHeaderChecked(headerID int64) error {
	_, err := repository.db.Exec(`INSERT INTO public.checked_headers (header_id, cat_file_chop_lump_checked)
			VALUES ($1, $2)
		ON CONFLICT (header_id) DO
			UPDATE SET cat_file_chop_lump_checked = $2`, headerID, true)
	return err
}

func (repository CatFileChopLumpRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	var result []core.Header
	err := repository.db.Select(
		&result,
		`SELECT headers.id, headers.block_number FROM headers
			LEFT JOIN checked_headers on headers.id = header_id
		WHERE (header_id ISNULL OR cat_file_chop_lump_checked IS FALSE)
			AND headers.block_number >= $1
			AND headers.block_number <= $2
			AND headers.eth_node_fingerprint = $3`,
		startingBlockNumber,
		endingBlockNumber,
		repository.db.Node.ID,
	)
	return result, err
}

func (repository *CatFileChopLumpRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
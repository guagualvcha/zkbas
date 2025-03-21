/*
 * Copyright © 2021 ZkBAS Protocol
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package blockwitness

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/bnb-chain/zkbas/types"
)

type (
	BlockWitnessModel interface {
		CreateBlockWitnessTable() error
		DropBlockWitnessTable() error
		GetLatestBlockWitnessHeight() (blockNumber int64, err error)
		GetBlockWitnessByNumber(height int64) (witness *BlockWitness, err error)
		UpdateBlockWitnessStatus(witness *BlockWitness, status int64) error
		GetLatestBlockWitness() (witness *BlockWitness, err error)
		CreateBlockWitness(witness *BlockWitness) error
	}

	defaultBlockWitnessModel struct {
		table string
		DB    *gorm.DB
	}

	BlockWitness struct {
		gorm.Model
		Height      int64 `gorm:"index:idx_height,unique"`
		WitnessData string
		Status      int64
	}
)

func NewBlockWitnessModel(db *gorm.DB) BlockWitnessModel {
	return &defaultBlockWitnessModel{
		table: TableName,
		DB:    db,
	}
}

func (*BlockWitness) TableName() string {
	return TableName
}

func (m *defaultBlockWitnessModel) CreateBlockWitnessTable() error {
	return m.DB.AutoMigrate(BlockWitness{})
}

func (m *defaultBlockWitnessModel) DropBlockWitnessTable() error {
	return m.DB.Migrator().DropTable(m.table)
}

func (m *defaultBlockWitnessModel) GetLatestBlockWitnessHeight() (blockNumber int64, err error) {
	var row *BlockWitness
	dbTx := m.DB.Table(m.table).Order("height desc").Limit(1).Find(&row)
	if dbTx.Error != nil {
		return 0, types.DbErrSqlOperation
	} else if dbTx.RowsAffected == 0 {
		return 0, types.DbErrNotFound
	}
	return row.Height, nil
}

func (m *defaultBlockWitnessModel) GetLatestBlockWitness() (witness *BlockWitness, err error) {
	dbTx := m.DB.Table(m.table).Where("status = ?", StatusPublished).Order("height asc").Limit(1).Find(&witness)
	if dbTx.Error != nil {
		return nil, types.DbErrSqlOperation
	} else if dbTx.RowsAffected == 0 {
		return nil, types.DbErrNotFound
	}
	return witness, nil
}

func (m *defaultBlockWitnessModel) GetBlockWitnessByNumber(height int64) (witness *BlockWitness, err error) {
	dbTx := m.DB.Table(m.table).Where("height = ?", height).Limit(1).Find(&witness)
	if dbTx.Error != nil {
		return nil, types.DbErrSqlOperation
	} else if dbTx.RowsAffected == 0 {
		return nil, types.DbErrNotFound
	}
	return witness, nil
}

func (m *defaultBlockWitnessModel) CreateBlockWitness(witness *BlockWitness) error {
	if witness.Height > 1 {
		_, err := m.GetBlockWitnessByNumber(witness.Height - 1)
		if err != nil {
			return fmt.Errorf("previous witness does not exist")
		}
	}

	dbTx := m.DB.Table(m.table).Create(witness)
	if dbTx.Error != nil {
		return types.DbErrSqlOperation
	}
	return nil
}

func (m *defaultBlockWitnessModel) UpdateBlockWitnessStatus(witness *BlockWitness, status int64) error {
	witness.Status = status
	witness.UpdatedAt = time.Now()
	dbTx := m.DB.Table(m.table).Save(witness)
	if dbTx.Error != nil {
		return types.DbErrSqlOperation
	}
	return nil
}

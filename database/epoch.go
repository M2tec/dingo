// Copyright 2025 Blink Labs Software
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package database

type Epoch struct {
	ID uint `gorm:"primarykey"`
	// NOTE: we would normally use this as the primary key, but GORM doesn't
	// like a primary key value of 0
	EpochId       uint64 `gorm:"uniqueIndex"`
	StartSlot     uint64
	Nonce         []byte
	EraId         uint
	SlotLength    uint
	LengthInSlots uint
}

func (Epoch) TableName() string {
	return "epoch"
}

func GetEpochLatest(db *Database) (Epoch, error) {
	return db.GetEpochLatest(nil)
}

func (d *Database) GetEpochLatest(txn *Txn) (Epoch, error) {
	tmpEpoch := Epoch{}
	if txn == nil {
		txn = d.Transaction(false)
	}
	epoch, err := txn.DB().Metadata().GetEpochLatest(txn.Metadata())
	if err != nil {
		return tmpEpoch, err
	}
	tmpEpoch.ID = epoch.ID
	tmpEpoch.EpochId = epoch.EpochId
	tmpEpoch.StartSlot = epoch.StartSlot
	tmpEpoch.Nonce = epoch.Nonce
	tmpEpoch.EraId = epoch.EraId
	tmpEpoch.SlotLength = epoch.SlotLength
	tmpEpoch.LengthInSlots = epoch.LengthInSlots
	return tmpEpoch, nil
}

func GetEpochsByEra(db *Database, eraId uint) ([]Epoch, error) {
	return db.GetEpochsByEra(eraId, nil)
}

func (d *Database) GetEpochsByEra(eraId uint, txn *Txn) ([]Epoch, error) {
	tmpEpochs := []Epoch{}
	if txn == nil {
		txn = d.Transaction(false)
	}
	epochs, err := txn.DB().Metadata().GetEpochsByEra(eraId, txn.Metadata())
	if err != nil {
		return tmpEpochs, err
	}
	for _, epoch := range epochs {
		tmpEpoch := Epoch{
			ID:            epoch.ID,
			EpochId:       epoch.EpochId,
			StartSlot:     epoch.StartSlot,
			Nonce:         epoch.Nonce,
			EraId:         epoch.EraId,
			LengthInSlots: epoch.LengthInSlots,
		}
		tmpEpochs = append(tmpEpochs, tmpEpoch)
	}
	return tmpEpochs, nil
}

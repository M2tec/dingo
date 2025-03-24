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

package dgraph

import (
	"github.com/blinklabs-io/dingo/database/plugin/metadata/dgraph/models"
	"github.com/dgraph-io/dgo/v240"
)

// GetPParams returns a list of protocol parameters for a given epoch
func (d *MetadataStoreDgraph) GetPParams(
	epoch uint64,
	txn *dgo.Dgraph,
) ([]models.PParams, error) {
	ret := []models.PParams{}
	if txn != nil {
		result := txn.Where("epoch = ?", epoch).Order("id DESC").Find(&ret)
		if result.Error != nil {
			return ret, result.Error
		}
	} else {
		result := d.DB().Where("epoch = ?", epoch).Order("id DESC").Find(&ret)
		if result.Error != nil {
			return ret, result.Error
		}
	}
	return ret, nil
}

// GetPParamUpdates returns a list of protocol parameter updates for a given epoch
func (d *MetadataStoreDgraph) GetPParamUpdates(
	epoch uint64,
	txn *dgo.Dgraph,
) ([]models.PParamUpdate, error) {
	ret := []models.PParamUpdate{}
	if txn != nil {
		result := txn.Where("epoch = ?", epoch).Order("id DESC").Find(&ret)
		if result.Error != nil {
			return ret, result.Error
		}
	} else {
		result := d.DB().Where("epoch = ?", epoch).Order("id DESC").Find(&ret)
		if result.Error != nil {
			return ret, result.Error
		}
	}
	return ret, nil
}

// SetPParams saves protocol parameters
func (d *MetadataStoreDgraph) SetPParams(
	params []byte,
	slot, epoch uint64,
	eraId uint,
	txn *dgo.Dgraph,
) error {
	tmpItem := models.PParams{
		Cbor:      params,
		AddedSlot: slot,
		Epoch:     epoch,
		EraId:     eraId,
	}
	if txn != nil {
		if result := txn.Create(&tmpItem); result.Error != nil {
			return result.Error
		}
	} else {
		if result := d.DB().Create(&tmpItem); result.Error != nil {
			return result.Error
		}
	}
	return nil
}

// SetPParamUpdate saves a protocol parameter update
func (d *MetadataStoreDgraph) SetPParamUpdate(
	genesis, update []byte,
	slot, epoch uint64,
	txn *dgo.Dgraph,
) error {
	tmpItem := models.PParamUpdate{
		GenesisHash: genesis,
		Cbor:        update,
		AddedSlot:   slot,
		Epoch:       epoch,
	}
	if txn != nil {
		if result := txn.Create(&tmpItem); result.Error != nil {
			return result.Error
		}
	} else {
		if result := d.DB().Create(&tmpItem); result.Error != nil {
			return result.Error
		}
	}
	return nil
}

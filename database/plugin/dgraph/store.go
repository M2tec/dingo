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

package metadata

import (
	"log/slog"

	"github.com/blinklabs-io/dingo/database/plugin/dgraph/dgraph"
	"github.com/blinklabs-io/dingo/database/plugin/dgraph/dgraph/models"
	"github.com/blinklabs-io/gouroboros/ledger"
	lcommon "github.com/blinklabs-io/gouroboros/ledger/common"
	ochainsync "github.com/blinklabs-io/gouroboros/protocol/chainsync"
	"github.com/dgraph-io/dgo/v240"
)

type MetadataStore interface {
	// Database
	Close() error
	DB() *dgo.Dgraph
	GetCommitTimestamp() (int64, error)
	SetCommitTimestamp(*dgo.Dgraph, int64) error
	Transaction() *dgo.Dgraph

	// Ledger state
	GetPoolRegistrations(
		lcommon.PoolKeyHash,
		*dgo.Dgraph,
	) ([]lcommon.PoolRegistrationCertificate, error)
	GetStakeRegistrations(
		[]byte, // stakeKey
		*dgo.Dgraph,
	) ([]lcommon.StakeRegistrationCertificate, error)
	GetTip(*dgo.Dgraph) (ochainsync.Tip, error)

	GetPParams(
		uint64, // epoch
		*dgo.Dgraph,
	) ([]models.PParams, error)
	GetPParamUpdates(
		uint64, // epoch
		*dgo.Dgraph,
	) ([]models.PParamUpdate, error)
	GetUtxo(
		[]byte, // txId
		uint32, // idx
		*dgo.Dgraph,
	) (models.Utxo, error)

	SetEpoch(
		uint64, // epoch
		uint64, // slot
		[]byte, // nonce
		uint, // era
		uint, // slotLength
		uint, // lengthInSlots
		*dgo.Dgraph,
	) error
	SetPoolRegistration(
		*lcommon.PoolRegistrationCertificate,
		uint64, // slot
		uint64, // deposit
		*dgo.Dgraph,
	) error
	SetPoolRetirement(
		*lcommon.PoolRetirementCertificate,
		uint64, // slot
		*dgo.Dgraph,
	) error
	SetPParams(
		[]byte, // pparams
		uint64, // slot
		uint64, // epoch
		uint, // era
		*dgo.Dgraph,
	) error
	SetPParamUpdate(
		[]byte, // genesis
		[]byte, // update
		uint64, // slot
		uint64, // epoch
		*dgo.Dgraph,
	) error
	SetStakeDelegation(
		*lcommon.StakeDelegationCertificate,
		uint64, // slot
		*dgo.Dgraph,
	) error
	SetStakeDeregistration(
		*lcommon.StakeDeregistrationCertificate,
		uint64, // slot
		*dgo.Dgraph,
	) error
	SetStakeRegistration(
		*lcommon.StakeRegistrationCertificate,
		uint64, // slot
		uint64, // deposit
		*dgo.Dgraph,
	) error
	SetTip(
		ochainsync.Tip,
		*dgo.Dgraph,
	) error

	// Helpers
	GetEpochLatest(*dgo.Dgraph) (models.Epoch, error)
	GetEpochsByEra(uint, *dgo.Dgraph) ([]models.Epoch, error)
	GetUtxosByAddress(ledger.Address, *dgo.Dgraph) ([]models.Utxo, error)
	DeleteUtxo(any, *dgo.Dgraph) error
	DeleteUtxos([]any, *dgo.Dgraph) error
}

// For now, this always returns a dgraph plugin
func New(
	pluginName, dataDir string,
	logger *slog.Logger,
) (MetadataStore, error) {
	dgraphPort := 9080
	return dgraph.New(dgraphPort, logger)
}

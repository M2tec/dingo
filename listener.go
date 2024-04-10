// Copyright 2024 Blink Labs Software
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

package node

import (
	"fmt"
	"net"

	ouroboros "github.com/blinklabs-io/gouroboros"
	"github.com/blinklabs-io/gouroboros/protocol/peersharing"
	"github.com/blinklabs-io/gouroboros/protocol/txsubmission"
)

type ListenerConfig struct {
	UseNtC   bool
	Listener net.Listener
	// TODO
}

func (n *Node) startListener(l ListenerConfig) {
	defaultConnOpts := []ouroboros.ConnectionOptionFunc{
		ouroboros.WithNetworkMagic(n.config.networkMagic),
		ouroboros.WithNodeToNode(!l.UseNtC),
		ouroboros.WithServer(true),
	}
	if l.UseNtC {
		// Node-to-client
		defaultConnOpts = append(
			defaultConnOpts,
			// TODO: add localtxsubmission
			// TODO: add localstatequery
			// TODO: add localtxmonitor
		)
	} else {
		// Node-to-node
		defaultConnOpts = append(
			defaultConnOpts,
			// Peer sharing
			ouroboros.WithPeerSharing(n.config.peerSharing),
			ouroboros.WithPeerSharingConfig(
				peersharing.NewConfig(
					peersharing.WithShareRequestFunc(n.peersharingShareRequest),
				),
			),
			// TxSubmission
			ouroboros.WithTxSubmissionConfig(
				txsubmission.NewConfig(
					txsubmission.WithInitFunc(n.txsubmissionServerInit),
				),
			),
			// TODO: add chain-sync
			// TODO: add block-fetch
		)
	}
	for {
		// Accept connection
		conn, err := l.Listener.Accept()
		if err != nil {
			n.config.logger.Error(fmt.Sprintf("accept failed: %s", err))
			continue
		}
		n.config.logger.Info(fmt.Sprintf("accepted connection from %s", conn.RemoteAddr()))
		// Setup Ouroboros connection
		connOpts := append(
			defaultConnOpts,
			ouroboros.WithConnection(conn),
		)
		oConn, err := ouroboros.NewConnection(connOpts...)
		if err != nil {
			n.config.logger.Error(fmt.Sprintf("failed to setup connection: %s", err))
			continue
		}
		// Add to connection manager
		// TODO: add tags for connection for later tracking
		n.connManager.AddConnection(oConn)
	}
}

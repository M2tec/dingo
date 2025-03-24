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
	"fmt"
	"io"
	"log/slog"

	"github.com/blinklabs-io/dingo/database/plugin"
	"github.com/dgraph-io/dgo/v240"
	"github.com/dgraph-io/dgo/v240/protos/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Register plugin
func init() {
	plugin.Register(
		plugin.PluginEntry{
			Type: plugin.PluginTypeMetadata,
			Name: "dgraph",
		},
	)
}

// MetadataStoreDgraph stores all data in dgraph. Data may not be persisted
type MetadataStoreDgraph struct {
	dgraphPort int
	db         *dgo.Dgraph
	logger     *slog.Logger
}

// New creates a new database
func New(
	dgraphPort int,
	logger *slog.Logger,
) (*MetadataStoreDgraph, error) {
	var dgraphDb *dgo.Dgraph
	var err error

	// Open dGraph db using dgo
	conn, err := grpc.NewClient("localhost:9080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Conn: ", conn)
	// defer conn.Close()
	dgraphDb = dgo.NewDgraphClient(api.NewDgraphClient(conn))

	db := &MetadataStoreDgraph{
		db:         dgraphDb,
		dgraphPort: dgraphPort,
		logger:     logger,
	}
	if err := db.init(); err != nil {
		// MetadataStoreDgraph is available for recovery, so return it with error
		return db, err
	}

	// Create table schemas
	// db.logger.Debug(fmt.Sprintf("creating table: %#v", &CommitTimestamp{}))
	// if err := db.db.AutoMigrate(&CommitTimestamp{}); err != nil {
	// 	return db, err
	// }
	// for _, model := range models.MigrateModels {
	// 	db.logger.Debug(fmt.Sprintf("creating table: %#v", model))
	// 	if err := db.db.AutoMigrate(model); err != nil {
	// 		return db, err
	// 	}
	// }
	return db, nil
}

func (d *MetadataStoreDgraph) init() error {
	if d.logger == nil {
		// Create logger to throw away logs
		// We do this so we don't have to add guards around every log operation
		d.logger = slog.New(slog.NewJSONHandler(io.Discard, nil))
	}
	// Configure tracing for GORM
	// if err := d.db.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
	// 	return err
	// }
	return nil
}

// AutoMigrate wraps the gorm AutoMigrate
// func (d *MetadataStoreDgraph) AutoMigrate(dst ...interface{}) error {
// 	return d.DB().AutoMigrate(dst...)
// }

// Close gets the database handle from our MetadataStore and closes it
func (d *MetadataStoreDgraph) Close() error {
	// get DB handle from gorm.DB
	db, err := d.DB().DB()
	if err != nil {
		return err
	}
	return db.Close()
}

// Create creates a record
// func (d *MetadataStoreDgraph) Create(value interface{}) *dgo.Dgraph {
// 	return d.DB().Create(value)
// }

// DB returns the database handle
func (d *MetadataStoreDgraph) DB() *dgo.Dgraph {
	return d.db
}

// First returns the first DB entry
// func (d *MetadataStoreDgraph) First(args interface{}) *dgo.Dgraph {
// 	return d.DB().First(args)
// }

// Order orders a DB query
// func (d *MetadataStoreDgraph) Order(args interface{}) *dgo.Dgraph {
// 	return d.DB().Order(args)
// }

// Transaction creates a gorm transaction
// func (d *MetadataStoreDgraph) Transaction() *dgo.Dgraph {
// 	return d.DB().Begin()
// }

// Where constrains a DB query
// func (d *MetadataStoreDgraph) Where(
// 	query interface{},
// 	args ...interface{},
// ) *dgo.Dgraph {
// 	return d.DB().Where(query, args...)
// }

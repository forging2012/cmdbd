// Copyright 2017 John Scherff
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

package main

import (
	`encoding/json`
	`path/filepath`
	`os`
)

// Configuration aliases.
var (
	db *Database
	qy *Queries
	ws *Server
	sl *Syslog
	sl, al, el *Logger
)

// Config contains infomation about the server process and log writers.
type Config struct {

	ConfigDir string
	SerialFmt string
	MetaUsbUrl string
	Configs   map[string]string

	Database *Database
	Queries  *Queries
	Syslog   *Syslog
	Logger   *Logger
	Router   *Router
	MetaUsb	 *MetaUsb
	Server   *Server
}

// NewConfig creates a new Config object and reads its config
// from the provided JSON configuration file.
func NewConfig(cf string) (this *Config, err error) {

	// Load the base configuration needed to load remaining configs.

	if this, err := loadConfig(cf); err != nil {
		return nil, err
	}

	// Prepend configuration directory to configuration filenames.

	for _, fn := range this.Configs {
		fn = filepath.Join(this.ConfigDir, fn)
	}

	// Create and initialize Database object.

	if database, err := NewDatabase(this.Configs[Database]); err != nil {
		return nil, err
	} else {
		this.Database = database
	}

	db = this.Database.DB

	// Create and initialize Queries object.

	if queries, err := NewQueries(this.Configs[Queries], db); err != nil {
		return nil, err
	} else {
		this.Queries = queries
	}

	qy = this.Queries

	// Create and initialize Syslog object.

	if syslog, err := NewSyslog(this.Configs[syslog]); err != nil {
		return nil, err
	} else {
		this.Syslog = syslog
	}

	sl = this.Syslog

	// Create and initialize Logger object.

	if logger, err := NewLogger(this.Configs[Logger], sl); err != nil {
		return nil, err
	} else {
		this.Logger = logger
	}

	sl = this.Logger.Logs[`system`]
	al = this.Logger.Logs[`access`]
	el = this.Logger.Logs[`error`]

	// Create and initialize Router object.

	if router, err := NewRouter(this.Configs[Router], al, el); err != nil {
		return nil, err
	} else {
		router.AddRoutes(usbCiRoutes).AddRoutes(usbMetaRoutes)
		this.Router = router
	}

	// Create and initialize MetaUsb object.

	if metausb, err := NewMetaUsb(this.Configs[MetaUsb]); err != nil {
		return nil, err
	} else {
		this.MetaUsb = metausb
	}

	// Create and initialize Server, Router, and Routes objects.

	if server, err := NewServer(this.Configs[Server]); err != nil {
		return nil, err
	} else {
		server.Handler = this.Router
		this.Server = server
	}

	ws = this.Server

	return this, nil
}

// loadConfig loads a JSON configuration file into an object.
func loadConfig(t interface{}, cf string) (interface{}, error) {

	if fh, err := os.Open(cf); err != nil {
		return nil, err
	} else {
		defer fh.Close()
		jd := json.NewDecoder(fh)
		err = jd.Decode(&t)
		return t, err
	}
}
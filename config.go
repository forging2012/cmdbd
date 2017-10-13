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
	`os`
)

// Config contains infomation about the server process and log writers.
type Config struct {
	Server   *Server
	Database *Database
	Queries  *Queries
	Router   *Router
	Logger   *Logger
	Syslog   *Syslog
	MetaCi   *MetaCi
}

// NewConfig creates a new Config object and reads its configuration from
// the provided JSON configuration file.
func NewConfig(cf string) (this *Config, err error) {

	fh, err := os.Open(cf)

	if err != nil {
		return this, err
	}

	defer fh.Close()

	this = &Config{}
	jd := json.NewDecoder(fh)
	err = jd.Decode(&this)

	return this, err
}

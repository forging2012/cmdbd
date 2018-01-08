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

package usbci

import (
	`github.com/jscherff/cmdbd/api`
	`github.com/jscherff/cmdbd/api/v2/cmdb/usbci`
)

// Endpoints is a collection of URL path to handler function mappings.
var Endpoints = api.Endpoints {

	api.Endpoint {
		Name:		`USBCI CheckIn Handler`,
		Path:		`/v1/usbci/checkin/{host}/{vid}/{pid}`,
		Method:		`POST`,
		HandlerFunc:	usbci.CheckIn,
		Protected:	true,
	},

	api.Endpoint {
		Name:		`USBCI CheckOut Handler`,
		Path:		`/v1/usbci/checkout/{host}/{vid}/{pid}/{sn}`,
		Method:		`GET`,
		HandlerFunc:	usbci.CheckOut,
		Protected:	true,
	},

	api.Endpoint {
		Name:		`USBCI NewSn Handler`,
		Path:		`/v1/usbci/newsn/{host}/{vid}/{pid}`,
		Method:		`POST`,
		HandlerFunc:	usbci.NewSn,
		Protected:	true,
	},

	api.Endpoint {
		Name:		`USBCI Audit Handler`,
		Path:		`/v1/usbci/audit/{host}/{vid}/{pid}/{sn}`,
		Method:		`POST`,
		HandlerFunc:	Audit,
		Protected:	true,
	},
}

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

package usbmeta

import (
	`github.com/jscherff/cmdbd/api`
	v2 `github.com/jscherff/cmdbd/api/v2/cmdb/usbmeta`
)

// Routes is a collection of HTTP verb/path-to-handler-function mappings.
var Routes = api.Routes {

	api.Route {
		Name:		`Metadata USB Vendor Handler`,
		Path:		`/v1/usbmeta/vendor/{vid}`,
		Method:		`GET`,
		HandlerFunc:	v2.Vendor,
		Protected:	false,
	},

	api.Route {
		Name:		`Metadata USB Product Handler`,
		Path:		`/v1/usbmeta/vendor/{vid}/{pid}`,
		Method:		`GET`,
		HandlerFunc:	v2.Product,
		Protected:	false,
	},

	api.Route {
		Name:		`Metadata USB Class Handler`,
		Path:		`/v1/usbmeta/class/{cid}`,
		Method:		`GET`,
		HandlerFunc:	v2.Class,
		Protected:	false,
	},

	api.Route {
		Name:		`Metadata USB SubClass Handler`,
		Path:		`/v1/usbmeta/subclass/{cid}/{sid}`,
		Method:		`GET`,
		HandlerFunc:	v2.SubClass,
		Protected:	false,
	},

	api.Route {
		Name:		`Metadata USB Protocol Handler`,
		Path:		`/v1/usbmeta/protocol/{cid}/{sid}/{pid}`,
		Method:		`GET`,
		HandlerFunc:	v2.Protocol,
		Protected:	false,
	},
}

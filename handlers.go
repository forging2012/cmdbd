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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/jscherff/gocmdb/usbci"
)

// USBDeviceHandler handles various 'actions' for device gocmdb agents.
func USBDeviceHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	action := vars["action"]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, HttpBodySizeLimit))

	if err != nil {
		panic(ErrorDecorator(err))
	}

	if err := r.Body.Close(); err != nil {
		panic(ErrorDecorator(err))
	}

	w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")

	dev := usbci.NewWSAPI()

	if err := json.Unmarshal(body, &dev); err != nil {

		conf.Log.Writer[Error].WriteError(ErrorDecorator(err))
		w.WriteHeader(http.StatusUnprocessableEntity)

		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(ErrorDecorator(err))
		}

		return
	}

	w.Header().Set("Content-Type", "applicaiton/json; charset=UTF8")

	switch action {

	case "serial":

		var sn = dev.GetSerialNum()

		if len(sn) != 0 {
			err = fmt.Errorf("device already has serial number %q", sn)
			conf.Log.Writer[Error].WriteError(ErrorDecorator(err))
			w.WriteHeader(http.StatusNoContent)
			return
		}

		var id int64

		if id, err = storeDevice(conf.Database.Stmt["SerialInsert"], dev); err != nil {
			conf.Log.Writer[Error].WriteError(ErrorDecorator(err))
		} else {
			sn = fmt.Sprintf("24F%04x", id)
			_, err = updateSerial(conf.Database.Stmt["SerialUpdate"], sn, id)
		}

		if err == nil {
			w.WriteHeader(http.StatusCreated)

			if err = json.NewEncoder(w).Encode(sn); err != nil {
				panic(ErrorDecorator(err))
			}
		}

	case "checkin":

		if _, err = storeDevice(conf.Database.Stmt["CheckinInsert"], dev); err == nil {
			w.WriteHeader(http.StatusAccepted)
		}

	case "audit":

		if _, err = storeDevice(conf.Database.Stmt["CheckinInsert"], dev); err != nil {
			conf.Log.Writer[Error].WriteError(ErrorDecorator(err))
		}

		if err = storeAudit(conf.Database.Stmt["AuditInsert"], dev); err == nil {
			w.WriteHeader(http.StatusAccepted)
		}
	}

	if err != nil {
		conf.Log.Writer[Error].WriteError(ErrorDecorator(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}


// AllowedMethodHandler restricts requests to methods listed in the AllowedMethods
// slice in the systemwide configuration.
func AllowedMethodHandler(h http.Handler, methods ...string) http.Handler {

	return http.HandlerFunc(

		func(w http.ResponseWriter, r *http.Request) {

			for _, m := range methods {
				if r.Method == m {
					h.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, fmt.Sprintf("Unsupported method %q", r.Method),
				http.StatusMethodNotAllowed)
		},
	)
}

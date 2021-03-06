
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

package cmdb

import (
	`encoding/json`
	`fmt`
	`net/http`
	`github.com/gorilla/mux`
	`github.com/jscherff/cmdbd/api`
	`github.com/jscherff/cmdbd/model/cmdb`
	`github.com/jscherff/cmdbd/service`
)

// Templates for system and error messages.
const (
	fmtHostInfo = `host '%s' at '%s'`
	fmtAuthMissingCreds = `auth failure - missing credentials`
	fmtAuthUserNotFound = `auth failure - user '%s' not found: %v`
	fmtAuthSuccess = `auth success - issuing token to '%s' on ` + fmtHostInfo
	fmtEventSuccess = `event logged for ` + fmtHostInfo
	fmtHealthSuccess = `health check success for host at '%s'`
)

// Package variables required for operation.
var (
	authSvc service.AuthSvc
	loggerSvc service.LoggerSvc
)

// Init initializes the package variables required for operation.
func Init(as service.AuthSvc, ls service.LoggerSvc) {
	authSvc, loggerSvc = as, ls
}

// SetauthToken authenticates client using basic authentication and
// issues a token for API authentication if successful.
func SetAuthToken(w http.ResponseWriter, r *http.Request) {

	var (
		ok bool
		passwd string
	)

	vars := mux.Vars(r)
	user := &cmdb.User{}

	if user.Username, passwd, ok = r.BasicAuth(); !ok {

		err := fmt.Errorf(fmtAuthMissingCreds)
		loggerSvc.ErrorLog().Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusUnauthorized)

	} else if err := user.Read(); err != nil {

		err = fmt.Errorf(fmtAuthUserNotFound, user.Username, err)
		loggerSvc.ErrorLog().Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusNotFound)

	} else if err := user.VerifyPassword(passwd); err != nil {

		loggerSvc.ErrorLog().Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusUnauthorized)

	} else if err := user.VerifyAccess(); err != nil {

		loggerSvc.ErrorLog().Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusUnauthorized)

	} else if token, err := authSvc.CreateToken(user); err != nil {

		loggerSvc.ErrorLog().Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else if tokenString, err := authSvc.CreateTokenString(token); err != nil {

		loggerSvc.ErrorLog().Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else if cookie, err := authSvc.CreateCookie(tokenString); err != nil {

		loggerSvc.ErrorLog().Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		loggerSvc.SystemLog().Printf(fmtAuthSuccess, user.Username, vars[`host`], r.RemoteAddr)
		http.SetCookie(w, cookie)
	}
}

// CreateEvent writes an event to the datastore event log.
func CreateEvent(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	event := &cmdb.Event{}
	event.HostName, event.RemoteAddr = vars[`host`], r.RemoteAddr

	if _, err := event.Create(); err != nil {

		loggerSvc.ErrorLog().Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		loggerSvc.SystemLog().Printf(fmtEventSuccess, event.HostName, event.RemoteAddr)
		w.WriteHeader(http.StatusCreated)
	}
}

// CheckHealth returns the health of the server.
func CheckHealth(w http.ResponseWriter, r *http.Request) {

	info := &cmdb.Info{
		Client:	r.RemoteAddr,
		Server:	r.URL.Host,
		Proto:	r.Proto,
		Method:	r.Method,
		Scheme:	r.URL.Scheme,
		Path:	r.URL.Path,
	}

	w.Header().Set(`Content-Type`, `application/json; charset=UTF8`)

	if err := info.Read(); err != nil {

		loggerSvc.ErrorLog().Print(api.AppendRequest(err, r))
		http.Error(w, err.Error(), http.StatusInternalServerError)

	} else {

		loggerSvc.SystemLog().Printf(fmtHealthSuccess, r.RemoteAddr)
		w.WriteHeader(http.StatusOK)

		if err = json.NewEncoder(w).Encode(info); err != nil {
			loggerSvc.ErrorLog().Panic(err)
		}
	}
}

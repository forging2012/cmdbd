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
	"path/filepath"
	"runtime"
	"fmt"
)

// ErrorDecorator prepends function filename, line number, and function name
// to error messages.
func ErrorDecorator(ue error) (de error) {

	var msg string

	pc, file, line, success := runtime.Caller(1)
	function := runtime.FuncForPC(pc)

	if success {
		msg = fmt.Sprintf("%s:%d: %s()", filepath.Base(file), line, function.Name())
	} else {
		msg = "unknown goroutine"
	}

	return fmt.Errorf("%s: %v", msg, ue)
}
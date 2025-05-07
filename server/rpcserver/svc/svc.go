// Copyright 2025 wsserver Author. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package svc

import (
	"errors"
	"reflect"

	"github.com/go-the-way/validator"
)

func Do[REQ, T any](req REQ, t T, callFunc func(req REQ) (err error), errFunc func(err error, t T)) (err0 error) {
	if err := validate[REQ](&req); err != nil {
		errFunc(err, t)
		return
	}
	if err := check[REQ](&req); err != nil {
		errFunc(err, t)
		return
	}
	if err := callFunc(req); err != nil {
		errFunc(err, t)
		return
	}
	return
}

func validate[REQ any](req *REQ) (err error) {
	v := validator.New(req).Lang()
	if vr := v.Validate(); !vr.Passed {
		err = errors.New(vr.Messages())
	}
	return
}

func check[REQ any](req *REQ) (err error) {
loop:
	for _, value := range []reflect.Value{
		reflect.ValueOf(req).MethodByName("Check"),
		reflect.ValueOf(&req).MethodByName("Check"),
	} {
		if value.IsValid() {
			if values := value.Call([]reflect.Value{}); values != nil && len(values) > 0 {
				for _, val := range values {
					if val.CanInterface() {
						if inter := val.Interface(); inter != nil {
							if err = inter.(error); err != nil {
								break loop
							}
						}
					}
				}
			}
		}
	}
	return
}

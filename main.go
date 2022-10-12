// Copyright 2022 wsserver Author. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net/http"

	"github.com/go-the-way/wsserver/config"
	"github.com/go-the-way/wsserver/listener"
	"github.com/go-the-way/wsserver/manager"

	_ "github.com/go-the-way/wsserver/handler"
	_ "github.com/go-the-way/wsserver/rpc"
)

func init() {
	var (
		cCh = make(chan *manager.C, 10000)
		dCh = make(chan *manager.C, 10000)
		_   = listener.NewCreator(cCh)
		_   = listener.NewDestroyer(dCh)
	)
	manager.Init(cCh, dCh)
}

func main() {
	serverAddr := config.GetServerAddr()
	s := &http.Server{Addr: serverAddr, Handler: config.GetServer()}
	fmt.Println("http server started on", serverAddr)
	fmt.Println(s.ListenAndServe())
}

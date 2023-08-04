// Copyright 2023 wsserver Author. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rpc

import (
	"fmt"
	"github.com/go-the-way/wsserver/rpc/service/client"
	"github.com/go-the-way/wsserver/rpc/service/sender"

	"github.com/go-the-way/wsserver/config"
	"github.com/smallnest/rpcx/server"
)

func init() { go serve() }

func serve() {
	s := server.NewServer()
	rpcAddr := config.GetRpcAddr()
	_ = s.Register(new(sender.Sender), "")
	_ = s.Register(new(client.Client), "")
	fmt.Println("rpc server started on", rpcAddr)
	fmt.Println(s.Serve("tcp", rpcAddr))
}

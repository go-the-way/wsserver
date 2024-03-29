// Copyright 2023 wsserver Author. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"context"
	"fmt"

	"github.com/smallnest/rpcx/client"
)

type (
	Args struct {
		Type     string         `json:"type"`
		ClientID []string       `json:"client_id"`
		Data     map[string]any `json:"data"`
	}
	Reply struct {
		Code int `json:"code"`
	}
)

func main() {
	args := Args{Type: "echo", ClientID: []string{"7ed8658c34384722a518e151cb6ccb85"}, Data: map[string]any{"seq": 10000}}
	d, err := client.NewPeer2PeerDiscovery("tcp@:86", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	cc := client.NewXClient("Sender", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer func() { _ = cc.Close() }()
	var reply Reply
	err = cc.Call(context.Background(), "Send", args, &reply)
	if err != nil {
		fmt.Println(err)
		return
	}
}

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

package listener

import (
	"fmt"
	c "github.com/go-the-way/wsserver/manager"
)

// creator 创建监听器
type creator struct{ ch chan *c.C }

func NewCreator(ch chan *c.C) Listener {
	cr := &creator{ch: ch}
	go cr.startCreate()
	return cr
}

func (cr *creator) C() chan<- *c.C { return cr.ch }

func (cr *creator) startCreate() {
	for {
		if client := <-cr.ch; client != nil {
			fmt.Println("created:", client.ID())
			client.Write(&c.WriteProto{Type: "connect", ClientID: client.ID(), Data: nil})
		}
	}
}

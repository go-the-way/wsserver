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

package listener

import (
	"fmt"
	c "github.com/go-the-way/wsserver/manager"
)

// joiner 加入组监听器
type joiner struct{ ch chan *c.C }

func NewJoiner(ch chan *c.C) Listener {
	jr := &joiner{ch: ch}
	go jr.startJoin()
	return jr
}

func (jr *joiner) C() chan<- *c.C { return jr.ch }

func (jr *joiner) startJoin() {
	for {
		if client := <-jr.ch; client != nil {
			fmt.Println("joined:", client.ID())
		}
	}
}

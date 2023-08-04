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

// destroyer 销毁监听器
type destroyer struct{ ch chan *c.C }

func NewDestroyer(ch chan *c.C) Listener {
	dr := &destroyer{ch: ch}
	go dr.startDestroy()
	return dr
}

func (dr *destroyer) C() chan<- *c.C { return dr.ch }

func (dr *destroyer) startDestroy() {
	for {
		client := <-dr.ch
		fmt.Println("destroyed:", client.ID())
	}
}

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

package manager

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	m                   = newClientManager()
	AddClient           = m.addClient
	Send                = m.sendToClient
	needAbortClientIdCh = make(chan string, 100)
)

type clientManager struct {
	mu      *sync.RWMutex
	m       map[string]*client
	closeCh chan string
	sendCh  chan SendProtocol
}

func newClientManager() *clientManager {
	cm := &clientManager{
		mu:      &sync.RWMutex{},
		m:       map[string]*client{},
		closeCh: make(chan string, 10000),
		sendCh:  make(chan SendProtocol, 10000),
	}
	go cm.startReceive()
	go cm.startCloseClient()
	go cm.startAbortTaskToNode()
	return cm
}

func (cm *clientManager) addClient(conn *websocket.Conn, clientId, group string, reqHeader http.Header) {
	c := newClient(conn, clientId, group, reqHeader, cm.closeCh)
	cm.runLock(func() { cm.m[c.id] = c })
}

func (cm *clientManager) sendToClient(proto SendProtocol) { cm.sendCh <- proto }
func (cm *clientManager) startAbortTaskToNode() {
	getTaskAbortProto := func(clientId, abortClientId string) map[string]any {
		return map[string]any{"type": "task_abort", "task_id": "TASK_ID", "task_client_id": "SERVER", "client_id": clientId, "data": map[string]any{"client_id": abortClientId}}
	}
	for {
		clientId := <-needAbortClientIdCh
		cm.runRLock(func() {
			for _, v := range cm.m {
				if v.isNode() {
					v.write(getTaskAbortProto(v.id, clientId))
				}
			}
		})
	}
}

func (cm *clientManager) startReceive() {
	for proto := range cm.sendCh {
		cm.runRLock(func() {
			c, exists := cm.m[proto.ClientID()]
			if !exists {
				return
			}
			c.write(proto.SendData())
		})
	}
}

func (cm *clientManager) closeClient(id string) { cm.closeCh <- id }

func (cm *clientManager) startCloseClient() {
	for clientId := range cm.closeCh {
		cm.runLock(func() {
			c, exists := cm.m[clientId]
			if !exists {
				return
			}
			c.close()
			delete(cm.m, c.id)
		})
	}
}

func (cm *clientManager) runLock(fn func())  { cm.lock(); defer cm.unlock(); fn() }
func (cm *clientManager) runRLock(fn func()) { cm.rLock(); defer cm.rUnlock(); fn() }
func (cm *clientManager) lock()              { cm.mu.Lock() }
func (cm *clientManager) unlock()            { cm.mu.Unlock() }
func (cm *clientManager) rLock()             { cm.mu.RLock() }
func (cm *clientManager) rUnlock()           { cm.mu.RUnlock() }

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

package manager

import (
	"errors"
	"fmt"
	"sync"

	ws "github.com/gorilla/websocket"
)

var (
	manager *clientManager
)

// clientManager 客户端管理器
type (
	clientManager struct {
		M  *sync.Map   // M: WS connection client sync.Map
		mu *sync.Mutex // mu: the sync mutex

		closeCh chan string // closeCh: the client close channel

		createCh, destroyCh chan<- *C // createCh, destroyCh: the client create and destroy channel
		joinCh              chan<- *C // joinCh: the client join group channel

		readCh chan *ReadProto // readCh: the client read channel

		writeCh chan *WriteProto // writeCh: the client write channel

		gReadCh chan *GReadProto // gReadCh: the group read channel

		gWriteCh chan *GWriteProto // gWriteCh: the group write channel

		bcReadCh chan *BCReadProto // bcReadCh: the broadcast read channel

		bcWriteCh chan *BCWriteProto // bcWriteCh: the broadcast write channel
	}
	proto struct {
		Type     string         `json:"type"`
		ClientID string         `json:"client_id"`
		Data     map[string]any `json:"data"`
	}
	gProto struct {
		Type  string         `json:"type"`
		Group string         `json:"group"`
		Data  map[string]any `json:"data"`
	}
	bcProto struct {
		Type string         `json:"type"`
		Data map[string]any `json:"data"`
	}
	ReadProto    proto
	WriteProto   proto
	GReadProto   gProto
	GWriteProto  gProto
	BCReadProto  bcProto
	BCWriteProto bcProto
)

// newClientManager 新建客户端管理器
func newClientManager(createCh, destroyCh, joinCh chan<- *C) *clientManager {
	return &clientManager{
		M:         &sync.Map{},
		mu:        &sync.Mutex{},
		closeCh:   make(chan string, 10000),
		createCh:  createCh,
		destroyCh: destroyCh,
		joinCh:    joinCh,
		readCh:    make(chan *ReadProto, 1000),
		writeCh:   make(chan *WriteProto, 1000),
		gReadCh:   make(chan *GReadProto, 1000),
		gWriteCh:  make(chan *GWriteProto, 1000),
		bcReadCh:  make(chan *BCReadProto, 1000),
		bcWriteCh: make(chan *BCWriteProto, 1000),
	}
}

// Init 初始化客户端管理器
func Init(createCh, destroyCh, joinCh chan<- *C) {
	manager = newClientManager(createCh, destroyCh, joinCh)
	go manager.startClose()
	go manager.startRead()
	go manager.startWrite()
	go manager.startGRead()
	go manager.startGWrite()
	go manager.startBCRead()
	go manager.startBCWrite()
}

// Connect 客户端连接
func Connect(conn *ws.Conn, group string) {
	if _, err := manager.Connect(conn, group); err != nil {
		fmt.Println("Connect", err)
	}
}

// SendToClient 发送客户端消息
func SendToClient(pro *WriteProto) { manager.writeCh <- pro }

// SendToGroup 发送组消息
func SendToGroup(gPro *GWriteProto) { manager.gWriteCh <- gPro }

// Broadcast 广播
func Broadcast(bcPro *BCWriteProto) { manager.bcWriteCh <- bcPro }

func (cm *clientManager) startRead() {
	for {
		select {
		case data, ok := <-cm.readCh:
			if !ok {
				return
			}
			fmt.Println("read:", data)
		}
	}
}

func (cm *clientManager) startWrite() {
	for {
		select {
		case pro := <-cm.writeCh:
			if value, ok := manager.M.Load(pro.ClientID); ok {
				go func(value any) { value.(*C).Write(pro) }(value)
			}
		}
	}
}

func (cm *clientManager) startGRead() {
	for {
		select {
		case data, ok := <-cm.gReadCh:
			if !ok {
				return
			}
			fmt.Println("read:", data)
		}
	}
}

func (cm *clientManager) startGWrite() {
	for {
		select {
		case pro := <-cm.gWriteCh:
			manager.M.Range(func(_, value any) bool {
				cc := value.(*C)
				if cc.groups.Contains(pro.Group) {
					go func(value any) { cc.Write(&WriteProto{pro.Type, cc.ID(), pro.Data}) }(value)
				}
				return true
			})
		}
	}
}

func (cm *clientManager) startBCRead() {
	for {
		select {
		case data, ok := <-cm.bcReadCh:
			if !ok {
				return
			}
			fmt.Println("read:", data)
		}
	}
}

func (cm *clientManager) startBCWrite() {
	for {
		select {
		case pro := <-cm.bcWriteCh:
			manager.M.Range(func(_, value any) bool {
				cc := value.(*C)
				go func(value any) { cc.Write(&WriteProto{pro.Type, cc.ID(), pro.Data}) }(value)
				return true
			})
		}
	}
}

func OnlineClient() []string {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	clients := make([]string, 0)
	manager.M.Range(func(_, v any) bool {
		clients = append(clients, v.(*C).ID())
		return true
	})
	return clients
}

// Connect WS客户端连接
func (cm *clientManager) Connect(conn *ws.Conn, group string) (*C, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	ct := newClient(conn, group, cm.closeCh, cm.readCh)

	clientID := ct.ID()

	if _, ok := cm.M.Load(ct); ok {
		return nil, errors.New("client_manager: client already exists")
	}

	// first store
	cm.M.Store(clientID, ct)

	// then attach to create
	if ch := cm.createCh; ch != nil {
		ch <- ct
	}

	if ch := cm.joinCh; ch != nil {
		ch <- ct
	}

	return ct, nil
}

// startClose 监听WS客户端断开连接
func (cm *clientManager) startClose() {
	for {
		select {
		case clientID, ok := <-cm.closeCh:
			if !ok {
				return
			}
			if value, ok := cm.M.Load(clientID); ok {
				if ct, ok := value.(*C); ok {
					// first to delete
					cm.M.Delete(clientID)
					// then attach to destroy
					if ch := cm.destroyCh; ch != nil {
						ch <- ct
					}
				}
			}
		}
	}
}

func JoinGroup(clientID, group string) error {
	return manager.joinGroup(clientID, group)
}

func (cm *clientManager) joinGroup(clientID, group string) error {
	if value, ok := cm.M.Load(clientID); ok {
		cc := value.(*C)
		if err := cc.JoinGroup(group); err != nil {
			return err
		} else {
			if ch := cm.joinCh; ch != nil {
				ch <- cc
			}
			return nil
		}
	} else {
		return errors.New("client: not exists")
	}
}

func LeaveGroup(clientID, group string) error {
	return manager.leaveGroup(clientID, group)
}

func (cm *clientManager) leaveGroup(clientID, group string) error {
	if value, ok := cm.M.Load(clientID); ok {
		cc := value.(*C)
		return cc.LeaveGroup(group)
	} else {
		return errors.New("client: not exists")
	}
}

func LeaveAllGroup(clientID string) error {
	return manager.leaveAllGroup(clientID)
}

func (cm *clientManager) leaveAllGroup(clientID string) error {
	if value, ok := cm.M.Load(clientID); ok {
		cc := value.(*C)
		cc.LeaveGroups()
		return nil
	} else {
		return errors.New("client: not exists")
	}
}

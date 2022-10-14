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

	"strings"
	"sync"
	"time"

	"github.com/go-the-way/wsserver/pkg/uuid"

	"github.com/go-the-way/streams/ts"

	ws "github.com/gorilla/websocket"
)

// Client WS客户端
type (
	client struct {
		mu   *sync.Mutex // mu: locker
		conn *ws.Conn    // conn: WS connection
		id   string      // id: WS ClientID

		groups ts.Set[string] // groups: The group names

		connectedTime time.Time // connectedTime: WS connection connected time

		closed     bool      // closed: WS connection closed?
		closedTime time.Time // closedTime: WS connection closed time

		closeCh chan<- string

		readCh     chan<- *ReadProto
		stopReadCh chan struct{}

		writeCh     chan *WriteProto
		stopWriteCh chan struct{}

		heartMu     *sync.Mutex
		stopHeartCh chan struct{}

		ticker *time.Ticker
	}
	C = client
)

// newClient 新建客户端
func newClient(conn *ws.Conn, group string, closeCh chan<- string, readCh chan<- *ReadProto) *client {
	c := &client{
		id:            uuid.RandID(),
		mu:            &sync.Mutex{},
		conn:          conn,
		groups:        ts.NewSetValue[string](group),
		connectedTime: time.Now(),
		closeCh:       closeCh,
		readCh:        readCh,
		stopReadCh:    make(chan struct{}, 1),
		writeCh:       make(chan *WriteProto, 1),
		stopWriteCh:   make(chan struct{}, 1),
		heartMu:       &sync.Mutex{},
		stopHeartCh:   make(chan struct{}, 1),
		ticker:        time.NewTicker(time.Second * 10),
	}
	go c.read()
	go c.write()
	go c.heart()
	return c
}

// ID 客户端ID
func (c *client) ID() string { return c.id }

// RemoteAddrIP 远程地址IP
func (c *client) RemoteAddrIP() string { return strings.Split(c.conn.RemoteAddr().String(), ":")[0] }

// RemoteAddrPort 远程地址端口
func (c *client) RemoteAddrPort() string { return strings.Split(c.conn.RemoteAddr().String(), ":")[1] }

// Write 写
func (c *client) Write(wp *WriteProto) { c.writeCh <- wp }

// JoinGroup 客户端加入组
func (c *client) JoinGroup(group string) (err error) {
	if group == "" {
		return errors.New("client: group name must be not empty")
	}

	if c.groups.Contains(group) {
		return errors.New("client: group was joined")
	}

	c.groups.Add(group)

	return nil
}

// LeaveGroup 客户端离开组
func (c *client) LeaveGroup(group string) (err error) {
	if group == "" {
		return errors.New("client: group name must be not empty")
	}

	if !c.groups.Contains(group) {
		return errors.New("client: group not joined")
	}

	c.groups.Delete(group)

	return
}

// LeaveGroups 客户端离开所有组
func (c *client) LeaveGroups() { c.groups.Clear() }

// InGroup 客户端是否在组
func (c *client) InGroup(group string) bool { return c.groups.Contains(group) }

//read 客户端读
func (c *client) read() {
	for {
		select {
		case <-c.stopReadCh:
			return
		default:
			rp := ReadProto{}
			if err := c.conn.ReadJSON(&rp); err != nil {
				c.close()
				return
			} else {
				rp.ClientID = c.ID()
				c.readCh <- &rp
			}
		}
	}
}

// write 客户端写
func (c *client) write() {
	for {
		select {
		case <-c.stopWriteCh:
			return
		case data, ok := <-c.writeCh:
			if ok {
				err := c.conn.WriteJSON(data)
				if err != nil {
					fmt.Println("client write", err)
				}
			} else {
				return
			}
		}
	}
}

// heart 心跳
func (c *client) heart() {
	for {
		select {
		case <-c.stopHeartCh:
			return
		case <-c.ticker.C:
			if try := c.heartMu.TryLock(); !try {
				continue
			}
			if err := c.conn.WriteControl(ws.PingMessage, []byte{}, time.Now().Add(time.Second)); err != nil {
				c.close()
				c.heartMu.Unlock()
				return
			}
			c.heartMu.Unlock()
		}
	}
}

// close 客户端关闭
func (c *client) close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return
	}
	c.closed = true
	c.closedTime = time.Now()
	_ = c.conn.Close()
	c.stopWriteCh <- struct{}{}
	c.stopHeartCh <- struct{}{}
	c.closeCh <- c.id
	c.ticker.Stop()
}

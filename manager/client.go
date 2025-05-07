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
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-the-way/wsserver/envs"
	"github.com/go-the-way/wsserver/logger"
	"github.com/go-the-way/wsserver/pkg"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type client struct {
	mu, hmu *sync.Mutex

	conn *websocket.Conn

	reqHeader http.Header

	id    string
	group string

	closed     bool
	closeCh    chan<- string
	readCh     chan map[string]any
	stopReadCh chan struct{}

	writeCh     chan any
	stopWriteCh chan struct{}

	stopHeartCh chan struct{}

	ht *time.Ticker

	connectedTime, closedTime time.Time
}

func newClient(conn *websocket.Conn, clientId, group string, reqHeader http.Header, closeCh chan<- string) *client {
	c := &client{
		mu:            &sync.Mutex{},
		hmu:           &sync.Mutex{},
		conn:          conn,
		reqHeader:     reqHeader,
		group:         group,
		closeCh:       closeCh,
		readCh:        make(chan map[string]any, 1),
		stopReadCh:    make(chan struct{}, 1),
		writeCh:       make(chan any, 1),
		stopWriteCh:   make(chan struct{}, 1),
		stopHeartCh:   make(chan struct{}, 1),
		ht:            time.NewTicker(time.Second * 1),
		connectedTime: time.Now(),
	}
	if c.isNode() {
		c.id = clientId
	} else {
		if envs.Debug && clientId != "" {
			c.id = clientId
		} else {
			c.id = pkg.Md5(uuid.NewString())
		}
	}
	go c.startRead()
	go c.startWrite()
	go c.startHeart()
	return c.onCreatedHandler()
}

func (c *client) remoteAddrIP() string {
	if xRealIp := c.reqHeader.Get("X-Real-Ip"); xRealIp != "" {
		return xRealIp
	}
	return strings.Split(c.conn.RemoteAddr().String(), ":")[0]
}

func (c *client) write(data any) { c.writeCh <- data }

func (c *client) isNode() bool   { return c.inGroup("node") }
func (c *client) isClient() bool { return !c.isNode() }

func (c *client) inGroup(group string) bool { return c.group == group }

func (c *client) startRead() {
	for {
		select {
		case <-c.stopReadCh:
			return
		case data := <-c.readCh:
			c.onReadHandler(data)
		default:
			mm := map[string]any{}
			if err := c.conn.ReadJSON(&mm); err != nil {
				c.debug("read json error: %v", err)
				c.close()
				return
			}
			c.readCh <- mm
		}
	}
}

func (c *client) startWrite() {
	for {
		select {
		case <-c.stopWriteCh:
			return
		case data := <-c.writeCh:
			if err := c.conn.WriteJSON(data); err != nil {
				c.debug("write json error: %v", err)
			}
		}
	}
}

func (c *client) startHeart() {
	defer c.ht.Stop()
	for {
		select {
		case <-c.stopHeartCh:
			return
		case <-c.ht.C:
			if try := c.hmu.TryLock(); !try {
				continue
			}
			if err := c.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second)); err != nil {
				c.close()
				c.hmu.Unlock()
				c.debug("write control message error: %v", err)
				return
			}
			c.hmu.Unlock()
		}
	}
}

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
	c.ht.Stop()
	c.closeCh <- c.id
	c.onClosedHandler()
}

func (c *client) debug(format string, v ...any) {
	logger.Debug(fmt.Sprintf("client["+c.id+"]%s", format), v...)
}

func (c *client) log(format string, v ...any) {
	logger.Log(fmt.Sprintf("client["+c.id+"]%s", format), v...)
}

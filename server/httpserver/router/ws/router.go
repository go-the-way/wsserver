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

package ws

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-the-way/svc"
	"github.com/go-the-way/wsserver/logger"
	"github.com/go-the-way/wsserver/manager"
	"github.com/go-the-way/wsserver/server/httpserver/app"
	"github.com/gorilla/websocket"
)

func init() {
	a := app.GetAppWithGroup("/api")
	a.GET("/ws", ws)
}

var ugr = websocket.Upgrader{
	HandshakeTimeout:  time.Second * 5,
	CheckOrigin:       func(r *http.Request) bool { return true },
	EnableCompression: true,
}

func ws(ctx *gin.Context) {
	svc.Query(ctx, func() (err error) {
		conn, err := ugr.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			logger.Log("websocket upgrade error: %v", err)
			return
		}
		reqHeader := ctx.Request.Header.Clone()
		group := reqHeader.Get("group")
		clientId := reqHeader.Get("client_id")
		manager.AddClient(conn, clientId, group, reqHeader)
		return svc.ErrNoReturn
	})
}

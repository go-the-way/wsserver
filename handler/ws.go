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

package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-the-way/wsserver/config"
	"github.com/go-the-way/wsserver/manager"

	"github.com/gorilla/websocket"
)

func init() {
	server := config.GetServer()
	server.HandleFunc("/ws", wsHandler)
}

// wsHandler ws处理
func wsHandler(w http.ResponseWriter, r *http.Request) {
	ugr := websocket.Upgrader{
		HandshakeTimeout:  time.Second * 5,
		CheckOrigin:       func(r *http.Request) bool { return true },
		EnableCompression: true,
	}
	if conn, err := ugr.Upgrade(w, r, nil); err != nil {
		fmt.Println("websocket upgrade", err)
	} else {
		h := r.Header.Clone()
		manager.Connect(conn, h.Get("group"))
	}
}

func writeJSON(w http.ResponseWriter, data any) {
	if buffer, err := json.Marshal(data); err != nil {
		writeError(w, err)
	} else {
		writeJSONBuffer(w, buffer)
	}
}

func writeJSONBuffer(w http.ResponseWriter, buffer []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buffer)
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(fmt.Sprintf("{\"error\":\"%v\"}\n", err)))
}

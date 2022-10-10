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

	"io/ioutil"
	"net/http"

	"github.com/go-the-way/wsserver/config"
	m "github.com/go-the-way/wsserver/manager"

	"github.com/go-the-way/streams"
	"github.com/go-the-way/streams/types"
)

func init() {
	server := config.GetServer()
	server.HandleFunc("/api/send_to_client", sendToClient)
	server.HandleFunc("/api/online_client", onlineClient)
}

type (
	msgProto struct {
		Type     string         `json:"type"`
		ClientID []string       `json:"client_ids,omitempty"`
		ClientId string         `json:"client_id,omitempty"`
		Data     map[string]any `json:"data"`
	}
	pRO = m.WriteProto
)

// send 发送消息
func sendToClient(w http.ResponseWriter, r *http.Request) {
	if readAll, err := ioutil.ReadAll(r.Body); err != nil {
		writeError(w, err)
	} else {
		proto := msgProto{}
		if err = json.Unmarshal(readAll, &proto); err != nil {
			writeError(w, err)
		} else {
			writeJSON(w, map[string]any{})
			go func(proto *msgProto) {
				set := types.MakeSet[string]()
				if cid := proto.ClientId; cid != "" {
					set.Add(cid)
				}
				if cid := proto.ClientID; cid != nil && len(cid) > 0 {
					streams.ForEach(cid, func(_ int, id string) { set.Add(id) })
				}
				set.Iterate(func(clientID string) { m.SendToClient(&pRO{Type: proto.Type, ClientID: clientID, Data: proto.Data}) })
			}(&proto)
		}
	}
}

func onlineClient(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]any{"clients": m.OnlineClient()})
}

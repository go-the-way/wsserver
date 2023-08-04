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

package sender

import (
	"context"

	"github.com/go-the-way/streams"
	"github.com/go-the-way/streams/ts"

	m "github.com/go-the-way/wsserver/manager"
)

type (
	Sender struct{}
	Args   struct {
		Type     string         `json:"type"`
		ClientID []string       `json:"client_id"`
		Data     map[string]any `json:"data"`
	}
	Reply struct {
		Code int `json:"code"`
	}
	pRO = m.WriteProto

	GArgs struct {
		Type  string         `json:"type"`
		Group string         `json:"group"`
		Data  map[string]any `json:"data"`
	}

	BCArgs struct {
		Type string         `json:"type"`
		Data map[string]any `json:"data"`
	}

	gPRO  = m.GWriteProto
	bcPRO = m.BCWriteProto
)

func (s *Sender) Send(_ context.Context, args Args, reply *Reply) error {
	set := ts.NewSet[string]()
	if cid := args.ClientID; cid != nil && len(cid) > 0 {
		streams.ForEach(cid, func(_ int, id string) { set.Add(id) })
	}
	set.Iterate(func(clientID string) { m.SendToClient(&pRO{Type: args.Type, ClientID: clientID, Data: args.Data}) })
	reply.Code = 200
	return nil
}

func (s *Sender) GSend(_ context.Context, args GArgs, reply *Reply) error {
	m.SendToGroup(&gPRO{Type: args.Type, Group: args.Group, Data: args.Data})
	reply.Code = 200
	return nil
}

func (s *Sender) Broadcast(_ context.Context, args BCArgs, reply *Reply) error {
	m.Broadcast(&bcPRO{Type: args.Type, Data: args.Data})
	reply.Code = 200
	return nil
}

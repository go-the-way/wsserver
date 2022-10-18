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

package service

import (
	"context"
	m "github.com/go-the-way/wsserver/manager"
)

type (
	GroupSender     struct{}
	GroupSenderArgs struct {
		Type  string         `json:"type"`
		Group string         `json:"group"`
		Data  map[string]any `json:"data"`
	}
	GroupSenderReply struct {
		Code int `json:"code"`
	}
	gPRO = m.GWriteProto
)

func (s *GroupSender) Send(_ context.Context, args GroupSenderArgs, reply *GroupSenderReply) error {
	m.SendToGroup(&gPRO{Type: args.Type, Group: args.Group, Data: args.Data})
	reply.Code = 200
	return nil
}

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

package client

import (
	"context"
	"github.com/go-the-way/wsserver/manager"
)

type (
	Client        struct{}
	JoinGroupArgs struct {
		ClientID string `json:"client_id"`
		Group    string `json:"group"`
	}
	LeaveGroupArgs    JoinGroupArgs
	LeaveAllGroupArgs struct {
		ClientID string `json:"client_id"`
	}
	Reply struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
)

func (s *Client) JoinGroup(_ context.Context, args JoinGroupArgs, reply *Reply) error {
	reply.Code = 200
	if err := manager.JoinGroup(args.ClientID, args.Group); err != nil {
		reply.Code = 500
		reply.Msg = err.Error()
	}
	return nil
}

func (s *Client) LeaveGroup(_ context.Context, args LeaveGroupArgs, reply *Reply) error {
	reply.Code = 200
	if err := manager.LeaveGroup(args.ClientID, args.Group); err != nil {
		reply.Code = 500
		reply.Msg = err.Error()
	}
	return nil
}

func (s *Client) LeaveAllGroup(_ context.Context, args LeaveAllGroupArgs, reply *Reply) error {
	reply.Code = 200
	if err := manager.LeaveAllGroup(args.ClientID); err != nil {
		reply.Code = 500
		reply.Msg = err.Error()
	}
	return nil
}

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

package sender

import (
	"context"
	"github.com/go-the-way/wsserver/server/rpcserver/svc"

	"github.com/go-the-way/wsserver/manager"
)

type Service struct{}

func NewService() *Service { return &Service{} }

func (s *Service) Send(_ context.Context, req SendReq, resp *Resp) (err error) {
	return svc.Do(req, resp, func(req SendReq) (err error) {
		for _, clientId := range req.ClientId {
			manager.Send(manager.SendTaskProtocol(req.Type, req.TaskId, req.TaskClientId, clientId, req.Data))
		}
		return
	}, func(err error, a *Resp) { a.Transform(err) })
}

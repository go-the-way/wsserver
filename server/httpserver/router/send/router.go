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

package send

import (
	"github.com/gin-gonic/gin"
	"github.com/go-the-way/svc"
	"github.com/go-the-way/wsserver/manager"
	"github.com/go-the-way/wsserver/server/httpserver/app"
)

type req struct {
	Type         string         `validate:"enum(ping|tcping|http|dns_resolve|mtr,任务类型不合法)" json:"type"`
	TaskId       string         `validate:"minlength(1,任务Id不能为空) maxlength(32,任务Id长度不能超过32)" json:"task_id"`
	TaskClientId string         `validate:"minlength(1,任务客户端Id不能为空) maxlength(32,任务客户端Id长度不能超过32)" json:"task_client_id"`
	ClientId     []string       `validate:"minlength(1,客户端Id不能为空) arr_minlength(1,客户端Id不能为空)" json:"client_id"`
	Data         map[string]any `json:"data"`
}

func init() {
	a := app.GetAppWithGroup("/api")
	a.POST("/send", send)
}

func send(ctx *gin.Context) {
	svc.BodyReq(ctx, req{}, func(req req) (err error) {
		for _, clientId := range req.ClientId {
			manager.Send(manager.SendTaskProtocol(req.Type, req.TaskId, req.TaskClientId, clientId, req.Data))
		}
		return
	})
}

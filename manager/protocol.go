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

type (
	SendProtocol interface {
		ClientID() (clientId string)
		SendData() (data any)
	}
	sendTaskProtocol struct {
		Type         string         `json:"type"`
		TaskId       string         `json:"task_id"`
		TaskClientId string         `json:"task_client_id"`
		ClientId     string         `json:"client_id"`
		Data         map[string]any `json:"data"`
	}
)

func SendTaskProtocol(Type string, taskId string, taskClientId string, clientId string, data map[string]any) *sendTaskProtocol {
	return &sendTaskProtocol{Type: Type, TaskId: taskId, TaskClientId: taskClientId, ClientId: clientId, Data: data}
}

func (p *sendTaskProtocol) ClientID() (clientId string) { return p.ClientId }
func (p *sendTaskProtocol) SendData() (data any)        { return p }

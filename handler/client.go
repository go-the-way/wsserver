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
	"net/http"

	"github.com/go-the-way/wsserver/config"

	m "github.com/go-the-way/wsserver/manager"
)

func init() {
	server := config.GetServer()
	server.HandleFunc("/api/online_client", onlineClient)
}

func onlineClient(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]any{"clients": m.OnlineClient()})
}

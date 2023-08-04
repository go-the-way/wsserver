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

package registry

import "time"

// Server the websocket server instance
type Server struct {
	ID           uint      // ID the server instance id
	Name         string    // Name the server instance name
	Endpoint     string    // Endpoint the server instance endpoint
	RegisterTime time.Time // RegisterTime the server instance register time
}

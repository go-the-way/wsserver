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

package app

import (
	"github.com/gin-gonic/gin"
	"github.com/go-the-way/svc"
)

var _app = svc.GetApp(gin.Recovery(), svc.Cors())

func GetApp() *gin.Engine { return _app }

func GetAppWithGroup(prefix string) gin.IRoutes { return GetApp().Group(prefix) }

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

package config

import (
	"os"
)

var (
	serverAddr = bindEnv("SERVER_ADDR", ":80")
	rpcAddr    = bindEnv("RPC_ADDR", ":86")
)

func GetServerAddr() string { return serverAddr }
func GetRpcAddr() string    { return rpcAddr }

func bindEnv(key, defVal string) string {
	if n := os.Getenv(key); n != "" {
		return n
	}
	return defVal
}

//
// func bindEnvInt(key string, defVal int) (i int) {
//	if n := bindEnv(key, ""); n == "" {
//		i = defVal
//	} else {
//		i, _ = strconv.Atoi(n)
//	}
//	return
// }

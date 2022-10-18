package rpc

import (
	"fmt"

	"github.com/go-the-way/wsserver/config"
	"github.com/go-the-way/wsserver/rpc/service"

	"github.com/smallnest/rpcx/server"
)

func init() { go serve() }

func serve() {
	s := server.NewServer()
	rpcAddr := config.GetRpcAddr()
	_ = s.Register(new(service.Sender), "")
	_ = s.Register(new(service.GroupSender), "")
	fmt.Println("rpc server started on", rpcAddr)
	fmt.Println(s.Serve("tcp", rpcAddr))
}

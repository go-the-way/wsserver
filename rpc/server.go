package rpc

import (
	"fmt"
	"github.com/go-the-way/wsserver/rpc/service/client"
	"github.com/go-the-way/wsserver/rpc/service/sender"

	"github.com/go-the-way/wsserver/config"
	"github.com/smallnest/rpcx/server"
)

func init() { go serve() }

func serve() {
	s := server.NewServer()
	rpcAddr := config.GetRpcAddr()
	_ = s.Register(new(sender.Sender), "")
	_ = s.Register(new(client.Client), "")
	fmt.Println("rpc server started on", rpcAddr)
	fmt.Println(s.Serve("tcp", rpcAddr))
}

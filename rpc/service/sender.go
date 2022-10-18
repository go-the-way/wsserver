package service

import (
	"context"

	"github.com/go-the-way/streams"
	"github.com/go-the-way/streams/ts"

	m "github.com/go-the-way/wsserver/manager"
)

type (
	Sender     struct{}
	SenderArgs struct {
		Type     string         `json:"type"`
		ClientID []string       `json:"client_id"`
		Data     map[string]any `json:"data"`
	}
	SenderReply struct {
		Code int `json:"code"`
	}
	pRO = m.WriteProto
)

func (s *Sender) Send(_ context.Context, args SenderArgs, reply *SenderReply) error {
	set := ts.NewSet[string]()
	if cid := args.ClientID; cid != nil && len(cid) > 0 {
		streams.ForEach(cid, func(_ int, id string) { set.Add(id) })
	}
	set.Iterate(func(clientID string) { m.SendToClient(&pRO{Type: args.Type, ClientID: clientID, Data: args.Data}) })
	reply.Code = 200
	return nil
}

package participant

import (
  "net/rpc"
  "net"
  "log"
)

type Participant struct {
  Objects map[string]Object
  Transactions map[int32]Transaction
  Address string
  Id int
}

var self Participant

func Start(hostname string, id int) error {
  log.Println("Starting participant")
  self = New(hostname, id)
  go self.setupRPC()
  return nil
}

func (p Participant) setupRPC() error {
  rpc.Register(&self)
  l, e := net.Listen("tcp", ":3000")
  if e != nil {
    log.Println("Error in setup RPC:", e)
    return e
  }
  go rpc.Accept(l)
  return nil
}

func New(addr string, id int) Participant {
  objs := make(map[string]Object, 0)
  trans := make(map[int32]Transaction, 0)
  return Participant{objs, trans, addr, id}
}

package participant

import (
	"log"
	"net"
	"net/rpc"
	"sync"
)

type Participant struct {
	Objects      map[string]*Object
	Transactions map[int32]*Transaction
	Address      string
	Id           int
}

var self Participant
var wg sync.WaitGroup

func Start(hostname string, id int) error {
	log.Println("Starting participant")
	self = New(hostname, id)
	go self.setupRPC()
	wg.Add(1)
	wg.Wait()
	return nil
}

func (p Participant) setupRPC() {
	part := new(Participant)
	rpc.Register(part)
	l, e := net.Listen("tcp", ":3000")
	if e != nil {
		log.Println("Error in setup RPC:", e)
	}
	log.Println("Setting up participant RPCs")
	go rpc.Accept(l)
}

func New(addr string, id int) Participant {
	objs := make(map[string]*Object, 0)
	trans := make(map[int32]*Transaction, 0)
	return Participant{objs, trans, addr, id}
}

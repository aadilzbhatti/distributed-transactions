package coordinator

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"node/participant"
	"sync"
)

var host string = "sp17-cs425-g26-0%d.cs.illinois.edu"
var mutex = &sync.Mutex{}
var self Coordinator
var graph *Graph

type Coordinator struct {
	Participants map[string]participant.Participant
}

func Start() error {
	log.Println("Starting coordinator..")
	self = New()

	// set up RPCs
	e := self.setupRPC()
	if e != nil {
		return e
	}

	// join up with participant servers
	for i := 2; i < 10; i++ {
		go self.joinParticipant(i)
	}

  // set up deadlock detection graph
	graph = NewGraph()

	// interface with client
	return nil
}

func New() Coordinator {
	parts := make(map[string]participant.Participant, 0)
	c := Coordinator{parts}
	return c
}

func (c Coordinator) setupRPC() error {
	coord := new(Coordinator)
	rpc.Register(coord)
	l, e := net.Listen("tcp", ":3000")
	if e != nil {
		log.Println("Error in setup RPC:", e)
		return e
	}
	go rpc.Accept(l)
	return nil
}

func (c Coordinator) joinParticipant(id int) {
  serverId := string(rune('A' + (id - 2)))
	log.Printf("Trying to join node %v\n", serverId)
	hostname := fmt.Sprintf("%s:%d", fmt.Sprintf(host, id), 3000)

	for {
		client, err := rpc.Dial("tcp", hostname)
		if err != nil {
			continue

		} else {
			var reply participant.Participant
			ja := participant.JoinArgs{}
			err = client.Call("Participant.Join", &ja, &reply)

			if err != nil {
				log.Println("Error in join: ", err)

			} else {
				mutex.Lock()
				c.Participants[serverId] = reply
				mutex.Unlock()
				log.Printf("Server %v joined the system\n", serverId)
			}

      graph.AddVertex(serverId)
			client.Close()
			return
		}
	}
}

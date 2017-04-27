package coordinator

import (
  "net/rpc"
  "net"
  "log"
  "fmt"
  "sync"
  "node/participant"
)

var host string = "sp17-cs425-g26-0%d.cs.illinois.edu"
var mutex = &sync.Mutex{}
var self Coordinator

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
  log.Printf("Trying to join node %v\n", id)
  hostname := fmt.Sprintf("%s:%d", fmt.Sprintf(host, id), 3000)
  for {
    client, err := rpc.Dial("tcp", hostname)
    if err != nil {
      log.Println(err)
      continue

    } else {
      var reply participant.Participant
      err = client.Call("Participant.Join", nil, &reply)
      log.Println("Did da join")

      if err != nil {
        log.Println("Error in join: ", err)

      } else {
        serverId := string(rune('A' + (id - 2)))

        mutex.Lock()
        c.Participants[serverId] = reply
        mutex.Unlock()
        log.Printf("Server %v joined the system\n", serverId)
      }
      client.Close()
      return
    }
  }
}

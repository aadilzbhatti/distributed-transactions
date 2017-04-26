package coordinator

import (
  "net/rpc"
  "log"
  "fmt"
  "sync"
  "node/participant"
)

var host string = "sp17-cs425-g26-0%d.cs.illinois.edu"
var mutex = &sync.Mutex{}

type Coordinator struct {
  participants map[string]Participant
}

func Start() error {
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
}

func New() Coordinator {
  parts := make(map[string]Participant, 0)
  c := Coordinator{parts}
  return c
}

func (c Coordinator) setupRPC() error {
  rpc.Register(&self)
  l, e := net.Listen("tcp", ":3000")
  if e != nil {
    log.Println("Error in setup RPC:", e)
    return e
  }
  go rpc.Accept(l)
  return nil
}

func (c Coordinator) joinParticipant(id int) {
  hostname := fmt.Sprintf("%s:%d", fmt.Sprintf(host, id), 4000)
  for {
    client, err := rpc.Dial("tcp", hostname)
    if err != nil {
      continue

    } else {
      defer client.Close()
      var reply Participant
      err = client.Call("Participant.Join", nil, &reply)

      if err != nil {
        log.Println("Error in join: ", err)

      } else {
        serverId := string(rune('A' + (id - 2)))

        mutex.Lock()
        c.participants[serverId] = reply
        mutex.Unlock()
        log.Printf("Server %v joined the system\n", serverId)
      }
      return
    }
  }
}

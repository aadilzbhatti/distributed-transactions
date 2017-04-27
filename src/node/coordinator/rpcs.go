package coordinator

import (
  "fmt"
  "time"
  "net/rpc"
  "log"
  "node/participant"
)

type BeginArgs struct {}

type SetArgs struct {
  Tid int32
  ServerId string
  Key string
  Value string
}

type CoordGetArgs struct {
  Tid int32
  ServerId string
  Key string
}

type AbortArgs struct {
  Tid int32
}

type CommitArgs struct {
  Tid int32
}

func (c Coordinator) Begin(ba *BeginArgs, reply *int32) error {
  log.Println("In Begin!")
  *reply = int32(time.Now().Unix())
  return nil
}

func (c Coordinator) Set(sa *SetArgs, reply *bool) error {
  log.Printf("%v\n", self.Participants)
  if p, ok := self.Participants[sa.ServerId]; ok {
    log.Println(p.Address)
    client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", p.Address, 3000))
    if err != nil {
      log.Println("Error in Set/Dial: ", err)
      return err
    }

    psa := participant.SetArgs{sa.Tid, sa.Key, sa.Value}
    err = client.Call("Participant.SetKey", &psa, &reply)
    if err != nil {
      log.Println("Error in Set/RPC: ", err)
      return err
    }
    return nil

  } else {
    return fmt.Errorf("No such server in system")
  }
}

func (c Coordinator) Get(ga *CoordGetArgs, reply *bool) error {
  if p, ok := self.Participants[ga.ServerId]; ok {
    client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", p.Address, 3000))
    if err != nil {
      log.Println("Error in Get/Dial: ", err)
      return err
    }

    pga := participant.GetArgs{ga.Tid, ga.Key}
    err = client.Call("Participant.GetKey", &pga, &reply)
    if err != nil {
      log.Println("Error in Get/RPC: ", err)
      return err
    }
    return nil

  } else {
    return fmt.Errorf("No such server in system")
  }
}

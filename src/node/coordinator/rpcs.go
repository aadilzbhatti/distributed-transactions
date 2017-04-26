package coordinator

import (
  "fmt"
  "time"
  "node/participant"
)

type BeginArgs struct {}

type CoordSetArgs struct {
  serverId string
  key string
  value string
}

type CoordGetArgs struct {
  serverId string
  key string
}

type AbortArgs struct {
  tid int32
}

type CommitArgs struct {
  tid int32
}

func (c Coordinator) Begin(ba *BeginArgs, reply *int32) error {
  *reply = int32(time.Now().Unix())
  return nil
}

func (c Coordinator) Set(sa *CoordSetArgs, reply *bool) error {
  if p, ok := c.participants[sa.serverId]; ok {
    client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", p.Address, 3000))
    if err != nil {
      log.Println("Error in Set/Dial: ", err)
      return err
    }

    psa := participant.SetArgs{sa.key, sa.value}
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
  if p, ok := c.participants[ga.serverId]; ok {
    client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", p.Address, 3000))
    if err != nil {
      log.Println("Error in Get/Dial: ", err)
      return err
    }

    pga := participant.GetArgs{ga.key}
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

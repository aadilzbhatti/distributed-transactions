package dtrans

import (
  "node/coordinator"
  "net/rpc"
  "log"
)

var chost string = "sp17-cs425-g26-01.cs.illinois.edu:3000"
var host string = "sp17-cs425-g26-0%d.cs.illinois.edu:3000"

func Begin() (error, int32) {
  client, err := rpc.Dial("tcp", chost)
  log.Println("Beginning...")
  if err != nil {
    log.Println("Error in Begin/Dial: ", err)
    return err, 0
  }

  var reply int32
  ba := coordinator.BeginArgs{}
  err = client.Call("Coordinator.Begin", &ba, &reply)
  return err, reply
}

func Set(serverId string, key string, value string, currId int32) error {
  client, err := rpc.Dial("tcp", chost)
  if err != nil {
    log.Println("Error in Set/Dial: ", err)
    return err
  }

  sa := coordinator.SetArgs{currId, serverId, key, value}
  var reply bool
  err = client.Call("Coordinator.Set", &sa, &reply)
  if err != nil {
    log.Println("Error in Set/RPC: ", err)
    return err
  }

  return nil
}

func Get(serverId string, key string, currId int32) string {
  // TODO RPC to Coordinator.Get
  return ""
}

func Abort() {
  // TODO RPC to Coordinator.Abort
}

func Commit() {
  // TODO RPC to Coordinator.Commit
}

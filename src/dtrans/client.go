package dtrans

import (
  "node/coordinator"
)

var coordinator string = "sp17-cs425-g26-01.cs.illinois.edu:3000"

func Begin() (error, int32) {
  client, err := rpc.Dial("tcp", coordintor)
  if err != nil {
    log.Println("Error in Begin/Dial: ", err)
    return (err, 0)
  }

  var reply int32
  err = client.Call("Coordinator.Begin", nil, &reply)
  return (err, reply)
}

func Set(string serverId, string key, string value) {
  // TODO RPC to Coordinator.Set
}

func Get(string serverId, string key) string {
  // TODO RPC to Coordinator.Get
}

func Abort() {
  // TODO RPC to Coordinator.Abort
}

func Commit() {
  // TODO RPC to Coordinator.Commit
}

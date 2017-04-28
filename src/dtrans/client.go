package dtrans

import (
	"log"
	"net/rpc"
	"node/coordinator"
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

	log.Printf("Calling on server %v\n", serverId)
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
	client, err := rpc.Dial("tcp", chost)
	if err != nil {
		log.Println("Error in Get/Dial: ", err)
		return ""
	}

	ga := coordinator.GetArgs{currId, serverId, key}
	var reply string
	err = client.Call("Coordinator.Get", &ga, &reply)
	if err != nil {
		log.Println("Error in Get/RPC: ", err)
		return ""
	}

	return reply
}

func Abort() error {
	client, err := rpc.Dial("tcp", chost)
	if err != nil {
		log.Println("Error in Abort/Dial:", err)
		return err
	}

	aa := coordinator.AbortArgs{currentId}
	var reply bool
	err = client.Call("Coordinator.Abort", &aa, &reply)
	if err != nil {
		log.Println("Error in Abort/RPC:", err)
		return err
	}

	return nil
}

func Commit() error {
	client, err := rpc.Dial("tcp", chost)
	if err != nil {
		log.Println("Error in Commit/Dial:", err)
		return err
	}

	ca := coordinator.CommitArgs{currentId}
	var reply bool
	err = client.Call("Coordinator.Commit", &ca, &reply)
	if err != nil {
		log.Println("Error in Commit/RPC:", err)
		return err
	}

  if !reply {
    return Abort()
  }
  
	return nil
}

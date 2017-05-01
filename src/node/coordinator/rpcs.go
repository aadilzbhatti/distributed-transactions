package coordinator

import (
	"fmt"
	"log"
	"net/rpc"
	"node/participant"
	"time"
)

type BeginArgs struct{}

type SetArgs struct {
	Tid      int32
  MyId string
	ServerId string
	Key      string
	Value    string
}

type GetArgs struct {
	Tid      int32
	MyId string
	ServerId string
	Key      string
}

type AbortArgs struct {
	Tid int32
}

type CommitArgs struct {
	Tid int32
}

func (c Coordinator) Begin(ba *BeginArgs, reply *int32) error {
	for _, s := range self.Participants {
		client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", s.Address, 3000))
		if err != nil {
			log.Println("Error in Begin/Dial:", err)
			return err
		}

		pba := participant.BeginArgs{}
		var r bool
		err = client.Call("Participant.Begin", &pba, &r)
		if err != nil {
			log.Println("Error in Begin/RPC:", err)
			return err
		}
		client.Close()
	}
	*reply = int32(time.Now().Unix())
	return nil
}

func (c Coordinator) Set(sa *SetArgs, reply *bool) error {
	otherId := string([]rune(sa.ServerId)[0] - 63)
	fmt.Println(otherId, []rune(sa.ServerId)[0])

	if p, ok := self.Participants[sa.ServerId]; ok {
		client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", p.Address, 3000))
		defer client.Close()

		if err != nil {
			log.Println("Error in Set/Dial: ", err)
			return err
		}

		// add new edge to Graph
		graph.AddVertex(sa.MyId)
		fmt.Println("It was called herEEE!")
		graph.AddEdge(sa.MyId, otherId, sa.Tid)
		fmt.Println("fuck me")

		// if cycle in Graph caused by this transaction
		if graph.DetectCycle(sa.Tid) {
			graph.RemoveEdge(sa.Tid)
			fmt.Println("About to abort")

			// abort this Transaction
			aa := AbortArgs{sa.Tid}
			var r bool
			c.Abort(&aa, &r)
			return fmt.Errorf("Transaction caused deadlock, aborted")
		}

		// otherwise continue
		psa := participant.SetArgs{sa.Tid, sa.Key, sa.Value}
		err = client.Call("Participant.SetKey", &psa, &reply)
		if err != nil {
			log.Println("Error in Set/RPC: ", err)
			return err
		}

		graph.RemoveEdge(sa.Tid)
		return nil

	} else {
		return fmt.Errorf("No such server in system\n")
	}
}

func (c Coordinator) Get(ga *GetArgs, reply *string) error {
	otherId := string([]rune(ga.ServerId)[0] - 63)

	if p, ok := self.Participants[ga.ServerId]; ok {
		client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", p.Address, 3000))
		defer client.Close()

		if err != nil {
			log.Println("Error in Get/Dial: ", err)
			return err
		}

		// add new edge to Graph
		graph.AddVertex(ga.MyId)
		graph.AddEdge(ga.MyId, otherId, ga.Tid)

		// if cycle in Graph caused by this transaction
		if graph.DetectCycle(ga.Tid) {
			graph.RemoveEdge(ga.Tid)

			// abort this Transaction
			aa := AbortArgs{ga.Tid}
			var r bool
			c.Abort(&aa, &r)
			return fmt.Errorf("Transaction caused deadlock, aborted")
		}

		// abort transaction
		pga := participant.GetArgs{ga.Tid, ga.Key}
		err = client.Call("Participant.GetKey", &pga, &reply)
		if err != nil {
			return err
		}

		// remove created edge
		graph.RemoveEdge(ga.Tid)
		return nil

	} else {
		return fmt.Errorf("No such server in system\n")
	}
}

func (c Coordinator) Commit(ca *CommitArgs, reply *bool) error {
  fmt.Println("COMMIT TIME!!")
	// check if we can commit
	cca := participant.CanCommitArgs{ca.Tid}
	for _, p := range self.Participants {
    fmt.Println("CHECKING WITH", p)
		client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", p.Address, 3000))
		if err != nil {
			log.Println("Error in Commit/Dial:", err)
			client.Close()
			return err
		}

		var check bool
		err = client.Call("Participant.CanCommit", &cca, &check)
		if err != nil {
			if err.Error() != "No such transaction in server" {
				log.Println("Error in Commit/RPC:", err)
				client.Close()
				return err
			} else {
				continue
			}
		}

		if !check {
			*reply = false
			log.Println("Someone said no!")
			client.Close()
			return nil
		}
		client.Close()
	}

	// if we can, we commit
	dca := participant.DoCommitArgs{ca.Tid}
	for _, p := range self.Participants {
    fmt.Println("COMMITTING TO", p)
		client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", p.Address, 3000))
		if err != nil {
			log.Println("Error in DoCommit/Dial:", err)
			client.Close()
			return err
		}

		var check bool
		err = client.Call("Participant.DoCommit", &dca, &check)
		if err != nil && err.Error() != "No such transaction in server" {
			log.Println("Error in DoCommit/RPC:", err)
			client.Close()
			return err
		}

		client.Close()
	}

	*reply = true
	return nil
}

func (c Coordinator) Abort(aa *AbortArgs, reply *bool) error {
	paa := participant.DoAbortArgs{aa.Tid}
	for _, p := range self.Participants {
		log.Println(p)
		client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", p.Address, 3000))
		if err != nil {
			log.Println("Error in Abort/Dial:", err)
			return err
		}

		var check bool
		err = client.Call("Participant.DoAbort", &paa, &check)
		if err != nil && err.Error() != "No such transaction in server" {
			log.Println("Error in DoAbort/RPC:", err)
			client.Close()
			return err
		}

		client.Close()
	}

	return nil
}

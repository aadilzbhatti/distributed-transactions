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
	ServerId string
	Key      string
	Value    string
}

type GetArgs struct {
	Tid      int32
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
	log.Println("In Begin!")
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
	log.Printf("%v\n", self.Participants)
	if p, ok := self.Participants[sa.ServerId]; ok {
		client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", p.Address, 3000))
		defer client.Close()

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

func (c Coordinator) Get(ga *GetArgs, reply *string) error {
	if p, ok := self.Participants[ga.ServerId]; ok {
		client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", p.Address, 3000))
		defer client.Close()

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

func (c Coordinator) Commit(ca *CommitArgs, reply *bool) error {
	// check if we can commit
	cca := participant.CanCommitArgs{ca.Tid}
	for _, p := range self.Participants {
		client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", p.Address, 3000))
		if err != nil {
			log.Println("Error in Commit/Dial:", err)
			client.Close()
			return err
		}

		var check bool
		err = client.Call("Participant.CanCommit", &cca, &check)
		if err != nil && err.Error() != "No such transaction in server" {
			log.Println("Error in Commit/RPC:", err)
			client.Close()
			return err
		}

		if !check {
			*reply = false
			client.Close()
			return nil
		}
		client.Close()
	}

	// if we can, we commit
	dca := participant.DoCommitArgs{ca.Tid}
	for _, p := range self.Participants {
		client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", p.Address, 3000))
		if err != nil {
			log.Println("Error in DoCommit/Dial:", err)
			client.Close()
			return err
		}

		var check bool
		err = client.Call("Participant.DoCommit", &dca, &check)
		if err != nil {
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
		client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", p.Address, 3000))
		if err != nil {
			log.Println("Error in Abort/Dial:", err)
			return err
		}

		var check bool
		err = client.Call("Participant.DoAbort", &paa, &check)
		if err != nil {
			log.Println("Error in DoAbort/RPC:", err)
			client.Close()
			return err
		}

		client.Close()
	}

	return nil
}

package participant

import (
	"fmt"
	"log"
	"sync"
)

var mutex = &sync.Mutex{}

type CanCommitArgs struct {
	Tid int32
}

type DoCommitArgs struct {
	Tid int32
}

type DoAbortArgs struct {
	Tid int32
}

type JoinArgs struct{}

type SetArgs struct {
	Tid   int32
	Key   string
	Value string
}

type GetArgs struct {
	Tid int32
	Key string
}

type BeginArgs struct{}

func (p *Participant) Join(ja *JoinArgs, reply *Participant) error {
	*reply = self
	return nil
}

func (p *Participant) Begin(ba *BeginArgs, reply *bool) error {
	for k := range self.Objects {
		self.Objects[k].start()
	}

	*reply = true
	log.Println("Initialized all objects for transaction")
	return nil
}

func (p *Participant) CanCommit(cca *CanCommitArgs, reply *bool) error {
	log.Println(self.Transactions, cca.Tid)
	if value, ok := self.Transactions[cca.Tid]; ok {
		log.Println("In here!")
		*reply = !value.hasFailed()
		return nil
	} 
	log.Println("Should not get here..")
	return fmt.Errorf("No such transaction in server")
}

func (p *Participant) DoCommit(dca *DoCommitArgs, reply *bool) error {
	if value, ok := self.Transactions[dca.Tid]; ok {
		for k := range self.Objects {
			self.Objects[k].stop()
		}
		value.commit()
		*reply = true
		return nil
	}
	return fmt.Errorf("No such transaction in server")
}

func (p *Participant) DoAbort(daa *DoAbortArgs, reply *bool) error {
	if value, ok := self.Transactions[daa.Tid]; ok {
		for k := range self.Objects {
			self.Objects[k].stop()
		}
		value.abort()
		*reply = true
		return nil
	}
	log.Println("No such transaction in server")
	return nil
}

func (p *Participant) SetKey(sa *SetArgs, reply *bool) error {
	log.Printf("In set!: %v\n", sa)
	if trans, ok := self.Transactions[sa.Tid]; ok {
		// we are executing a running transaction
		log.Println(trans)

	} else {
		log.Println("In here?")
		// we need to start a new transaction
		t := Transaction{sa.Tid, false}
		self.Transactions[sa.Tid] = &t
	}
	if _, ok := self.Objects[sa.Key]; ok {
		self.Objects[sa.Key].setKey(sa.Value, sa.Tid)
		log.Printf("Just reset %v to %v=%v\n", sa.Key, sa.Key, self.Objects[sa.Key])
	} else {
		mutex.Lock()
		self.Objects[sa.Key] = NewObject(sa.Key, sa.Value, sa.Tid)
		mutex.Unlock()
	}
	*reply = true
	log.Printf("Finished setting %v = %v\n", sa.Key, sa.Value)
	log.Println(self.Objects[sa.Key])
	return nil
}

func (p *Participant) GetKey(ga *GetArgs, reply *string) error {
	if _, ok := self.Transactions[ga.Tid]; ok {
		// we are executing a running Transaction

	} else {
		// we need to start a new transaction
		t := Transaction{ga.Tid, false}
		self.Transactions[ga.Tid] = &t
	}
	if v, ok := self.Objects[ga.Key]; ok {
		*reply = v.getKey(ga.Tid)
	} else {
		reply = nil
		return fmt.Errorf("No such object in server")
	}
	return nil
}

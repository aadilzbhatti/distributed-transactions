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
	mutex.Lock()
	defer mutex.Unlock()

	for k := range self.Objects {
		self.Objects[k].start()
	}

	*reply = true
	log.Println("Initialized all objects for transaction")
	return nil
}

func (p *Participant) CanCommit(cca *CanCommitArgs, reply *bool) error {
	mutex.Lock()
	defer mutex.Unlock()

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
	mutex.Lock()
	defer mutex.Unlock()

	if value, ok := self.Transactions[dca.Tid]; ok {
		for k, v := range self.Transactions[dca.Tid].updates {
			self.Objects[k] = v
		}
		for k, _ := range self.Objects {
			self.Objects[k].stop()
			self.held[k].holding = false
			self.held[k].cond.Broadcast()
		}
		value.commit()
		*reply = true
		return nil
	}
	return fmt.Errorf("No such transaction in server")
}

func (p *Participant) DoAbort(daa *DoAbortArgs, reply *bool) error {
	mutex.Lock()
	defer mutex.Unlock()

	if trans, ok := self.Transactions[daa.Tid]; ok {
		for k := range self.Objects {
			self.Objects[k].stop()
      self.Objects[k].resetKey(trans.initial[k].Value, daa.Tid)
			self.held[k].holding = false
			self.held[k].cond.Broadcast()
		}
		trans.abort()
		*reply = true
		return nil
	}
	return fmt.Errorf("No such transaction in server")
}

func (p *Participant) SetKey(sa *SetArgs, reply *bool) error {
	mutex.Lock()
	defer mutex.Unlock()

	log.Printf("In set: %v\n", sa)
	if trans, ok := self.Transactions[sa.Tid]; ok {
		// we are executing a running transaction
		log.Println(trans)
	} else {
		// we need to start a new transaction
		self.Transactions[sa.Tid] = NewTransaction(sa.Tid)

    // set initial state of transaction
    for k, _ := range self.Objects {
			self.Transactions[sa.Tid].addObject(k, *self.Objects[k])
		}
	}

	if _, ok := self.Transactions[sa.Tid].updates[sa.Key]; ok {
		self.Transactions[sa.Tid].updateObject(sa.Key, sa.Value)
		log.Printf("Reseting %v to %v=%v\n", sa.Key, sa.Key, self.Objects[sa.Key])
		self.held[sa.Key] = NewHeld(sa.Key, sa.Tid)

	} else {
		obj := NewObject(sa.Key, sa.Value, sa.Tid)
		self.Transactions[sa.Tid].updates[sa.Key] = obj
		h := NewHeld(sa.Key, sa.Tid)
		self.held[sa.Key] = h
	}

	*reply = true
	log.Printf("Finished setting %v = %v\n", sa.Key, sa.Value)
	log.Println(self.Transactions[sa.Tid].updates[sa.Key])
	return nil
}

func (p *Participant) GetKey(ga *GetArgs, reply *string) error {
	mutex.Lock()
	defer mutex.Unlock()

	if trans, ok := self.Transactions[ga.Tid]; ok {
		// we are executing a running Transaction
    log.Println(trans)
	} else {
		// we need to start a new transaction
		self.Transactions[ga.Tid] = NewTransaction(ga.Tid)

		// set initial state of transaction
		for k, v := range self.Objects {
			if !v.running {
				self.Transactions[ga.Tid].addObject(k, *self.Objects[k])
			}
		}
	}

	if v, ok := self.Transactions[ga.Tid].updates[ga.Key]; ok {
		if v.currTrans != ga.Tid && v.currTrans != 0 {
			log.Println(v)
			*reply = self.Transactions[ga.Tid].initial[ga.Key].Value
			log.Println("Ah!", *reply)
			return nil
		} else if _, ok2 := self.Transactions[ga.Tid].updates[ga.Key]; !ok2 {
			log.Println(v)
			*reply = self.Transactions[ga.Tid].initial[ga.Key].Value
			log.Println("Achoo!", *reply)
			return nil
		}
		*reply = self.Transactions[ga.Tid].updates[ga.Key].getKey()
		log.Println("o", *reply)

	} else if v, ok := self.Objects[ga.Key]; ok {
		if v.currTrans != ga.Tid && v.currTrans != 0 {
			log.Println(v)
			*reply = self.Transactions[ga.Tid].initial[ga.Key].Value
			log.Println("Ah!", *reply)
			return nil
		} else if _, ok2 := self.Transactions[ga.Tid].updates[ga.Key]; !ok2 {
			log.Println(v)
			*reply = v.getKey()
			log.Println("Achoo!", *reply)
			return nil
		}
		*reply = self.Transactions[ga.Tid].updates[ga.Key].getKey()
		log.Println("o", *reply)
	} else {
		*reply = "NOT FOUND"
		return fmt.Errorf("No such object in server")
	}

	return nil
}

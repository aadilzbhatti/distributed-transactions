package participant

import (
  "fmt"
  "sync"
  "log"
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

type JoinArgs struct {}

type SetArgs struct {
  Tid int32
  Key string
  Value string
}

type GetArgs struct {
  Tid int32
  Key string
}

func (p *Participant) Join(ja *JoinArgs, reply *Participant) error {
  log.Println("In join!")
  *reply = self
  return nil
}

func (p *Participant) CanCommit(cca *CanCommitArgs, reply *bool) error {
  if value, ok := self.Transactions[cca.Tid]; ok {
    *reply = !value.Failed
    return nil
  }
  return fmt.Errorf("No such transaction in server")
}

func (p *Participant) DoCommit(dca *DoCommitArgs, reply *bool) error {
  if value, ok := self.Transactions[dca.Tid]; ok {
    value.Commit()
    *reply = true
    return nil
  }
  return fmt.Errorf("No such transaction in server")
}

func (p *Participant) DoAbort(daa *DoAbortArgs, reply *bool) error {
  if value, ok := self.Transactions[daa.Tid]; ok {
    value.Abort()
    *reply = true
    return nil
  }
  return fmt.Errorf("No such transaction in server")
}

func (p *Participant) SetKey(sa *SetArgs, reply *bool) error {
  log.Printf("In set!: %v\n", sa)
  if trans, ok := self.Transactions[sa.Tid]; ok {
    // we are executing a running transaction
    log.Println(trans)

  } else {
    // we need to start a new transaction
    t := Transaction{sa.Tid, false}
    self.Transactions[sa.Tid] = t
  }
  if _, ok := self.Objects[sa.Key]; ok {
    self.Objects[sa.Key].setKey(sa.Value)
  } else {
    mutex.Lock()
    self.Objects[sa.Key] = NewObject(sa.Key, sa.Value)
    mutex.Unlock()
  }
  *reply = true
  return nil
}

func (p *Participant) GetKey(ga *GetArgs, reply *string) error {
  if _, ok := self.Transactions[ga.Tid]; ok {
    // we are executing a running Transaction

  } else {
    // we need to start a new transaction
    t := Transaction{ga.Tid, false}
    self.Transactions[ga.Tid] = t
  }
  if v, ok := self.Objects[ga.Key]; ok {
    *reply = v.getKey()
  } else {
    reply = nil
    return fmt.Errorf("No such object in server")
  }
  return nil
}

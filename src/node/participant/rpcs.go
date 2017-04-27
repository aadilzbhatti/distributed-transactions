package participant

import (
  "fmt"
  "sync"
  "log"
)

var mutex = &sync.Mutex{}

type CanCommitArgs struct {
  tid int32
}

type DoCommitArgs struct {
  tid int32
}

type DoAbortArgs struct {
  tid int32
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
  *reply = *p
  return nil
}

func (p *Participant) CanCommit(cca *CanCommitArgs, reply *bool) error {
  if value, ok := p.Transactions[cca.tid]; ok {
    *reply = !value.failed
    return nil
  }
  return fmt.Errorf("No such transaction in server")
}

func (p *Participant) DoCommit(dca *DoCommitArgs, reply *bool) error {
  if value, ok := p.Transactions[dca.tid]; ok {
    value.Commit()
    *reply = true
    return nil
  }
  return fmt.Errorf("No such transaction in server")
}

func (p *Participant) DoAbort(daa *DoAbortArgs, reply *bool) error {
  if value, ok := p.Transactions[daa.tid]; ok {
    value.Abort()
    *reply = true
    return nil
  }
  return fmt.Errorf("No such transaction in server")
}

func (p *Participant) SetKey(sa *SetArgs, reply *bool) error {
  if trans, ok := p.Transactions[sa.Tid]; ok {
    // we are executing a running transaction
    log.Println(trans)

  } else {
    // we need to start a new transaction
    t := Transaction{sa.Tid, false}
    p.Transactions[sa.Tid] = t
  }
  if _, ok := p.Objects[sa.Key]; ok {
    p.Objects[sa.Key].setKey(sa.Value)
  } else {
    mutex.Lock()
    p.Objects[sa.Key] = NewObject(sa.Key, sa.Value)
    mutex.Unlock()
  }
  *reply = true
  return nil
}

func (p *Participant) GetKey(ga *GetArgs, reply *string) error {
  if _, ok := p.Transactions[ga.Tid]; ok {
    // we are executing a running Transaction

  } else {
    // we need to start a new transaction
    t := Transaction{ga.Tid, false}
    p.Transactions[ga.Tid] = t
  }
  if v, ok := p.Objects[ga.Key]; ok {
    *reply = v.getKey()
  } else {
    reply = nil
    return fmt.Errorf("No such object in server")
  }
  return nil
}

package participant

type CanCommitArgs struct {
  tid int
}

type DoCommitArgs struct {
  tid int
}

type DoAbortArgs struct {
  tid int
}

type JoinArgs struct {}

type SetArgs struct {
  key string
  value string
}

type GetArgs struct {
  key string
}

func (p *Participant) Join(ja *JoinArgs, reply *Participant) error {
  *reply = p
  return nil
}

func (p *Participant) CanCommit(cca *CanCommitArgs, reply *bool) error {
  if value, ok := p.transactions[cca.tid]; ok {
    *reply = value.successful
    return nil
  }
  return fmt.Errorf("No such transaction in server")
}

func (p *Participant) DoCommit(dca *DoCommitArgs, reply *bool) error {
  if value, ok := p.transactions[dca.tid]; ok {
    value.Commit()
    *reply = true
    return nil
  }
  return fmt.Errorf("No such transaction in server")
}

func (p *Participant) DoAbort(daa *DoAbortArgs, reply *bool) error {
  if value, ok := p.transactions[dca.tid]; ok {
    value.Abort()
    *reply = true
    return nil
  }
  return fmt.Errorf("No such transaction in server")
}

func (p *Participant) SetKey(sa *SetArgs, reply *bool) error {
  // TODO
}

func (p *Participant) GetKey(ga *GetArgs, reply *bool) error {
  // TODO
}

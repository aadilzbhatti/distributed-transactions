package participant

type Participant struct {
  Objects map[string]Object
  Transactions map[int32]Transaction
  Address string
  Id int
}

func Start(hostname string, id int) error {
  // self := New(hostname, id)
  // // set up RPCs
  // if self == nil {
  //   return nil
  // }
  return nil
}

func New(addr string, id int) Participant {
  objs := make(map[string]Object, 0)
  trans := make(map[int32]Transaction, 0)
  return Participant{objs, trans, addr, id}
}

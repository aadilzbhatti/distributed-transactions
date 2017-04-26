package participant

type Participant struct {
  Objects map[string]string
  Transactions map[int]Transaction
  Address string
  Id int
}

func Start(hostname string, id int) error {
  self := New(hostname, id)
  // set up RPCs
}

func New(addr string, id int) Participant {
  objs := make(map[string]string, 0)
  trans := make(map[int]Transaction, 0)
  return Particpant{objs, trans, addr, id}
}

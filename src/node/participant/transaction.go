package participant

type Transaction struct {
  tid int32
  failed bool
}

func (t Transaction) Commit() {

}

func (t Transaction) Abort() {

}

package participant

type Transaction struct {
  Tid int32
  Failed bool
}

func (t Transaction) Commit() {

}

func (t Transaction) Abort() {

}

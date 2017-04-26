package participant 

type Transaction struct {
  tid int
  successful bool
}

func (t Transaction) Commit() {

}

func (t Transaction) Abort() {

}

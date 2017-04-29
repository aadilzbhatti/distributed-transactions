package participant

type Transaction struct {
	Tid    int32
	failed bool
  initial map[string]Object
	updates map[string]*Object
}

func NewTransaction(tid int32) *Transaction {
  return &Transaction{tid, false, make(map[string]Object, 0), make(map[string]*Object, 0)}
}

func (t *Transaction) addObject(key string, obj Object) {
	t.initial[key] = obj
}

func (t *Transaction) updateObject(key string, value *Object) {
	t.updates[key] = value
}

func (t *Transaction) hasFailed() bool {
	return t.failed
}

func (t *Transaction) commit() {
	t.failed = false
}

func (t *Transaction) abort() {
	t.failed = true
}

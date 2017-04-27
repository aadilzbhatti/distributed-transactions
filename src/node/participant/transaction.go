package participant

type Transaction struct {
	Tid    int32
	failed bool
}

func (t Transaction) hasFailed() bool {
	return t.failed
}

func (t *Transaction) commit() {
	t.failed = false
}

func (t *Transaction) abort() {
	t.failed = true
}

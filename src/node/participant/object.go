package participant

import (
	"fmt"
	"sync"
)

type Object struct {
	Key       string
	Value     string
	lock      *sync.RWMutex
	cond      *sync.Cond
	running   bool
	currTrans int32
}

func (o *Object) start() {
	o.lock.Lock()
	// o.running = true
	o.lock.Unlock()
}

func (o *Object) stop() {
	o.lock.Lock()
	o.running = false
	o.currTrans = 0
	o.cond.Broadcast()
	o.lock.Unlock()
}

func (o *Object) setKey(value string, trans int32) {
	o.lock.Lock()
	for o.running && trans != o.currTrans {
		o.cond.Wait()
	}
	fmt.Printf("In setKey: %v is value\n", value)
	o.Value = value
	o.running = true
	o.currTrans = trans
	o.lock.Unlock()
	fmt.Println(o)
}

func (o *Object) getKey(trans int32) string {
	o.lock.Lock()
	for o.running && trans != o.currTrans {
		o.cond.Wait()
	}
	o.lock.Unlock()
	fmt.Println("In getKey!")
	o.lock.RLock() // TODO fix this -- causes error with cond since unlocked?
	var res string
	res = o.Value
	o.lock.RUnlock()
	return res
}

func NewObject(key string, value string, trans int32) *Object {
	m := &sync.RWMutex{}
	c := sync.NewCond(m)
	return &Object{key, value, m, c, true, trans}
}

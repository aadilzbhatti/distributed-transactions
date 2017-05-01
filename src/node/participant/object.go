package participant

import (
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

func NewObject(key string, value string, trans int32) *Object {
	if _, ok := self.held[key]; ok {
		self.held[key].lock.Lock()
		for self.held[key].holding && self.held[key].currId != trans {
			self.held[key].cond.Wait()
		}
		self.held[key].lock.Unlock()
	}
	m := &sync.RWMutex{}
	c := sync.NewCond(m)
	return &Object{key, value, m, c, true, trans}
}

func (o *Object) stop() {
	o.lock.Lock()
	o.running = false
	o.currTrans = 0
	o.cond.Broadcast()
	o.lock.Unlock()
}

func (o *Object) resetKey(value string, trans int32) {
	o.lock.Lock()
	for o.running && trans != o.currTrans {
		o.cond.Wait()
	}
	o.Value = value
	o.currTrans = 0
	o.lock.Unlock()
}

func (o *Object) setKey(key string, value string, trans int32) {
	if _, ok := self.held[key]; ok {
		self.held[key].lock.Lock()
		for self.held[key].holding && self.held[key].currId != trans && self.held[key].currId != 0 {
			self.held[key].cond.Wait()
		}
		o.Value = value
		o.currTrans = trans
		self.held[key].lock.Unlock()
	} else {
		o.Value = value
		o.currTrans = trans
	}
}

func (o *Object) getKey(trans int32) string {
	key := o.Key
	if _, ok := self.held[key]; ok {
		self.held[key].lock.Lock()
		for self.held[key].holding && self.held[key].currId != trans && self.held[key].currId != 0 {
			self.held[key].cond.Wait()
		}
		self.held[key].lock.Unlock()
		o.lock.RLock()
		res := o.Value
		o.lock.RUnlock()
		return res
	}
	return o.Value
}

func (o *Object) copyObject(other Object) {
	o.Key = other.Key
	o.Value = other.Value
}

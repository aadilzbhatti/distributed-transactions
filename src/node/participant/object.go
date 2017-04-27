package participant

import (
  "sync"
)

type Object struct {
  Key string
  Value string
  Lock *sync.RWMutex
}

func (o Object) setKey(value string) {
  o.lock.Lock()
  o.value = value
  o.lock.Unlock()
}

func (o Object) getKey() string {
  var res string
  o.lock.RLock()
  res = o.value
  o.lock.RUnlock()
  return res
}

func NewObject(key string, value string) Object {
  var mutex = &sync.RWMutex{}
  return Object{key, value, mutex}
}

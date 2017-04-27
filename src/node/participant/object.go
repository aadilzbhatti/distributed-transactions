package participant

import (
  "sync"
)

type Object struct {
  Key string
  Value string
  lock *sync.RWMutex
}

func (o Object) setKey(value string) {
  o.lock.Lock()
  o.Value = value
  o.lock.Unlock()
}

func (o Object) getKey() string {
  var res string
  o.lock.RLock()
  res = o.Value
  o.lock.RUnlock()
  return res
}

func NewObject(key string, value string) Object {
  m := &sync.RWMutex{}
  return Object{key, value, m}
}

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
  o.Lock.Lock()
  o.Value = value
  o.Lock.Unlock()
}

func (o Object) getKey() string {
  var res string
  o.Lock.RLock()
  res = o.Value
  o.Lock.RUnlock()
  return res
}

func NewObject(key string, value string) Object {
  var mutex = &sync.RWMutex{}
  return Object{key, value, mutex}
}

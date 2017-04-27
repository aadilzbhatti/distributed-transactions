package participant

import (
  "sync"
)

var lock = &sync.RWMutex{}

type Object struct {
  Key string
  Value string
  // Lock *sync.RWMutex
}

func (o Object) setKey(value string) {
  lock.Lock()
  o.Value = value
  lock.Unlock()
}

func (o Object) getKey() string {
  var res string
  lock.RLock()
  res = o.Value
  lock.RUnlock()
  return res
}

func NewObject(key string, value string) Object {
  // var mutex = &sync.RWMutex{}
  return Object{key, value}
}

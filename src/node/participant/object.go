package participant

import (
  "sync"
  "fmt"
)

type Object struct {
  Key string
  Value string
  lock *sync.RWMutex
}

func (o Object) setKey(value string) {
  fmt.Printf("In setKey: %v is value\n", value)
  o.lock.Lock()
  o.Value = value
  o.lock.Unlock()
  fmt.Println(o)
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

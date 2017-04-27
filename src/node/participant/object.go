package participant

import (
  "sync"
  "fmt"
)

type Object struct {
  Key string
  Value string
  lock *sync.RWMutex
  writtenTo bool
}

func (o *Object) start() {
  o.lock.Lock()
  o.writtenTo = false
  o.lock.Unlock()
}

func (o *Object) setKey(value string) {
  for o.writtenTo {

  }
  fmt.Printf("In setKey: %v is value\n", value)
  o.lock.Lock()
  o.Value = value
  o.writtenTo = true
  o.lock.Unlock()
  fmt.Println(o)
}

func (o *Object) getKey() string {
  fmt.Println("In getKey!")
  var res string
  o.lock.RLock()
  res = o.Value
  o.lock.RUnlock()
  return res
}

func NewObject(key string, value string) *Object {
  m := &sync.RWMutex{}
  return &Object{key, value, m, false}
}

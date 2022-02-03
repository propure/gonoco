package dataStructure

import (
    "sync"

    "github.com/cheekybits/genny/generic"
)

type Key generic.Type
type Value generic.Type

type ValueDictionary struct {
    data map[Key]Value
    mux  sync.RWMutex
}

func (s *ValueDictionary) Set(key Key, value Value) {
    if s.data == nil {
        s.data = map[Key]Value{}
    }
    s.mux.Lock()
    defer s.mux.Unlock()
    s.data[key] = value
}

func (s *ValueDictionary) Delete(key Key) bool {
    s.mux.Lock()
    defer s.mux.Unlock()

    _, ok := s.data[key]
    if ok {
        delete(s.data, key)
    }

    return ok
}

func (s *ValueDictionary) Has(key Key) bool {
    s.mux.RLock()
    s.mux.RUnlock()

    _, result := s.data[key]

    return result
}

func (s *ValueDictionary) Get(key Key) Value {
    s.mux.RLock()
    s.mux.RUnlock()

    result, _ := s.data[key]

    return result
}

func (s *ValueDictionary) Clear() {
    s.mux.Lock()
    defer s.mux.Unlock()
    s.data = map[Key]Value{}
}

func (s *ValueDictionary) Size() int {
    return len(s.data)
}

func (s *ValueDictionary) Keys() []Key {
    s.mux.RLock()
    s.mux.RUnlock()

    keys := make([]Key, len(s.data))
    for k := range s.data {
        keys = append(keys, k)
    }

    return keys
}

func (s *ValueDictionary) Values() []Value {
    s.mux.RLock()
    s.mux.RUnlock()

    values := make([]Value, len(s.data))
    for _, v := range s.data {
        values = append(values, v)
    }

    return values
}

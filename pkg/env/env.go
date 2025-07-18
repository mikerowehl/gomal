package env

import (
	"github.com/mikerowehl/gomal/pkg/reader"
)

type Bindings struct {
	Data  map[reader.MalSymbol]reader.MalType
	Outer *Bindings
}

func (b *Bindings) Set(k reader.MalSymbol, v reader.MalType) {
	b.Data[k] = v
}

func (b *Bindings) Get(k reader.MalSymbol) (reader.MalType, bool) {
	if v, ok := b.Data[k]; ok {
		return v, ok
	}
	if b.Outer == nil {
		return nil, false
	}
	return b.Outer.Get(k)
}

func NewBindings(outer *Bindings) *Bindings {
	b := Bindings{
		Data:  make(map[reader.MalSymbol]reader.MalType),
		Outer: outer,
	}
	return &b
}

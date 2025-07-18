package env

import (
	"fmt"

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

func NewLambda(outer *Bindings, binds reader.MalList, exprs reader.MalList) (*Bindings, error) {
	if len(binds) != len(exprs) {
		return nil, fmt.Errorf("NewLambda: binds and exprs length mismatch")
	}
	b := Bindings{
		Data:  make(map[reader.MalSymbol]reader.MalType),
		Outer: outer,
	}
	for i := range binds {
		k, ok := binds[i].(reader.MalSymbol)
		if !ok {
			return nil, fmt.Errorf("NewLambda: lambda binding not a symbol: %v", binds[i])
		}
		b.Set(k, exprs[i])
	}
	return &b, nil
}

func NewLambdaVec(outer *Bindings, binds reader.MalVector, exprs reader.MalList) (*Bindings, error) {
	if len(binds) != len(exprs) {
		return nil, fmt.Errorf("NewLambdaVec: binds and exprs length mismatch")
	}
	b := Bindings{
		Data:  make(map[reader.MalSymbol]reader.MalType),
		Outer: outer,
	}
	for i := range binds {
		k, ok := binds[i].(reader.MalSymbol)
		if !ok {
			return nil, fmt.Errorf("NewLambda: lambda binding not a symbol: %v", binds[i])
		}
		b.Set(k, exprs[i])
	}
	return &b, nil
}

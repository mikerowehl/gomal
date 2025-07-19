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
	b := Bindings{
		Data:  make(map[reader.MalSymbol]reader.MalType),
		Outer: outer,
	}
	i := 0
	done := false
	for i < len(binds) && !done {
		k, ok := binds[i].(reader.MalSymbol)
		if !ok {
			return nil, fmt.Errorf("NewLambda: lambda binding not a symbol: %v", binds[i])
		}
		if k == reader.MalSymbol("&") {
			if i+1 >= len(binds) {
				return nil, fmt.Errorf("NewLambda: variadic marker with no following symbol")
			}
			variadic, ok := binds[i+1].(reader.MalSymbol)
			if !ok {
				return nil, fmt.Errorf("NewLambda: variadic market not a symbol")
			}
			b.Set(variadic, exprs[i:])
			done = true
		} else {
			b.Set(k, exprs[i])
		}
		i++
	}
	return &b, nil
}

func NewLambdaVec(outer *Bindings, binds reader.MalVector, exprs reader.MalList) (*Bindings, error) {
	b := Bindings{
		Data:  make(map[reader.MalSymbol]reader.MalType),
		Outer: outer,
	}
	i := 0
	done := false
	for i < len(binds) && !done {
		k, ok := binds[i].(reader.MalSymbol)
		if !ok {
			return nil, fmt.Errorf("NewLambda: lambda binding not a symbol: %v", binds[i])
		}
		if k == reader.MalSymbol("&") {
			if i+1 >= len(binds) {
				return nil, fmt.Errorf("NewLambda: variadic marker with no following symbol")
			}
			variadic, ok := binds[i+1].(reader.MalSymbol)
			if !ok {
				return nil, fmt.Errorf("NewLambda: variadic market not a symbol")
			}
			b.Set(variadic, exprs[i:])
			done = true
		} else {
			b.Set(k, exprs[i])
		}
		i++
	}
	return &b, nil
}

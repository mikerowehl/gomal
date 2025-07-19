package core

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mikerowehl/gomal/pkg/reader"
)

func listCast(raw reader.MalType, min int) (reader.MalList, error) {
	l, ok := raw.(reader.MalList)
	if !ok {
		return nil, fmt.Errorf("listCast: expected list: %v", raw)
	}
	if len(l) < min {
		return nil, fmt.Errorf("listCast: expected at least %d elements, got %d", min, len(l))
	}
	return l, nil
}

var add = func(a reader.MalType) (reader.MalType, error) {
	l, err := listCast(a, 2)
	if err != nil {
		return nil, err
	}
	i1, ok := l[0].(int)
	if !ok {
		return nil, fmt.Errorf("add: non int first arg to add: %v", l[0])
	}
	i2, ok := l[1].(int)
	if !ok {
		return nil, fmt.Errorf("add: non int second arg to add: %v", l[1])
	}
	return reader.MalType(i1 + i2), nil
}

var sub = func(a reader.MalType) (reader.MalType, error) {
	l, err := listCast(a, 2)
	if err != nil {
		return nil, err
	}
	i1, ok := l[0].(int)
	if !ok {
		return nil, fmt.Errorf("sub: non int first arg: %v", l[0])
	}
	i2, ok := l[1].(int)
	if !ok {
		return nil, fmt.Errorf("sub: non int second arg: %v", l[1])
	}
	return reader.MalType(i1 - i2), nil
}

var mul = func(a reader.MalType) (reader.MalType, error) {
	l, err := listCast(a, 2)
	if err != nil {
		return nil, err
	}
	i1, ok := l[0].(int)
	if !ok {
		return nil, fmt.Errorf("mul: non int first arg: %v", l[0])
	}
	i2, ok := l[1].(int)
	if !ok {
		return nil, fmt.Errorf("mul: non int second arg: %v", l[1])
	}
	return reader.MalType(i1 * i2), nil
}

var div = func(a reader.MalType) (reader.MalType, error) {
	l, err := listCast(a, 2)
	if err != nil {
		return nil, err
	}
	i1, ok := l[0].(int)
	if !ok {
		return nil, fmt.Errorf("div: non int first arg: %v", l[0])
	}
	i2, ok := l[1].(int)
	if !ok {
		return nil, fmt.Errorf("div: non int second arg: %v", l[1])
	}
	return reader.MalType(i1 / i2), nil
}

var prstr = func(a reader.MalType) (reader.MalType, error) {
	l, err := listCast(a, 0)
	if err != nil {
		return nil, err
	}
	s := []string{}
	for i := range l {
		s = append(s, reader.Pr_str(l[i], true))
	}
	return reader.MalType(strings.Join(s, " ")), nil
}

var str = func(a reader.MalType) (reader.MalType, error) {
	l, err := listCast(a, 0)
	if err != nil {
		return nil, err
	}
	s := []string{}
	for i := range l {
		s = append(s, reader.Pr_str(l[i], false))
	}
	return reader.MalType(strings.Join(s, "")), nil
}

var prn = func(a reader.MalType) (reader.MalType, error) {
	l, err := listCast(a, 0)
	if err != nil {
		return nil, err
	}
	s := []string{}
	for i := range l {
		s = append(s, reader.Pr_str(l[i], true))
	}
	fmt.Println(strings.Join(s, " "))
	return nil, nil
}

var printl = func(a reader.MalType) (reader.MalType, error) {
	l, err := listCast(a, 0)
	if err != nil {
		return nil, err
	}
	s := []string{}
	for i := range l {
		s = append(s, reader.Pr_str(l[i], false))
	}
	fmt.Println(strings.Join(s, " "))
	return nil, nil
}

var list = func(a reader.MalType) (reader.MalType, error) {
	return listCast(a, 0)
}

var islist = func(a reader.MalType) (reader.MalType, error) {
	l, err := listCast(a, 1)
	if err != nil {
		return nil, err
	}
	_, ok := l[0].(reader.MalList)
	if !ok {
		return reader.MalType(false), nil
	}
	return reader.MalType(true), nil
}

var empty = func(a reader.MalType) (reader.MalType, error) {
	l, err := listCast(a, 1)
	if err != nil {
		return nil, err
	}
	switch lp := l[0].(type) {
	case reader.MalList:
		return reader.MalType(bool(len(lp) == 0)), nil
	case reader.MalVector:
		return reader.MalType(bool(len(lp) == 0)), nil
	case reader.MalHashmap:
		return reader.MalType(bool(len(lp) == 0)), nil
	default:
		return reader.MalType(false), nil
	}
}

var count = func(a reader.MalType) (reader.MalType, error) {
	l, err := listCast(a, 1)
	if err != nil {
		return nil, err
	}
	switch lp := l[0].(type) {
	case reader.MalList:
		return reader.MalType(len(lp)), nil
	case reader.MalVector:
		return reader.MalType(len(lp)), nil
	default:
		return reader.MalType(0), nil
	}
}

type iterator struct {
	v    reader.MalType
	curr int
}

func NewMalIterator(v reader.MalType) (*iterator, error) {
	switch lv := v.(type) {
	case reader.MalList:
		ni := iterator{
			v:    lv,
			curr: 0,
		}
		return &ni, nil
	case reader.MalVector:
		ni := iterator{
			v:    lv,
			curr: 0,
		}
		return &ni, nil
	}
	return nil, fmt.Errorf("expected iterable value")
}

func (i *iterator) next() (reader.MalType, bool) {
	switch lv := i.v.(type) {
	case reader.MalList:
		if i.curr >= len(lv) {
			return nil, true
		} else {
			ret := lv[i.curr]
			i.curr++
			return ret, false
		}
	case reader.MalVector:
		if i.curr >= len(lv) {
			return nil, true
		} else {
			ret := lv[i.curr]
			i.curr++
			return ret, false
		}
	}
	return nil, true
}

func seriesCompare(s1 reader.MalType, s2 reader.MalType) (reader.MalType, error) {
	i1, err := NewMalIterator(s1)
	if err != nil {
		return nil, err
	}
	i2, err := NewMalIterator(s2)
	if err != nil {
		return nil, err
	}
	for {
		v1, done1 := i1.next()
		v2, done2 := i2.next()
		if done1 && done2 {
			return reader.MalType(true), nil
		}
		if done1 || done2 {
			return reader.MalType(false), nil
		}
		check, err := equalItems(v1, v2)
		if err != nil {
			return nil, err
		}
		if check != reader.MalType(true) {
			return reader.MalType(false), nil
		}
	}
}

func equalItems(v1 reader.MalType, v2 reader.MalType) (reader.MalType, error) {
	v1t := reflect.TypeOf(v1)
	v2t := reflect.TypeOf(v2)
	if v1t == reader.MalListType || v1t == reader.MalVectorType {
		if v2t == reader.MalListType || v2t == reader.MalVectorType {
			return seriesCompare(v1, v2)
		}
		return reader.MalType(false), nil
	}
	if v1t == reader.MalHashmapType || v2t == reader.MalHashmapType {
		return reader.MalType(false), nil
	}
	return reader.MalType(v1 == v2), nil
}

var equal = func(a reader.MalType) (reader.MalType, error) {
	l, err := listCast(a, 2)
	if err != nil {
		return nil, fmt.Errorf("equal: invalid arguments: %v", a)
	}
	return equalItems(l[0], l[1])
}

var lt = func(a reader.MalType) (reader.MalType, error) {
	l, err := listCast(a, 2)
	if err != nil {
		return nil, err
	}
	i1, ok := l[0].(int)
	if !ok {
		return nil, fmt.Errorf("lt: non int first arg: %v", l[0])
	}
	i2, ok := l[1].(int)
	if !ok {
		return nil, fmt.Errorf("lt: non int second arg: %v", l[1])
	}
	return reader.MalType(i1 < i2), nil
}

var lte = func(a reader.MalType) (reader.MalType, error) {
	l, err := listCast(a, 2)
	if err != nil {
		return nil, err
	}
	i1, ok := l[0].(int)
	if !ok {
		return nil, fmt.Errorf("lte: non int first arg: %v", l[0])
	}
	i2, ok := l[1].(int)
	if !ok {
		return nil, fmt.Errorf("lte: non int second arg: %v", l[1])
	}
	return reader.MalType(i1 <= i2), nil
}

var gt = func(a reader.MalType) (reader.MalType, error) {
	l, err := listCast(a, 2)
	if err != nil {
		return nil, err
	}
	i1, ok := l[0].(int)
	if !ok {
		return nil, fmt.Errorf("gt: non int first arg: %v", l[0])
	}
	i2, ok := l[1].(int)
	if !ok {
		return nil, fmt.Errorf("gt: non int second arg: %v", l[1])
	}
	return reader.MalType(i1 > i2), nil
}

var gte = func(a reader.MalType) (reader.MalType, error) {
	l, err := listCast(a, 2)
	if err != nil {
		return nil, err
	}
	i1, ok := l[0].(int)
	if !ok {
		return nil, fmt.Errorf("gte: non int first arg: %v", l[0])
	}
	i2, ok := l[1].(int)
	if !ok {
		return nil, fmt.Errorf("gte: non int second arg: %v", l[1])
	}
	return reader.MalType(i1 >= i2), nil
}
var NS = map[reader.MalSymbol]reader.MalFunc{
	"+":       add,
	"-":       sub,
	"*":       mul,
	"/":       div,
	"pr-str":  prstr,
	"str":     str,
	"prn":     prn,
	"println": printl,
	"list":    list,
	"list?":   islist,
	"empty?":  empty,
	"count":   count,
	"=":       equal,
	"<":       lt,
	"<=":      lte,
	">":       gt,
	">=":      gte,
}

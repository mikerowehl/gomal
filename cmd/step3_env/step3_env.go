package main

import (
	"bufio"
	"fmt"
	"github.com/mikerowehl/gomal/pkg/env"
	"github.com/mikerowehl/gomal/pkg/reader"
	"os"
)

// Returned values are:
//
//	string - token/line
//	bool - eof, true means end of input
//	error - set to nil unless there's an error
func READ(scanner *bufio.Scanner) (reader.MalType, bool, error) {
	fmt.Print("user> ")

	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", false, fmt.Errorf("READ unable to get line: %w", err)
		} else {
			return "", true, nil
		}
	}
	val, err := reader.Read_str(scanner.Text())
	return val, false, err
}

func APPLY(v reader.MalType) (reader.MalType, error) {
	if l, ok := v.(reader.MalList); !ok {
		fmt.Println("Applying something that isn't a list")
	} else {
		a, ok := l[0].(reader.MalFunc)
		if !ok {
			return nil, fmt.Errorf("Error converting apply function")
		}
		return a(l[1:])
	}
	return nil, fmt.Errorf("Error applying")
}

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

func EVAL(v reader.MalType, b *env.Bindings) (reader.MalType, error) {
	debug, ok := b.Get(reader.MalSymbol("DEBUG-EVAL"))
	if ok && debug != nil && debug != false {
		fmt.Print("EVAL: ")
		reader.Pr_str(v, true)
		fmt.Println()
	}
	switch t := v.(type) {
	case reader.MalSymbol:
		entry, ok := b.Get(t)
		if !ok {
			return nil, fmt.Errorf("eval: %v not found", t)
		}
		return entry, nil
	case reader.MalList:
		if len(t) == 0 {
			return v, nil
		}
		car := t[0].(reader.MalSymbol)
		if car == "def!" {
			key := t[1].(reader.MalSymbol)
			val, err := EVAL(t[2], b)
			if err != nil {
				return nil, fmt.Errorf("eval: error evaluating def! value: %w", err)
			}
			b.Set(key, val)
			return val, nil
		} else if car == "let*" {
			nenv := env.NewBindings(b)
			switch lt := t[1].(type) {
			case reader.MalList:
				for i := 0; i < len(lt); i += 2 {
					key := lt[i].(reader.MalSymbol)
					val, err := EVAL(lt[i+1], nenv)
					if err != nil {
						return nil, fmt.Errorf("eval: error evaluating let* value: %w", err)
					}
					nenv.Set(key, val)
				}
				return EVAL(t[2], nenv)
			case reader.MalVector:
				for i := 0; i < len(lt); i += 2 {
					key := lt[i].(reader.MalSymbol)
					val, err := EVAL(lt[i+1], nenv)
					if err != nil {
						return nil, fmt.Errorf("eval: error evaluating let* vector: %w", err)
					}
					nenv.Set(key, val)
				}
				return EVAL(t[2], nenv)
			default:
				return nil, fmt.Errorf("eval: invalid type for let* bindings: %v", t[1])
			}
		}
		evaled := reader.MalList{}
		for _, entry := range t {
			n, err := EVAL(entry, b)
			if err != nil {
				return nil, fmt.Errorf("eval: error evaluating list for apply: %w", err)
			}
			evaled = append(evaled, n)
		}
		app, err := APPLY(evaled)
		if err != nil {
			return nil, fmt.Errorf("eval: error during apply: %w", err)
		}
		return app, nil
	case reader.MalVector:
		evaled := reader.MalVector{}
		for _, entry := range t {
			n, err := EVAL(entry, b)
			if err != nil {
				return nil, fmt.Errorf("eval: error evaluating vector contents: %w", err)
			}
			evaled = append(evaled, n)
		}
		return evaled, nil
	case reader.MalHashmap:
		evaled := reader.MalHashmap{}
		for i, entry := range t {
			if (i % 2) == 1 {
				n, err := EVAL(entry, b)
				if err != nil {
					return nil, fmt.Errorf("eval: error evaluating hashmap values: %w", err)
				}
				evaled = append(evaled, n)
			} else {
				evaled = append(evaled, entry)
			}
		}
		return evaled, nil
	}
	return v, nil
}

func PRINT(v reader.MalType) {
	reader.Pr_str(v, true)
	fmt.Println()
}

func rep() {
	scanner := bufio.NewScanner(os.Stdin)
	env := env.NewBindings(nil)
	env.Set(reader.MalSymbol("+"), reader.MalFunc(add))
	env.Set(reader.MalSymbol("-"), reader.MalFunc(sub))
	env.Set(reader.MalSymbol("*"), reader.MalFunc(mul))
	env.Set(reader.MalSymbol("/"), reader.MalFunc(div))
	for {
		line, eof, err := READ(scanner)
		if eof {
			fmt.Println()
			return
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
		} else {
			e, err := EVAL(line, env)
			if err != nil {
				fmt.Printf("ERR: %v\n", err)
			} else {
				PRINT(e)
			}
		}
	}
}

func main() {
	rep()
}

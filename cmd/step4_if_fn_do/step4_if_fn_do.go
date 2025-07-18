package main

import (
	"bufio"
	"fmt"
	"github.com/mikerowehl/gomal/pkg/core"
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

func truthy(v reader.MalType) bool {
	switch t := v.(type) {
	case bool:
		return t
	default:
		return true
	}
}

func EVAL(v reader.MalType, b *env.Bindings) (reader.MalType, error) {
	debug, ok := b.Get(reader.MalSymbol("DEBUG-EVAL"))
	if ok && debug != nil && debug != false {
		fmt.Printf("EVAL: %s\n", reader.Pr_str(v, true))
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
		car, ok := t[0].(reader.MalSymbol)
		if ok {
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
			} else if car == "fn*" {
				switch binds := t[1].(type) {
				case reader.MalList:
					fn := func(exprs reader.MalType) (reader.MalType, error) {
						nexpr, ok := exprs.(reader.MalList)
						if !ok {
							return nil, fmt.Errorf("eval: bad type for expressions in fn*: %v", exprs)
						}
						nenv, err := env.NewLambda(b, binds, nexpr)
						if err != nil {
							return nil, fmt.Errorf("eval: error creating new fn* bindings: %w", err)
						}
						return EVAL(t[2], nenv)
					}
					return reader.MalFunc(fn), nil
				case reader.MalVector:
					fn := func(exprs reader.MalType) (reader.MalType, error) {
						nexpr, ok := exprs.(reader.MalList)
						if !ok {
							return nil, fmt.Errorf("eval: bad type for expressions in fn*: %v", exprs)
						}
						nenv, err := env.NewLambdaVec(b, binds, nexpr)
						if err != nil {
							return nil, fmt.Errorf("eval: error creating new fn* bindings: %w", err)
						}
						return EVAL(t[2], nenv)
					}
					return reader.MalFunc(fn), nil
				}
			} else if car == "do" {
				var retVal reader.MalType
				var err error
				for i := 1; i < len(t); i++ {
					retVal, err = EVAL(t[i], b)
					if err != nil {
						return nil, err
					}
				}
				return retVal, nil
			} else if car == "if" {
				cond, err := EVAL(t[1], b)
				if err != nil {
					return nil, err
				}
				var evalExpr reader.MalType
				if cond != nil && truthy(cond) {
					if len(t) >= 3 {
						evalExpr = t[2]
					} else {
						return nil, nil
					}
				} else {
					if len(t) >= 4 {
						evalExpr = t[3]
					} else {
						return nil, nil
					}
				}
				return EVAL(evalExpr, b)
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
	fmt.Println(reader.Pr_str(v, true))
}

func rep() {
	scanner := bufio.NewScanner(os.Stdin)
	env := env.NewBindings(nil)
	for k, v := range core.NS {
		env.Set(k, v)
	}
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

package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type ValueKind int

const (
	Integer ValueKind = iota
	Float
	Boolean
	String
	Symbol
	Name
	PairValue
	Null
	Function
)

type Value struct {
	kind ValueKind
	val  interface{}
}

type Node struct {
	item *Value
	next *Node
}

type Pair struct {
	first  *Value
	second *Value
}

var nullValue *Value = &Value{
	kind: Null,
	val:  nil,
}

func (vk ValueKind) String() string {
	var kind string
	switch vk {
	case Integer:
		kind = "Integer"
	case Float:
		kind = "Float"
	case Boolean:
		kind = "Boolean"
	case String:
		kind = "String"
	case Symbol:
		kind = "Symbol"
	case Name:
		kind = "Name"
	case PairValue:
		kind = "Pair"
	case Function:
		kind = "Function"
	case Null:
		kind = "Null"
	}

	return kind
}

func (v *Value) String() string {
	switch v.kind {
	case Integer:
		return fmt.Sprintf("%d", v.val)
	case Float:
		return fmt.Sprintf("%f", v.val)
	case Boolean:
		return fmt.Sprintf("%v", v.val)
	case String:
		return fmt.Sprintf("\"%s\"", v.val)
	case Symbol:
		return fmt.Sprintf("'%s", v.val)
	case Name:
		return fmt.Sprintf("%s", v.val)
	case Function:
		return fmt.Sprintf("%v", v.val)
	case PairValue:
		p := v.val.(*Pair)
		f := p.first
		s := p.second
		if isPair(f) == false && isPair(s) == false {
			return fmt.Sprintf("(%s . %s)", f, s)
		} else {
			var w strings.Builder
			printList(&w, v)
			return fmt.Sprintf("%s", w.String())
		}
	case Null:
		return fmt.Sprintf("<null>")
	default:
		panic(fmt.Sprintf("invalid value %v", v))
	}
}

func cons(f *Value, s *Value) *Value {
	return &Value{
		kind: PairValue,
		val: &Pair{
			first:  f,
			second: s,
		},
	}
}

func car(v *Value) *Value {
	if v == nil {
		panic("car: value is nil")
	}

	if v.kind != PairValue {
		panic("car: value is not a pair")
	}

	p, ok := v.val.(*Pair)
	if !ok {
		panic("car: value is not a proper pair")
	}

	return p.first
}

func setCar(p *Value, val *Value) {
	pair, ok := p.val.(*Pair)
	if !ok {
		panic("setCar: p is not a pair")
	}
	set(&pair.first, val)
}

func cdr(v *Value) *Value {
	if v == nil {
		panic("cdr: value is null")
	}

	if v.kind != PairValue {
		panic("cdr: value is not a pair")
	}

	p, ok := v.val.(*Pair)
	if !ok {
		panic("cdr: value is not a proper pair")
	}

	return p.second
}

func setCdr(p *Value, val *Value) {
	pair, ok := p.val.(*Pair)
	if !ok {
		panic("setCar: p is not a pair")
	}
	set(&pair.second, val)
}

func list(items ...*Value) *Value {
	if len(items) == 0 {
		return nullValue
	}

	var first *Value
	var current *Value
	null := &Value{
		kind: Null,
		val:  nil,
	}

	for i, item := range items {
		if i == 0 {
			first = cons(item, null)
			current = first
			continue
		}

		next := cons(item, null)
		current.val.(*Pair).second = next
		current = next
	}

	return first
}

func listAppend(l1 *Value, l2 *Value) *Value {
	if isNull(l1) {
		return l2
	}
	return cons(car(l1), listAppend(cdr(l1), l2))
}

func listLen(l *Value) int {
	if isNull(l) {
		return 0
	}
	return 1 + listLen(cdr(l))
}

func set(dest **Value, src *Value) {
	*dest = src
}

func printList(output io.Writer, v *Value) {
	if v.kind == Null {
		fmt.Fprintf(output, "()\n")
		return
	}

	_, ok := v.val.(*Pair)
	if !ok {
		panic("list: value is not a pair")
	}

	current := v
	fmt.Fprintf(output, "(")

	for {
		if current.kind == Null {
			break
		}

		fmt.Fprintf(output, "%v", current.val.(*Pair).first)
		current = current.val.(*Pair).second

		if current.kind != Null {
			fmt.Fprintf(output, " ")
		}
	}
	fmt.Fprintf(output, ")")
}

func isPair(v *Value) bool {
	if v == nil || v.kind == Null {
		return false
	}

	_, ok := v.val.(*Pair)
	if ok {
		return true
	}

	return false
}

func isNull(v *Value) bool {
	if v == nil {
		panic("not a value")
	}

	if v.kind == Null {
		return true
	}

	return false
}

func isNumber(v *Value) bool {
	if v == nil {
		panic("not a value")
	}
	if v.kind == Integer || v.kind == Float {
		return true
	}
	return false
}

func isString(v *Value) bool {
	if v == nil {
		panic("not a value")
	}
	if v.kind == String {
		return true
	}
	return false
}

func isSymbol(v *Value) bool {
	if v == nil {
		panic("not a value")
	}
	if v.kind == Symbol {
		return true
	}
	return false
}

func isName(v *Value) bool {
	if v == nil {
		panic("not a value")
	}
	if v.kind == Name {
		return true
	}
	return false
}

func isEqual(v1 *Value, v2 *Value) bool {
	if v1 == nil {
		panic("not a value")
	}
	if v2 == nil {
		panic("not a value")
	}
	if v1.kind != v2.kind {
		return false
	}

	switch v1.kind {
	case Integer:
		return v1.val.(int64) == v2.val.(int64)
	case Float:
		return v1.val.(float64) == v2.val.(float64)
	case Boolean:
		return v1.val.(bool) == v2.val.(bool)
	case String:
		return v1.val.(string) == v2.val.(string)
	case Symbol:
		return v1.val.(string) == v2.val.(string)
	case Name:
		return v1.val.(string) == v2.val.(string)
	case PairValue:
		panic("cannot compare pairs")
	case Function:
		return true
	case Null:
		return true
	}

	panic("unreachable")
}

func isTrue(v *Value) bool {
	if v == nil {
		panic("not a value")
	}
	if v.kind != Boolean {
		panic("not a Boolean")
	}
	return v.val.(bool)
}

func _map(fun func(*Value) *Value, items *Value) *Value {
	if isNull(items) {
		return nullValue
	}

	item := fun(car(items))
	return cons(item, _map(fun, cdr(items)))
}

func not(b *Value) bool {
	if b == nil {
		panic("not a value")
	}
	if b.kind != Boolean {
		panic("not a Boolean")
	}
	val := b.val.(bool)
	return !val
}

func cadr(v *Value) *Value {
	return car(cdr(v))
}

func caddr(v *Value) *Value {
	return car(cdr(cdr(v)))
}

func caadr(v *Value) *Value {
	return car(car(cdr(v)))
}

func cadddr(v *Value) *Value {
	return car(cdr(cdr(cdr(v))))
}

func cddr(v *Value) *Value {
	return cdr(cdr(v))
}

func cdddr(v *Value) *Value {
	return cdr(cdr(cdr(v)))
}

func cdadr(v *Value) *Value {
	return cdr(car(cdr(v)))
}

func gt(args *Value) *Value {
	v1 := car(args)
	v2 := cadr(args)
	var a float64
	var b float64

	if v1.kind == Integer {
		a = float64(v1.val.(int64))
	} else if v1.kind == Float {
		a = v1.val.(float64)
	} else {
		panic(fmt.Sprintf("comparison is only allowed on numbers %s", v1))
	}

	if v2.kind == Integer {
		b = float64(v2.val.(int64))
	} else if v2.kind == Float {
		b = v2.val.(float64)
	} else {
		panic(fmt.Sprintf("comparison is only allowed on numbers %s", v2))
	}

	if a > b {
		return make_true()
	}
	return make_false()
}

func lt(args *Value) *Value {
	v1 := car(args)
	v2 := cadr(args)
	var a float64
	var b float64

	if v1.kind == Integer {
		a = float64(v1.val.(int64))
	} else if v1.kind == Float {
		a = v1.val.(float64)
	} else {
		panic(fmt.Sprintf("comparison is only allowed on numbers %s", v1))
	}

	if v2.kind == Integer {
		b = float64(v2.val.(int64))
	} else if v2.kind == Float {
		b = v2.val.(float64)
	} else {
		panic(fmt.Sprintf("comparison is only allowed on numbers %s", v2))
	}

	if a < b {
		return make_true()
	}
	return make_false()
}

func or(args *Value) *Value {
	var res bool

	for {
		if isNull(args) {
			break
		}
		v := car(args)
		if isTrue(v) {
			res = true
			break
		}
		args = cdr(args)
	}
	if res {
		return make_true()
	}
	return make_false()
}

func and(args *Value) *Value {
	res := true

	for {
		if isNull(args) {
			break
		}
		v := car(args)
		if !isTrue(v) {
			res = false
			break
		}
		args = cdr(args)
	}
	if res {
		return make_true()
	}
	return make_false()
}

func eq(args *Value) *Value {
	v1 := car(args)
	v2 := cadr(args)

	if isNumber(v1) && isNumber(v2) {
		var a float64
		var b float64

		if v1.kind == Integer {
			a = float64(v1.val.(int64))
		} else {
			a = v1.val.(float64)
		}
		if v2.kind == Integer {
			b = float64(v2.val.(int64))
		} else {
			b = v2.val.(float64)
		}

		if a == b {
			return make_true()
		}
		return make_false()
	}

	if isEqual(v1, v2) {
		return make_true()
	}
	return make_false()
}

func plus(args *Value) *Value {
	var sum float64

	for {
		if isNull(args) {
			break
		}
		v := car(args)
		if v.kind == Integer {
			sum += float64(v.val.(int64))
		} else if v.kind == Float {
			sum += v.val.(float64)
		}

		args = cdr(args)
	}

	return &Value{
		kind: Float,
		val:  sum,
	}
}

func minus(args *Value) *Value {
	var res float64
	first := car(args)
	if first.kind == Integer {
		res = float64(first.val.(int64))
	} else {
		res = first.val.(float64)
	}

	args = cdr(args)

	for {
		if isNull(args) {
			break
		}
		v := car(args)
		if v.kind == Integer {
			res -= float64(v.val.(int64))
		} else if v.kind == Float {
			res -= v.val.(float64)
		}

		args = cdr(args)
	}

	return &Value{
		kind: Float,
		val:  res,
	}
}

func mul(args *Value) *Value {
	var res float64 = 1.0

	for {
		if isNull(args) {
			break
		}
		v := car(args)
		if v.kind == Integer {
			res *= float64(v.val.(int64))
		} else if v.kind == Float {
			res *= v.val.(float64)
		}

		args = cdr(args)
	}

	return &Value{
		kind: Float,
		val:  res,
	}
}

func display(args *Value) *Value {
	var v *Value
	if isPair(args) {
		v = car(args)
	} else {
		v = args
	}
	fmt.Fprintf(os.Stdout, "%s", v)
	return constant("ok")
}

func newline(args *Value) *Value {
	fmt.Println()
	return constant("ok")
}

package main

import (
	"fmt"
	"os"
)

func tests() {
	one := &Value{
		kind: Integer,
		val:  1,
	}
	two := &Value{
		kind: Integer,
		val:  2,
	}
	three := &Value{
		kind: Integer,
		val:  3,
	}

	// cons tests
	p := cons(one, two)
	fmt.Printf("car: %s\n", car(p))
	fmt.Printf("cdr: %s\n", cdr(p))

	// printList tests
	l := list(one, two)
	printList(os.Stdout, l)
	fmt.Println()

	l2 := cons(one, cons(two, cons(three, nullValue)))
	printList(os.Stdout, l2)
	fmt.Println()

	l3 := list()
	printList(os.Stdout, l3)
	fmt.Println()

	// listAppend tests
	l4 := list(one, two)
	l5 := list(three)
	fmt.Print("listAppend: ")
	printList(os.Stdout, listAppend(l4, l5))
	fmt.Println()

	// listLen tests
	fmt.Printf("listLen: %d\n", listLen(l4))
	fmt.Printf("listLen: %d\n", listLen(l5))
	fmt.Printf("listLen: %d\n", listLen(listAppend(l4, l5)))

	// _map tests
	plusOne := func(v *Value) *Value {
		val := v.val.(int)
		return &Value{
			kind: Integer,
			val:  val + 1,
		}
	}
	items := cons(one, cons(two, cons(three, nullValue)))
	printList(os.Stdout, _map(plusOne, items))
	fmt.Println()

	// set tests
	fmt.Printf("%s\n", one)
	x := &Value{kind: Integer, val: 55}
	set(&one, x)
	fmt.Printf("%s\n", one)

	// set car tests
	fmt.Printf("%s\n", p)
	setCar(p, x)
	fmt.Printf("%s\n", p)

}

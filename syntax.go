package main

import "fmt"

// self evaluating items (numbers and strings)
func is_self_evaluating(exp *Value) *Value {
	if isNumber(exp) || isString(exp) {
		return make_true()
	}
	return make_false()
}

// variables and quotations
func is_variable(exp *Value) *Value {
	if isName(exp) {
		return make_true()
	}
	return make_false()
}

func is_quoted(exp *Value) *Value {
	if isSymbol(exp) {
		return make_true()
	}
	return make_false()
}

func text_of_quotation(exp *Value) *Value {
	return cadr(exp)
}

func is_tagged_list(exp *Value, tag string) bool {
	if isPair(exp) {
		name := car(exp)
		return name.val.(string) == tag
	}
	return false
}

// assignments with set!
func is_assignment(exp *Value) *Value {
	if is_tagged_list(exp, "set!") {
		return make_true()
	}
	return make_false()
}

func assignment_variable(exp *Value) *Value {
	return cadr(exp)
}
func assignment_value(exp *Value) *Value {
	return caddr(exp)
}

// definitions
func is_definition(exp *Value) *Value {
	if is_tagged_list(exp, "define") {
		return make_true()
	}
	return make_false()
}

func definition_variable(exp *Value) *Value {
	if isName(cadr(exp)) {
		return cadr(exp)
	}
	return caadr(exp)
}

func definition_value(exp *Value) *Value {
	if isName(cadr(exp)) {
		return caddr(exp)
	}
	return make_lambda(cdadr(exp), cddr(exp)) // parameters and body
}

// lambda expressions
func is_lambda(exp *Value) *Value {
	if is_tagged_list(exp, "lambda") {
		return make_true()
	}
	return make_false()
}

func lambda_parameters(exp *Value) *Value {
	return cadr(exp)
}
func lambda_body(exp *Value) *Value {
	return cddr(exp)
}

func make_lambda(parameters *Value, body *Value) *Value {
	lamb := &Value{
		kind: Name,
		val:  "lambda",
	}
	return cons(lamb, cons(parameters, body))
}

// if conditionals
func is_if(exp *Value) *Value {
	if is_tagged_list(exp, "if") {
		return make_true()
	}
	return make_false()
}
func if_predicate(exp *Value) *Value {
	return cadr(exp)
}
func if_consequent(exp *Value) *Value {
	return caddr(exp)
}

func if_alternative(exp *Value) *Value {
	if !isNull(cdddr(exp)) {
		return cadddr(exp)
	}
	return &Value{
		kind: Boolean,
		val:  false,
	}
}

func make_if(predicate *Value, consequent *Value, alternative *Value) *Value {
	if_name := &Value{
		kind: Name,
		val:  "if",
	}
	return list(if_name, predicate, consequent, alternative)
}

// begin expressions
func is_begin(exp *Value) *Value {
	if is_tagged_list(exp, "begin") {
		return make_true()
	}
	return make_false()
}
func begin_actions(exp *Value) *Value {
	return cdr(exp)
}
func is_last_exp(seq *Value) *Value {
	if isNull(cdr(seq)) {
		return make_true()
	}
	return make_false()
}
func first_exp(seq *Value) *Value {
	return car(seq)
}
func rest_exps(seq *Value) *Value {
	return cdr(seq)
}

func sequence_to_exp(seq *Value) *Value {
	if isNull(seq) {
		return seq
	} else if is_last_exp(seq).val.(bool) == true {
		return first_exp(seq)
	} else {
		return make_begin(seq)
	}
}

func make_begin(seq *Value) *Value {
	begin_name := &Value{
		kind: Name,
		val:  "begin",
	}
	return cons(begin_name, seq)
}

// let expressions
func is_let(exp *Value) bool {
	return is_tagged_list(exp, "let")
}

func let_bindings(exp *Value) *Value {
	return cadr(exp)
}
func let_body(exp *Value) *Value {
	return cddr(exp)
}

func make_application(operator *Value, operands *Value) *Value {
	return cons(operator, operands)
}

// let to combination transformation
func let_to_combination(exp *Value) *Value {
	lamb := make_lambda(
		_map(car, let_bindings(exp)),
		let_body(exp))

	return make_application(lamb, _map(cadr, let_bindings(exp)))
}

// procedure applications
func is_application(exp *Value) *Value {
	if isPair(exp) {
		return make_true()
	}
	return make_false()
}

func operator(exp *Value) *Value {
	return car(exp)
}
func operands(exp *Value) *Value {
	return cdr(exp)
}
func has_no_operands(ops *Value) *Value {
	if isNull(ops) {
		return make_true()
	}
	return make_false()
}
func first_operand(ops *Value) *Value {
	return car(ops)
}
func rest_operands(ops *Value) *Value {
	return cdr(ops)
}
func is_last_operand(ops *Value) *Value {
	if isNull(cdr(ops)) {
		return make_true()
	}
	return make_false()
}

// derived expression with cond
func is_cond(exp *Value) bool {
	return is_tagged_list(exp, "cond")
}
func cond_clauses(exp *Value) *Value {
	return cdr(exp)
}
func is_cond_else_clause(clause *Value) bool {
	pred := cond_predicate(clause)
	if pred.val.(string) == "else" {
		return true
	}
	return false
}

func cond_predicate(clause *Value) *Value {
	return car(clause)
}
func cond_actions(clause *Value) *Value {
	return cdr(clause)
}

func cond_to_if(exp *Value) *Value {
	return expand_clauses(cond_clauses(exp))
}

func expand_clauses(clauses *Value) *Value {
	if isNull(clauses) {
		return &Value{
			kind: Boolean,
			val:  false,
		}
	}

	first := car(clauses)
	rest := cdr(clauses)
	if is_cond_else_clause(first) {
		if isNull(rest) {
			return sequence_to_exp(cond_actions(first))
		} else {
			panic(fmt.Sprintf("ELSE clause isn't last -- cond_to_if %s", clauses))
		}
	} else {
		return make_if(
			cond_predicate(first),
			sequence_to_exp(cond_actions(first)),
			expand_clauses(rest))
	}
}

func is_true(v *Value) *Value {
	if isTrue(v) {
		return make_true()
	}
	return make_false()
}

func make_prim(p func(args *Value) *Value) *Value {
	return &Value{
		kind: Function,
		val:  p,
	}
}

func make_name(n string) *Value {
	return &Value{
		kind: Name,
		val:  n,
	}
}

func make_false() *Value {
	return &Value{
		kind: Boolean,
		val:  false,
	}
}

func make_true() *Value {
	return &Value{
		kind: Boolean,
		val:  true,
	}
}

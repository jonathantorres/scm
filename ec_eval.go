package main

import "fmt"

type Register struct {
	name     string
	contents *Value
}

func (r *Register) String() string {
	return fmt.Sprintf("%s: %s", r.name, r.contents)
}

// the machine stack
var stack *Stack

// machine registers
var exp *Register = &Register{name: "exp"}
var env *Register = &Register{name: "env"}
var unev *Register = &Register{name: "unev"}
var argl *Register = &Register{name: "argl"}
var proc *Register = &Register{name: "proc"}
var cont *Register = &Register{name: "cont"}
var val *Register = &Register{name: "val"}

func startEval(v *Value) {
	initialize_stack()
	assign(exp, v)
	assign(env, get_global_environment())
	assign(cont, label(print_result))
	go_to(label(eval_dispatch))
}

func eval_dispatch() {
	if test(is_self_evaluating(reg(exp))) {
		ev_self_eval()
		return
	}

	if test(is_variable(reg(exp))) {
		ev_variable()
		return
	}

	if test(is_quoted(reg(exp))) {
		ev_quoted()
		return
	}

	if test(is_assignment(reg(exp))) {
		ev_assignment()
		return
	}

	if test(is_definition(reg(exp))) {
		ev_definition()
		return
	}

	if test(is_if(reg(exp))) {
		ev_if()
		return
	}

	if test(is_lambda(reg(exp))) {
		ev_lambda()
		return
	}

	if test(is_begin(reg(exp))) {
		ev_begin()
		return
	}

	if test(is_application(reg(exp))) {
		ev_application()
		return
	}

	go_to(label(unknown_expression_type))
}

func ev_lambda() {
	assign(unev, lambda_parameters(reg(exp)))
	assign(exp, lambda_body(reg(exp)))
	assign(val, make_procedure(reg(unev), reg(exp), reg(env)))
	go_to(reg(cont))
}

func ev_application() {
	save(*cont)
	save(*env)
	assign(unev, operands(reg(exp)))
	save(*unev)
	assign(exp, operator(reg(exp)))
	assign(cont, label(ev_appl_did_operator))
	go_to(label(eval_dispatch))
}

func ev_appl_did_operator() {
	restore(unev)
	restore(env)
	assign(argl, empty_arglist())
	assign(proc, reg(val))
	if test(has_no_operands(reg(unev))) {
		apply_dispatch()
		return
	}
	save(*proc)
	go_to(label(ev_appl_operand_loop))
}

func ev_appl_operand_loop() {
	save(*argl)
	assign(exp, first_operand(reg(unev)))
	if test(is_last_operand(reg(unev))) {
		ev_appl_last_arg()
		return
	}
	save(*env)
	save(*unev)
	assign(cont, label(ev_appl_accumulate_arg))
	go_to(label(eval_dispatch))
}

func ev_appl_accumulate_arg() {
	restore(unev)
	restore(env)
	restore(argl)
	assign(argl, adjoin_arg(reg(val), reg(argl)))
	assign(unev, rest_operands(reg(unev)))
	go_to(label(ev_appl_operand_loop))
}

func ev_appl_last_arg() {
	assign(cont, label(ev_appl_accum_last_arg))
	go_to(label(eval_dispatch))
}

func ev_appl_accum_last_arg() {
	restore(argl)
	assign(argl, adjoin_arg(reg(val), reg(argl)))
	restore(proc)
	go_to(label(apply_dispatch))
}

func apply_dispatch() {
	if test(is_primitive_procedure(reg(proc))) {
		primitive_apply()
		return
	}
	if test(is_compound_procedure(reg(proc))) {
		compound_apply()
		return
	}
	go_to(label(unknown_procedure_type))
}

func primitive_apply() {
	assign(val, apply_primitive_procedure(reg(proc), reg(argl)))
	restore(cont)
	go_to(reg(cont))
}

func compound_apply() {
	assign(unev, procedure_parameters(reg(proc)))
	assign(env, procedure_environment(reg(proc)))
	assign(env, extend_environment(reg(unev), reg(argl), reg(env)))
	assign(unev, procedure_body(reg(proc)))
	go_to(label(ev_sequence))
}

func ev_begin() {
	assign(unev, begin_actions(reg(exp)))
	save(*cont)
	go_to(label(ev_sequence))
}

func ev_sequence() {
	assign(exp, first_exp(reg(unev)))
	if test(is_last_exp(reg(unev))) {
		ev_sequence_last_exp()
		return
	}
	save(*unev)
	save(*env)
	assign(cont, label(ev_sequence_continue))
	go_to(label(eval_dispatch))
}

func ev_sequence_continue() {
	restore(env)
	restore(unev)
	assign(unev, rest_exps(reg(unev)))
	go_to(label(ev_sequence))
}

func ev_sequence_last_exp() {
	restore(cont)
	go_to(label(eval_dispatch))
}

func ev_if() {
	save(*exp)
	save(*env)
	save(*cont)
	assign(cont, label(ev_if_decide))
	assign(exp, if_predicate(reg(exp)))
	go_to(label(eval_dispatch))
}

func ev_if_decide() {
	restore(cont)
	restore(env)
	restore(exp)
	if test(is_true(reg(val))) {
		ev_if_consequent()
		return
	}
	go_to(label(ev_if_alternative))
}

func ev_if_alternative() {
	assign(exp, if_alternative(reg(exp)))
	go_to(label(eval_dispatch))
}

func ev_if_consequent() {
	assign(exp, if_consequent(reg(exp)))
	go_to(label(eval_dispatch))
}

func ev_assignment() {
	assign(unev, assignment_variable(reg(exp)))
	save(*unev)
	assign(exp, assignment_value(reg(exp)))
	save(*env)
	save(*cont)
	assign(cont, label(ev_assignment_1))
	go_to(label(eval_dispatch))
}

func ev_assignment_1() {
	restore(cont)
	restore(env)
	restore(unev)
	set_variable_value(reg(unev), reg(val), reg(env))
	assign(val, constant("ok"))
	go_to(reg(cont))
}

func ev_definition() {
	assign(unev, definition_variable(reg(exp)))
	save(*unev)
	assign(exp, definition_value(reg(exp)))
	save(*env)
	save(*cont)
	assign(cont, label(ev_definition_1))
	go_to(label(eval_dispatch))
}

func ev_definition_1() {
	restore(cont)
	restore(env)
	restore(unev)
	define_variable(reg(unev), reg(val), reg(env))
	assign(val, constant("ok"))
	go_to(reg(cont))
}

func unknown_expression_type() {
	assign(val, constant("unknown expression type error"))
	go_to(label(signal_error))
}

func unknown_procedure_type() {
	restore(cont)
	assign(val, constant("unknown procedure type error"))
	go_to(label(signal_error))
}

func signal_error() {
	user_print(reg(val))
	go_to(label(done))
}

func ev_self_eval() {
	assign(val, reg(exp))
	go_to(reg(cont))
}

func ev_variable() {
	assign(val, lookup_variable_value(reg(exp), reg(env)))
	go_to(reg(cont))
}

func ev_quoted() {
	assign(val, text_of_quotation(reg(exp)))
	go_to(reg(cont))
}

func done() {
	// nothing to do
}

func print_result() {
	user_print(reg(val))
	go_to(label(done))
}

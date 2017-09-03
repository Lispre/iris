// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package runtime

import (
	"math"

	"github.com/ta2gch/iris/runtime/env"
	"github.com/ta2gch/iris/runtime/ilos"
	"github.com/ta2gch/iris/runtime/ilos/class"
	"github.com/ta2gch/iris/runtime/ilos/instance"
)

// Listp returns t if obj is a list (instance of class list); otherwise, returns nil.
// obj may be any ISLISP object.
func Listp(e env.Environment, obj ilos.Instance) (ilos.Instance, ilos.Instance) {
	if ilos.InstanceOf(class.Cons, obj) {
		return T, nil
	}
	return Nil, nil
}

// CreateList returns a list of length i. If initial-element is given, the elements of the new list
// are initialized with this object; otherwise, the initialization is implementation defined. An
// error shall be signaled if the requested list cannot be allocated (error-id. cannot-create-list).
// An error shall be signaled if i is not a non-negative integer (error-id. domain-error).
//initial-element may be any ISLISP object.
func CreateList(e env.Environment, i ilos.Instance, initialElement ...ilos.Instance) (ilos.Instance, ilos.Instance) {
	if err := ensure(e, class.Integer, i); err != nil {
		return nil, err
	}
	if len(initialElement) > 1 {
		return SignalCondition(e, instance.NewArityError(e), Nil)
	}
	elm := Nil
	if len(initialElement) == 1 {
		elm = initialElement[0]
	}
	cons := Nil
	for j := 0; j < int(i.(instance.Integer)); j++ {
		cons = instance.NewCons(elm, cons)
	}
	return cons, nil
}

// List returns a new list whose length is the number of arguments and whose elements are the
// arguments in the same order as in the list-form. An error shall be signaled if the requested list
// cannot be allocated (error-id. cannot-create-list). Each obj may be any ISLISP object.
func List(e env.Environment, objs ...ilos.Instance) (ilos.Instance, ilos.Instance) {
	cons := Nil
	for i := len(objs) - 1; i >= 0; i-- {
		cons = instance.NewCons(objs[i], cons)
	}
	return cons, nil
}

// Reverse returns a list whose elements are those of the given list, but in reverse
// order. An error shall be signaled if list is not a list (error-id. domain-error).
//
// For reverse, no side-effect to the given list occurs. The resulting list is permitted but not
// required to share structure with the input list.
func Reverse(e env.Environment, list ilos.Instance) (ilos.Instance, ilos.Instance) {
	if err := ensure(e, class.List, list); err != nil {
		return nil, err
	}
	cons := Nil
	for _, car := range list.(instance.List).Slice() {
		cons = instance.NewCons(car, cons)
	}
	return cons, nil
}

// Nreverse returns a list whose elements are those of the given list, but in reverse
// order. An error shall be signaled if list is not a list (error-id. domain-error).
//
// For nreverse, the conses which make up the top level of the given list are permitted, but not
// required, to be side-effected in order to produce this new list. nreverse should never be called
// on a literal object.
func Nreverse(e env.Environment, list ilos.Instance) (ilos.Instance, ilos.Instance) {
	// TODO: tests literal object
	if err := ensure(e, class.List, list); err != nil {
		return nil, err
	}
	cons := Nil
	for _, car := range list.(instance.List).Slice() {
		cons = instance.NewCons(car, cons)
	}
	return cons, nil
}

// Append returns the result of appending all of the lists, or () if given no lists. An error shall
// be signaled if any list is not a list (error-id. domain-error).
//
// This function does not modify its arguments. It is implementation defined whether and when the
// result shares structure with its list arguments.
//
// An error shall be signaled if the list cannot be allocated (error-id. cannot-create-list).
func Append(e env.Environment, lists ...ilos.Instance) (ilos.Instance, ilos.Instance) {
	// Ref: https://github.com/sbcl/sbcl/blob/fe4faef65315c6ad52b3b89b62b6c6497cb78d09/src/code/list.lisp#L364

	result, err := List(e, Nil)
	if err != nil {
		return nil, err
	}
	cdr := result
	if err := ensure(e, class.List, lists...); err != nil {
		return nil, err
	}
	for _, list := range lists {
		for _, elt := range list.(instance.List).Slice() {
			it, err := List(e, elt)
			if err != nil {
				return nil, err
			}
			cdr.(*instance.Cons).Cdr = it
			cdr = cdr.(*instance.Cons).Cdr
		}
	}
	return result.(*instance.Cons).Cdr, nil
}

// Member returnes the first sublist of list whose car is obj  if list contains at least one
// occurrence of obj (as determined by eql).  Otherwise, nil is returned. An error shall be signaled
// if list is not a list (error-id. domain-error).
func Member(e env.Environment, obj, list ilos.Instance) (ilos.Instance, ilos.Instance) {
	if err := ensure(e, class.List, list); err != nil {
		return nil, err
	}
	for idx, elt := range list.(instance.List).Slice() {
		if obj == elt { // eql
			return list.(instance.List).NthCdr(idx), nil
		}
	}
	return Nil, nil
}

// Mapcar operates on successive elements of the lists. function is applied to the first element of
// each list, then to the second element of each list, and so on. The iteration terminates when the
// shortest list runs out, and excess elements in other lists are ignored. The value returned by
// mapcar is a list of the results of successive calls to function.
func Mapcar(e env.Environment, function, list1 ilos.Instance, lists ...ilos.Instance) (ilos.Instance, ilos.Instance) {
	lists = append([]ilos.Instance{list1}, lists...)
	if err := ensure(e, class.Function, function); err != nil {
		return nil, err
	}
	if err := ensure(e, class.List, lists...); err != nil {
		return nil, err
	}
	max := 0.0
	for _, list := range lists {
		max = math.Max(max, float64(len(list.(instance.List).Slice())))
	}
	result := []ilos.Instance{}
	for i := 0; i < int(max); i++ {
		arguments := make([]ilos.Instance, len(lists))
		for j, list := range lists {
			arguments[j] = list.(instance.List).Nth(i)
		}
		ret, err := function.(instance.Applicable).Apply(e.NewDynamic(), arguments...)
		if err != nil {
			return nil, err
		}
		result = append(result, ret)
	}
	return List(e, result...)
}

// Mapc is like mapcar except that the results of applying function are not accumulated;
// list1 is returned.
func Mapc(e env.Environment, function, list1 ilos.Instance, lists ...ilos.Instance) (ilos.Instance, ilos.Instance) {
	lists = append([]ilos.Instance{list1}, lists...)
	if err := ensure(e, class.Function, function); err != nil {
		return nil, err
	}
	if err := ensure(e, class.List, lists...); err != nil {
		return nil, err
	}
	max := 0.0
	for _, list := range lists {
		max = math.Max(max, float64(len(list.(instance.List).Slice())))
	}
	for i := 0; i < int(max); i++ {
		arguments := make([]ilos.Instance, len(lists))
		for j, list := range lists {
			arguments[j] = list.(instance.List).Nth(i)
		}
		if _, err := function.(instance.Applicable).Apply(e.NewDynamic(), arguments...); err != nil {
			return nil, err
		}
	}
	return list1, nil
}

// Mapcan is like mapcar respectively, except that the results of applying
// function are combined into a list by the use of an operation that performs a destructive form of
// append rather than list.
func Mapcan(e env.Environment, function, list1 ilos.Instance, lists ...ilos.Instance) (ilos.Instance, ilos.Instance) {
	lists = append([]ilos.Instance{list1}, lists...)
	if err := ensure(e, class.Function, function); err != nil {
		return nil, err
	}
	if err := ensure(e, class.List, lists...); err != nil {
		return nil, err
	}
	max := 0.0
	for _, list := range lists {
		max = math.Max(max, float64(len(list.(instance.List).Slice())))
	}
	result := []ilos.Instance{}
	for i := 0; i < int(max); i++ {
		arguments := make([]ilos.Instance, len(lists))
		for j, list := range lists {
			arguments[j] = list.(instance.List).Nth(i)
		}
		ret, err := function.(instance.Applicable).Apply(e.NewDynamic(), arguments...)
		if err != nil {
			return nil, err
		}
		result = append(result, ret)
	}
	return Append(e, result...)
}

// Maplist is like mapcar except that function is applied to successive sublists of the lists.
// function is first applied to the lists themselves, and then to the cdr of each list, and then to
// the cdr of the cdr of each list, and so on.
func Maplist(e env.Environment, function, list1 ilos.Instance, lists ...ilos.Instance) (ilos.Instance, ilos.Instance) {
	lists = append([]ilos.Instance{list1}, lists...)
	if err := ensure(e, class.Function, function); err != nil {
		return nil, err
	}
	if err := ensure(e, class.List, lists...); err != nil {
		return nil, err
	}
	max := 0.0
	for _, list := range lists {
		max = math.Max(max, float64(len(list.(instance.List).Slice())))
	}
	result := []ilos.Instance{}
	for i := 0; i < int(max); i++ {
		arguments := make([]ilos.Instance, len(lists))
		for j, list := range lists {
			arguments[j] = list.(instance.List).NthCdr(i)
		}
		ret, err := function.(instance.Applicable).Apply(e.NewDynamic(), arguments...)
		if err != nil {
			return nil, err
		}
		result = append(result, ret)
	}
	return List(e, result...)
}

// Mapl is like maplist except that the results of applying function are not accumulated;
// list1 is returned.
func Mapl(e env.Environment, function, list1 ilos.Instance, lists ...ilos.Instance) (ilos.Instance, ilos.Instance) {
	lists = append([]ilos.Instance{list1}, lists...)
	if err := ensure(e, class.Function, function); err != nil {
		return nil, err
	}
	if err := ensure(e, class.List, lists...); err != nil {
		return nil, err
	}
	max := 0.0
	for _, list := range lists {
		max = math.Max(max, float64(len(list.(instance.List).Slice())))
	}
	for i := 0; i < int(max); i++ {
		arguments := make([]ilos.Instance, len(lists))
		for j, list := range lists {
			arguments[j] = list.(instance.List).NthCdr(i)
		}
		if _, err := function.(instance.Applicable).Apply(e.NewDynamic(), arguments...); err != nil {
			return nil, err
		}
	}
	return list1, nil
}

// Mapcon is like maplist respectively, except that the results of applying
// function are combined into a list by the use of an operation that performs a destructive form of
// append rather than list.
func Mapcon(e env.Environment, function, list1 ilos.Instance, lists ...ilos.Instance) (ilos.Instance, ilos.Instance) {
	lists = append([]ilos.Instance{list1}, lists...)
	if err := ensure(e, class.Function, function); err != nil {
		return nil, err
	}
	if err := ensure(e, class.List, lists...); err != nil {
		return nil, err
	}
	max := 0.0
	for _, list := range lists {
		max = math.Max(max, float64(len(list.(instance.List).Slice())))
	}
	result := []ilos.Instance{}
	for i := 0; i < int(max); i++ {
		arguments := make([]ilos.Instance, len(lists))
		for j, list := range lists {
			arguments[j] = list.(instance.List).NthCdr(i)
		}
		ret, err := function.(instance.Applicable).Apply(e.NewDynamic(), arguments...)
		if err != nil {
			return nil, err
		}
		result = append(result, ret)
	}
	return Append(e, result...)
}

// Assoc returns the first cons if assocation-list contains at least one cons whose car is
// obj (as determined by eql). Otherwise, nil is returned. An error shall be signaled
// if association-list is not a list of conses (error-id. domain-error).
func Assoc(e env.Environment, obj, associationList ilos.Instance) (ilos.Instance, ilos.Instance) {
	if err := ensure(e, class.List, associationList); err != nil {
		return nil, err
	}
	for _, pair := range associationList.(instance.List).Slice() {
		if err := ensure(e, class.Cons, pair); err != nil {
			return nil, err
		}
		if pair.(*instance.Cons).Car == obj { // eql
			return pair.(*instance.Cons).Cdr, nil
		}
	}
	return Nil, nil
}

// Null returns t if obj is nil; otherwise, returns nil obj may be any ISLISP object.
func Null(e env.Environment, obj ilos.Instance) (ilos.Instance, ilos.Instance) {
	if obj == Nil {
		return T, nil
	}
	return Nil, nil
}
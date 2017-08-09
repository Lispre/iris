package instance

import (
	"github.com/k0kubun/pp"
	"github.com/ta2gch/iris/runtime/ilos"
	"github.com/ta2gch/iris/runtime/ilos/class"
)

//
// General Array Star
//

type GeneralArrayStar struct {
	dimension [128]int
	array     map[[128]int]ilos.Instance
}

func NewGeneralArrayStar(key [128]int) ilos.Instance {
	return &GeneralArrayStar{key, map[[128]int]ilos.Instance{}}
}

func (*GeneralArrayStar) Class() ilos.Class {
	return class.GeneralArrayStar
}

func (a *GeneralArrayStar) GetSlotValue(key ilos.Instance, _ ilos.Class) (ilos.Instance, bool) {
	if symbol, ok := key.(Symbol); ok && symbol == "LENGTH" {
		cons := NewNull()
		for i := 128; i > 0; i-- {
			if a.dimension[i-1] != 0 {
				cons = NewCons(NewInteger(a.dimension[i-1]), cons)
			}
		}
		return cons, true
	}
	if Of(class.List, key) {
		dim := [128]int{}
		idx := 0
		cdr := key
		for Of(class.Cons, cdr) {
			dim[idx] = int(UnsafeCar(cdr).(Integer))
			cdr = UnsafeCdr(cdr)
			idx++
		}
		for i := 0; i < 128; i++ {
			if a.dimension[i] != 0 && dim[i] >= a.dimension[i] {
				return nil, false
			}
		}
		return a.array[dim], true
	}
	return nil, false
}

func (a *GeneralArrayStar) SetSlotValue(key ilos.Instance, value ilos.Instance, _ ilos.Class) bool {
	if Of(class.List, key) {
		dim := [128]int{}
		idx := 0
		cdr := key
		for Of(class.Cons, cdr) {
			dim[idx] = int(UnsafeCar(cdr).(Integer))
			cdr = UnsafeCdr(cdr)
			idx++
		}
		for i := 0; i < 128; i++ {
			if a.dimension[i] != 0 && dim[i] >= a.dimension[i] {
				return false
			}
		}
		a.array[dim] = value
		return true
	}
	return false
}

func (a *GeneralArrayStar) String() string {
	return pp.Sprint(a.array)
}

//
// General Vector
//

type GeneralVector []ilos.Instance

func NewGeneralVector(n int) ilos.Instance {
	return GeneralVector(make([]ilos.Instance, n))
}

func (GeneralVector) Class() ilos.Class {
	return class.GeneraVector
}

func (i GeneralVector) GetSlotValue(key ilos.Instance, _ ilos.Class) (ilos.Instance, bool) {
	if symbol, ok := key.(Symbol); ok && symbol == "LENGTH" {
		return NewInteger(len(i)), true
	}
	if index, ok := key.(Integer); ok && int(index) < len(i) {
		return i[int(index)], true
	}
	return nil, false
}

func (i GeneralVector) SetSlotValue(key ilos.Instance, value ilos.Instance, _ ilos.Class) bool {
	if index, ok := key.(Integer); ok && int(index) < len(i) {
		i[int(index)] = value
		return true
	}
	return false
}

func (i GeneralVector) String() string {
	return pp.Sprint([]ilos.Instance(i))
}

//
// String
//

type String []rune

func NewString(a string) ilos.Instance {
	return String([]rune(a))
}

func (String) Class() ilos.Class {
	return class.String
}

func (i String) GetSlotValue(key ilos.Instance, _ ilos.Class) (ilos.Instance, bool) {
	if symbol, ok := key.(Symbol); ok && symbol == "LENGTH" {
		return NewInteger(len(i)), true
	}
	if index, ok := key.(Integer); ok && int(index) < len(i) {
		return NewCharacter(i[int(index)]), true
	}
	return nil, false
}

func (i String) SetSlotValue(key ilos.Instance, value ilos.Instance, _ ilos.Class) bool {
	if index, ok := key.(Integer); ok && int(index) < len(i) {
		if character, ok := value.(Character); ok {
			i[index] = rune(character)
			return true
		}
	}
	return false
}

func (i String) String() string {
	return "\"" + string(i) + "\""
}

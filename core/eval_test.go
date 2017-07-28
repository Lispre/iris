package core

import (
	"reflect"
	"strings"
	"testing"

	"github.com/ta2gch/gazelle/core/class"
	"github.com/ta2gch/gazelle/core/class/cons"
	"github.com/ta2gch/gazelle/core/class/function"
	"github.com/ta2gch/gazelle/core/environment"
	"github.com/ta2gch/gazelle/reader/parser"
	"github.com/ta2gch/gazelle/reader/tokenizer"
)

func read(s string) *class.Instance {
	e, _ := parser.Parse(tokenizer.New(strings.NewReader(s)))
	return e
}

func TestEval(t *testing.T) {
	testEnv := environment.New()
	testEnv.Var["PI"] = class.Float.New(3.14)
	testEnv.Fun["INC"] = function.New(func(args *class.Instance, global *environment.Env) (*class.Instance, error) {
		car, _ := cons.Car(args)
		return class.Integer.New(car.Value().(int) + 1), nil
	})
	type args struct {
		obj    *class.Instance
		local  *environment.Env
		global *environment.Env
	}
	tests := []struct {
		name    string
		args    args
		want    *class.Instance
		wantErr bool
	}{
		{
			name:    "local variable",
			args:    args{class.Symbol.New("PI"), testEnv, nil},
			want:    class.Float.New(3.14),
			wantErr: false,
		},
		{
			name:    "local function",
			args:    args{read("(inc (inc 1))"), testEnv, nil},
			want:    class.Integer.New(3),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Eval(tt.args.obj, tt.args.local, tt.args.global)
			if (err != nil) != tt.wantErr {
				t.Errorf("Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}

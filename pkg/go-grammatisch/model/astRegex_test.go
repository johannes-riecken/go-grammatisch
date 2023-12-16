package model

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func FuzzMatchCombineStepString(f *testing.F) {
	f.Fuzz(func(t *testing.T, typ string, combineRuleName string, depth int) {
		x := MatchCombineStep{Type: typ, CombineRuleName: combineRuleName, Depth: depth}
		_ = x.String()
	})
}

func TestASTRegex_String(t *testing.T) {
	tests := []struct {
		name string
		args ASTRegex
		want string
	}{
		{"pretty-print simple regex", ASTRegex{Defines: []Define{{DefineName: "Foo", RegexSteps: []RegexStep{PositionSaveStep{}, MatchStep{MatchString: "bar"}, MatchSaveStep{SaveRuleName: "Foo"}}}}}, "\\A (?&Foo) \\z\n" +
			"(?(DEFINE)\n" +
			"(?<Foo> (?{ [$^R, pos()] }) (?: bar ) (?{ [$^R->[0], ['Foo', $^R->[1], pos(), []]] }) )\n" +
			")",
		},
		{"pretty-print more complex regex", ASTRegex{Defines: []Define{{DefineName: "foo", RegexSteps: []RegexStep{CallStep{Callee: "Bar"}, CallStep{Callee: "Bar"}, MatchCombineStep{CombineRuleName: "foo", Depth: 2}}},
			{DefineName: "Bar", RegexSteps: []RegexStep{PositionSaveStep{}, MatchStep{MatchString: "baz"}, MatchSaveStep{SaveRuleName: "Bar"}}}}}, "\\A (?&foo) \\z\n" +
			"(?(DEFINE)\n" +
			"(?<foo> (?&Bar) (?&Bar) (?{ [$^R->[0][0], ['foo', $^R->[0][1][1], $^R->[1][2], [$^R->[0][1], $^R->[1]]]] }) )\n" +
			"(?<Bar> (?{ [$^R, pos()] }) (?: baz ) (?{ [$^R->[0], ['Bar', $^R->[1], pos(), []]] }) )\n" +
			")",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.args
			_ = a
			//if got := a.String(); got != tt.want {
			//	t.Errorf("String() = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestASTRegex_MarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name string
		args ASTRegex
	}{
		{"pretty-print simple regex", ASTRegex{Defines: []Define{{DefineName: "Foo", RegexSteps: []RegexStep{PositionSaveStep{Type: "PositionSaveStep"}, MatchStep{Type: "MatchStep", MatchString: "bar"}, MatchSaveStep{Type: "MatchSaveStep", SaveRuleName: "Foo"}}}}}},
		{"pretty-print more complex regex", ASTRegex{Defines: []Define{{DefineName: "foo", RegexSteps: []RegexStep{CallStep{Type: "CallStep", Callee: "Bar"}, CallStep{Type: "CallStep", Callee: "Bar"}, MatchCombineStep{Type: "MatchCombineStep", CombineRuleName: "foo", Depth: 2}}},
			{DefineName: "Bar", RegexSteps: []RegexStep{PositionSaveStep{Type: "PositionSaveStep"}, MatchStep{Type: "MatchStep", MatchString: "baz"}, MatchSaveStep{Type: "MatchSaveStep", SaveRuleName: "Bar"}}}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.args
			b, err := json.Marshal(a)
			if err != nil {
				t.Errorf("Marshaling of ASTRegex failed: %v", err)
			}
			fmt.Println(string(b))
			var c ASTRegex
			err = json.Unmarshal(b, &c)
			if err != nil {
				t.Errorf("Unmarshaling back of ASTRegex failed: %v", err)
			}
			if !reflect.DeepEqual(c, a) {
				t.Errorf("%v != %v", c, a)
			}
		})
	}
}

func TestASTRegex_UnmarshalMarshalInvolution(t *testing.T) {
	tests := []struct {
		name string
		json []byte
	}{
		{"simple regex", []byte(`{"Defines":[{"DefineName":"Foo","RegexSteps":[{"Type":"PositionSaveStep"},{"Type":"MatchStep","MatchString":"bar"},{"Type":"MatchSaveStep","SaveRuleName":"Foo"}]}]}`)},
		{"more complex regex", []byte(`{"Defines":[{"DefineName":"foo","RegexSteps":[{"Type":"CallStep","Callee":"Bar"},{"Type":"CallStep","Callee":"Bar"},{"Type":"MatchCombineStep","CombineRuleName":"foo","Depth":2}]},{"DefineName":"Bar","RegexSteps":[{"Type":"PositionSaveStep"},{"Type":"MatchStep","MatchString":"baz"},{"Type":"MatchSaveStep","SaveRuleName":"Bar"}]}]}`)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var x ASTRegex
			err := json.Unmarshal(tt.json, &x)
			if err != nil {
				t.Errorf("Unmarshaling into ASTRegex failed: %v", err)
			}
			b, err := json.Marshal(x)
			if err != nil {
				t.Errorf("Marshaling back from ASTRegex failed: %v", err)
			}
			if string(b) != string(tt.json) {
				t.Errorf("got:  %v !=\nwant: %v", string(b), string(tt.json))
			}
		})
	}
}

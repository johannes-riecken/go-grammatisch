package model

import (
	"testing"
	"encoding/json"
)

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
		{"pretty-print more complex regex", ASTRegex{Defines: []Define{{DefineName: "foo", RegexSteps: []RegexStep{CallStep{callee: "Bar"}, CallStep{callee: "Bar"}, MatchCombineStep{combineRuleName: "foo", depth: 2}}},
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
			if got := a.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestASTRegex_MarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name string
		args ASTRegex
	}{
		{"pretty-print simple regex", ASTRegex{Defines: []Define{{DefineName: "Foo", RegexSteps: []RegexStep{PositionSaveStep{}, MatchStep{MatchString: "bar"}, MatchSaveStep{SaveRuleName: "Foo"}}}}}},
		{"pretty-print more complex regex", ASTRegex{Defines: []Define{{DefineName: "foo", RegexSteps: []RegexStep{CallStep{callee: "Bar"}, CallStep{callee: "Bar"}, MatchCombineStep{combineRuleName: "foo", depth: 2}}},
		{DefineName: "Bar", RegexSteps: []RegexStep{PositionSaveStep{}, MatchStep{MatchString: "baz"}, MatchSaveStep{SaveRuleName: "Bar"}}}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.args
			b, err := json.Marshal(a)
			if err != nil {
				t.Errorf("Marshaling of ASTRegex failed: %v", err)
			}
			var c ASTRegex
			err = json.Unmarshal(b, &c)
			if err != nil {
				t.Errorf("Unmarshaling back of ASTRegex failed: %v", err)
			}
		})
	}
}

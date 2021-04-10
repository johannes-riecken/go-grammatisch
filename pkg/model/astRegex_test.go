package model

import "testing"

func TestASTRegex_String(t *testing.T) {
	tests := []struct {
		name   string
		args ASTRegex
		want   string
	}{
		{"pretty-print simple regex", ASTRegex{[]Define{{"Foo", []RegexStep{PositionSaveStep{},MatchStep{"bar"},MatchSaveStep{"Foo"}}}}}, "\\A (?&Foo) \\z\n" +
					"(?(DEFINE)\n" +
					"(?<Foo> (?{ [$^R, pos()] }) (?: bar ) (?{ [$^R->[0], ['Foo', $^R->[1], pos(), []]] }) )\n" +
					")",
		},
		{"prety-print more complex regex", ASTRegex{[]Define{{"foo",[]RegexStep{CallStep{"Bar"},CallStep{"Bar"},MatchCombineStep{"foo",2}}},
			{"Bar",[]RegexStep{PositionSaveStep{},MatchStep{"baz"},MatchSaveStep{"Bar"}}}}}, "\\A (?&foo) \\z\n" +
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

package model

import (
	"reflect"
	"testing"
)

func TestGrammar_ToRegex(t *testing.T) {
	tests := []struct {
		name string
		args Grammar
		want ASTRegex
	}{
		{name: "convert simple grammar", args: 		Grammar{
			RuleSpecs: []RuleSpec{
				{
					RuleRef: "Foo",
					Alternatives: []Alternative{Alternative{
						Elements: []Element{
							Quoted{
								Quoted: "'bar'",
							},
						},
					}},
				},
			},
		},
		want: ASTRegex{[]Define{{"Foo", []RegexStep{PositionSaveStep{}, MatchStep{"bar"}, MatchSaveStep{"Foo"}}}}},
			},
			{name: "convert more complex grammar", args: Grammar{
			[]RuleSpec{{"foo",[]Alternative{{[]Element{RuleRef{"Bar"}, RuleRef{"Bar"}}}}},{"Bar",[]Alternative{{[]Element{Quoted{"'baz'"}}}}},
			/*
			Just (Grammar [RuleSpec "foo" [Alternative [RuleRef "Bar" Nothing, RuleRef "Bar" Nothing]], RuleSpec "Bar" [Alternative [Quoted "'baz'" Nothing]]])
			 */
			}}, want: ASTRegex{[]Define{{"foo",[]RegexStep{CallStep{"Bar"}, CallStep{"Bar"}, MatchCombineStep{"foo",2}}},
				{"Bar",[]RegexStep{PositionSaveStep{}, MatchStep{"baz"}, MatchSaveStep{"Bar"}}}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := tt.args
			if got := g.ToRegex(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToRegex() = %v, want %v", got, tt.want)
			}
		})
	}
}

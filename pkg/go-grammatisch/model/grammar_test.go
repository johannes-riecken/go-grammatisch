package model

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestGrammar_ToRegex(t *testing.T) {
	tests := []struct {
		name string
		args Grammar
		want ASTRegex
	}{
		{name: "convert simple grammar", args: Grammar{
			RuleSpecs: []RuleSpec{
				{
					RuleRef: "Foo",
					Alternatives: []Alternative{{
						Elements: []Element{
							Quoted{
								Quoted: "'bar'",
							},
						},
					}},
				},
			},
		},
		want: ASTRegex{Defines: []Define{{DefineName: "Foo", RegexSteps: []RegexStep{PositionSaveStep{}, MatchStep{MatchString: "bar"}, MatchSaveStep{SaveRuleName: "Foo"}}}}},
		},
		{name: "convert more complex grammar", args: Grammar{
			RuleSpecs: []RuleSpec{{RuleRef: "foo", Alternatives: []Alternative{{Elements: []Element{RuleRef{RuleRefName: "Bar"}, RuleRef{RuleRefName: "Bar"}}}}}, {RuleRef: "Bar", Alternatives: []Alternative{{Elements: []Element{Quoted{Quoted: "'baz'"}}}}}}},/*
			Just (Grammar [RuleSpec "foo" [Alternative [RuleRef "Bar" Nothing, RuleRef "Bar" Nothing]], RuleSpec "Bar" [Alternative [Quoted "'baz'" Nothing]]])
		*/
		want: ASTRegex{Defines: []Define{{DefineName: "foo", RegexSteps: []RegexStep{CallStep{callee: "Bar"}, CallStep{callee: "Bar"}, MatchCombineStep{combineRuleName: "foo", depth: 2}}},
		{DefineName: "Bar", RegexSteps: []RegexStep{PositionSaveStep{}, MatchStep{MatchString: "baz"}, MatchSaveStep{SaveRuleName: "Bar"}}}}}},
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

func TestGrammar_MarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name string
		args Grammar
	}{
		{name: "convert simple grammar", args: Grammar{
			RuleSpecs: []RuleSpec{
				{
					RuleRef: "Foo",
					Alternatives: []Alternative{{
						Elements: []Element{
							Quoted{
								Quoted: "'bar'",
							},
						},
					}},
				},
			},
		},
		},
		{name: "convert more complex grammar", args: Grammar{
			RuleSpecs: []RuleSpec{{RuleRef: "foo", Alternatives: []Alternative{{Elements: []Element{RuleRef{RuleRefName: "Bar"}, RuleRef{RuleRefName: "Bar"}}}}}, {RuleRef: "Bar", Alternatives: []Alternative{{Elements: []Element{Quoted{Quoted: "'baz'"}}}}}}},/*
			Just (Grammar [RuleSpec "foo" [Alternative [RuleRef "Bar" Nothing, RuleRef "Bar" Nothing]], RuleSpec "Bar" [Alternative [Quoted "'baz'" Nothing]]])
		*/
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := tt.args
			b, err := json.Marshal(g)
			if err != nil {
				t.Errorf("Marshaling of Grammar failed: %v", err)
			}
			fmt.Println(string(b))
			var c Grammar
			err = json.Unmarshal(b, &c)
			if err != nil {
				t.Errorf("Unmarshaling back of Grammar failed: %v", err)
			}
		})
	}
}

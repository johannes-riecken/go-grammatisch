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
		want: ASTRegex{Defines: []Define{{DefineName: "foo", RegexSteps: []RegexStep{CallStep{Callee: "Bar"}, CallStep{Callee: "Bar"}, MatchCombineStep{CombineRuleName: "foo", Depth: 2}}},
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
								Type: "Quoted",
								Quoted: "'bar'",
							},
						},
					}},
				},
			},
		},
		},
		{name: "convert more complex grammar", args: Grammar{
			RuleSpecs: []RuleSpec{{RuleRef: "foo", Alternatives: []Alternative{{Elements: []Element{
				RuleRef{Type: "RuleRef", RuleRefName: "Bar"},
				RuleRef{Type: "RuleRef", RuleRefName: "Bar"},
			}}}}, {RuleRef: "Bar", Alternatives: []Alternative{{Elements: []Element{Quoted{Type: "Quoted", Quoted: "'baz'"}}}}}}},/*
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
			if !reflect.DeepEqual(g, c) {
				t.Errorf("%+v != %+v", g, c)
			}
		})
	}
}

func TestGrammar_UnmarshalMarshalInvolution(t *testing.T) {
	tests := []struct {
		name string
		json []byte
	}{
		{"simple grammar", []byte(`{"RuleSpecs":[{"RuleRef":"Foo","Alternatives":[{"Elements":[{"Type":"Quoted","Quoted":"'bar'"}]}]}]}`)},
		{"more complex grammar", []byte(`{"RuleSpecs":[{"RuleRef":"foo","Alternatives":[{"Elements":[{"Type":"RuleRef","RuleRefName":"Bar"},{"Type":"RuleRef","RuleRefName":"Bar"}]}]},{"RuleRef":"Bar","Alternatives":[{"Elements":[{"Type":"Quoted","Quoted":"'baz'"}]}]}]}`)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var x Grammar
			err := json.Unmarshal(tt.json, &x)
			if err != nil {
				t.Errorf("Unmarshaling into Grammar failed: %v", err)
			}
			b, err := json.Marshal(x)
			if err != nil {
				t.Errorf("Marshaling back from Grammar failed: %v", err)
			}
			if string(b) != string(tt.json) {
				t.Errorf("got:  %v !=\nwant: %v", string(b), string(tt.json))
			}
		})
	}
}

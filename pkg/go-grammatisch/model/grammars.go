package model

import "fmt"

type Grammars string

func (g Grammars) Process() (Grammar, error) {
	switch string(g) {
	case "calculator":
		return Grammar{
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
		}, nil
	case "simple":
		return Grammar{
			RuleSpecs: []RuleSpec{
				{RuleRef: "foo", Alternatives: []Alternative{{Elements: []Element{RuleRef{RuleRefName: "Bar"}, RuleRef{RuleRefName: "Bar"}}}}},
				{RuleRef: "Bar", Alternatives: []Alternative{{Elements: []Element{Quoted{Quoted: "'baz'"}}}}},
			},
		}, nil
	default:
		return Grammar{}, fmt.Errorf("unknown grammar choice: %v", g)
	}
}

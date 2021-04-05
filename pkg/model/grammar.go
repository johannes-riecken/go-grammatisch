package model

type Grammar struct {
	RuleSpecs []RuleSpec
}

type RuleSpec struct {
	RuleRef string
	Alternatives []Alternative
}

type Alternative struct {
	Elements []Element
}

type Element interface {
	ElementMarker()
}

type RuleRef struct {
	ruleRefName string
}

func (RuleRef) ElementMarker() {}

type Quoted struct {
	Quoted string
}

func (Quoted) ElementMarker() {}

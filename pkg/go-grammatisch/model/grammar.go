package model

import (
	"log"
)

type Grammar struct {
	RuleSpecs []RuleSpec
}

func (g *Grammar) Process() ASTRegex {
	defines := make([]Define, len(g.RuleSpecs))
	for i, x := range g.RuleSpecs {
		defines[i] = x.ToRegex()
	}
	return ASTRegex{Defines: defines}
}

type RuleSpec struct {
	RuleRef      string
	Alternatives []Alternative
}

func (r *RuleSpec) ToRegex() Define {
	ret := Define{}
	ret.DefineName = r.RuleRef
	ret.RegexSteps = make([]RegexStep, len(r.Alternatives[0].Elements))
	for i, x := range r.Alternatives[0].Elements {
		ret.RegexSteps[i] = x.ToRegex()
	}
	ret.RegexSteps = altToRegexPost(r.RuleRef, ret.RegexSteps)
	return ret
}

func altToRegexPost(name string, regexSteps []RegexStep) []RegexStep {
	return appendCtorStep(name, insertPosSaveStep(isToken(name), regexSteps))
}

func insertPosSaveStep(isToken bool, regexSteps []RegexStep) []RegexStep {
	if isToken {
		return append([]RegexStep{PositionSaveStep{}}, regexSteps...)
	}
	return regexSteps
}

func appendCtorStep(name string, regexSteps []RegexStep) []RegexStep {
	callStepCount := 0
	for _, x := range regexSteps {
		if _, ok := x.(CallStep); ok {
			callStepCount++
		}
	}
	if callStepCount > 0 {
		return append(regexSteps, MatchCombineStep{CombineRuleName: name, Depth: callStepCount})
	} else {
		return append(regexSteps, MatchSaveStep{SaveRuleName: name})
	}
}

type Alternative struct {
	Elements []Element
}

func (a *Alternative) ToRegex() []RegexStep {
	altRegexes := make([]RegexStep, len(a.Elements))
	for i, x := range a.Elements {
		altRegexes[i] = x.ToRegex()
	}
	return altRegexes
}

type Element interface {
	ElementMarker()
	ToRegex() RegexStep
}

func (q Quoted) ToRegex() RegexStep {
	return MatchStep{MatchString: unquote(q.Quoted)}
}

func (r RuleRef) ToRegex() RegexStep {
	return CallStep{Callee: r.RuleRefName}
}

func isToken(s string) bool {
	if s == "" {
		log.Fatal("isToken failed for empty string")
	}
	return s[0] >= 'A' && s[0] <= 'Z'
}

func unquote(s string) string {
	return s[1 : len(s)-1]
}

type RuleRef struct {
	Type string
	RuleRefName string
}

func (RuleRef) ElementMarker() {}

type Quoted struct {
	Type string
	Quoted string
}

func (Quoted) ElementMarker() {}

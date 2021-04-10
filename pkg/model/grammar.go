package model

import (
	"log"
)

type Grammar struct {
	ruleSpecs []RuleSpec
}

func (g *Grammar) ToRegex() ASTRegex {
	defines := make([]Define, len(g.ruleSpecs))
	for i, x := range g.ruleSpecs {
		defines[i] = x.ToRegex()
	}
	return ASTRegex{Defines: defines}
}

type RuleSpec struct {
	ruleRef string
	alternatives []Alternative
}

func (r *RuleSpec) ToRegex() Define {
	ret := Define{}
	ret.defineName = r.ruleRef
	ret.regexSteps = make([]RegexStep, len(r.alternatives))
	for i, x := range r.alternatives[0].elements {
		ret.regexSteps[i] = x.ToRegex()
	}
	ret.regexSteps = altToRegexPost(r.ruleRef, ret.regexSteps)
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
		return append(regexSteps, MatchCombineStep{combineRuleName: name, depth: callStepCount})
	} else {
		return append(regexSteps, MatchSaveStep{saveRuleName: name})
	}
}

type Alternative struct {
	elements []Element
}

func (a *Alternative) ToRegex() []RegexStep {
	altRegexes := make([]RegexStep, len(a.elements))
	for i, x := range a.elements {
		altRegexes[i] = x.ToRegex()
	}
	return altRegexes
}

type Element interface {
	ElementMarker()
	ToRegex() RegexStep
}

func (q Quoted) ToRegex() RegexStep {
	return MatchStep{matchString: unquote(q.quoted)}
}

func (r RuleRef) ToRegex() RegexStep {
	return CallStep{callee: r.ruleRefName}
}

func isToken(s string) bool {
	if s == "" {
		log.Fatal("isToken failed for empty string")
	}
	return s[0] >= 'A' && s[0] <= 'Z'
}

func unquote(s string) string {
	return s[1:len(s) - 1]
}

type RuleRef struct {
	ruleRefName string
}

func (RuleRef) ElementMarker() {}

type Quoted struct {
	quoted string
}

func (Quoted) ElementMarker() {}

package model

type Grammar struct {
	RuleSpecs []RuleSpec
}

func (g *Grammar) ToRegex() ASTRegex {
	defines := make([]Define, len(g.RuleSpecs))
	for i, x := range g.RuleSpecs {
		defines[i] = x.ToRegex()
	}
	return ASTRegex{Defines: defines}
}

type RuleSpec struct {
	RuleRef string
	Alternatives []Alternative
}

func (r *RuleSpec) ToRegex() Define {
	ret := Define{}
	ret.DefineName = r.RuleRef
	ret.RegexSteps = make([]RegexStep, len(r.Alternatives))
	for i, x := r.Alternatives {
		ret.RegexSteps[i] = x.ToRegex(r.RuleRef)
	}
}

type Alternative struct {
	Elements []Element
}

func (a *Alternative) ToRegex(name string) {
	altRegexes := make([]RegexStep, len(a.Elements))
	for i, x := range a.Elements {
		altRegexes[i] = x.ToRegex()
	}
	return altToRegexPost(name, altRegexes)
}

func altToRegexPost(name, altRegexes) {
	return appenCtorStep(name, insertPosSaveStep(isToken(name)))
}

func insertPosSaveStep(isToken bool, regexSteps []RegexStep) []RegexStep {
	if isToken {
		return append([]RegexStep{PositionSaveStep}, regexSteps...)
	}
	return regexSteps
}

func appendCtorStep(name string, regexSteps []RegexStep) []RegexStep {


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

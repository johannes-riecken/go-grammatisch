package model

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type ASTRegex struct {
	Defines []Define
}

func (a ASTRegex) String() string {
	if len(a.Defines) == 0 {
		log.Fatal("ASTRegex must contain at least one Define")
	}
	buf := new(strings.Builder)
	buf.WriteString(fmt.Sprintf("\\A (?&%v) \\z\n(?(DEFINE)\n", a.Defines[0].DefineName))
	for _, x := range a.Defines {
		buf.WriteString(x.String())
		buf.WriteString("\n")
	}
	buf.WriteString(")")
	return buf.String()
}

func (a *ASTRegex) Match(inputDoc string) (string, error) {
	inputFile, err := os.Create("input.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create input file: %v", err)
	}
	defer inputFile.Close()
	_, err = inputFile.WriteString(inputDoc)
	if err != nil {
		return "", fmt.Errorf("failed to write to input file: %v", err)
	}
	regexFile, err := os.Create("regex.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create regex file: %v", err)
	}
	defer regexFile.Close()
	_, err = regexFile.WriteString(a.String())
	if err != nil {
		return "", fmt.Errorf("failed to write to regex file: %v", err)
	}
	ast, err := exec.Command("perl", "generateSyntaxTree.pl").Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute perl script: %v", err)
	}
	return string(ast), nil
}

type Define struct {
	DefineName string
	RegexSteps []RegexStep
}

func (d *Define) String() string {
	if len(d.RegexSteps) == 0 {
		log.Fatalf("Define must contain at least one RegexStep")
	}
	buf := new(strings.Builder)
	buf.WriteString("(?<")
	buf.WriteString(d.DefineName)
	buf.WriteString("> ")
	for _, x := range d.RegexSteps {
		buf.WriteString(x.String())
		buf.WriteString(" ")
	}
	buf.WriteString(")")
	return buf.String()
}

type RegexStep interface {
	RegexStepMarker()
	String() string
}

type PositionSaveStep struct {
	Type       string
}

func (PositionSaveStep) RegexStepMarker() {}

type CallStep struct {
	Type   string
	Callee string
}

func (CallStep) RegexStepMarker() {}

type MatchCombineStep struct {
	Type            string
	CombineRuleName string
	Depth           int
}

func (MatchCombineStep) RegexStepMarker() {}

type MatchSaveStep struct {
	Type       string
	SaveRuleName string
}

func (MatchSaveStep) RegexStepMarker() {}

type MatchStep struct {
	Type       string
	MatchString string
}

func (MatchStep) RegexStepMarker() {}

func (PositionSaveStep) String() string {
	return "(?{ [$^R, pos()] })"
}

func (c CallStep) String() string {
	return fmt.Sprintf("(?&%v)", c.Callee)
}

func (m MatchCombineStep) String() string {
	beginIdx := 1
	endIdx := 2
	arr := unfoldAnnotatedRegexTree(m.Depth).indicesList()
	i0 := indexOfR(arr[0])
	i1 := indexOfR(append(arr[1], beginIdx))
	i2 := indexOfR(append(arr[len(arr)-1], endIdx))
	children := new(strings.Builder)
	for i, x := range arr[1:] {
		children.WriteString(indexOfR(x))
		if i < len(arr[1:])-1 {
			children.WriteString(", ")
		}
	}
	return fmt.Sprintf("(?{ [%v, ['%v', %v, %v, [%v]]] })", i0, m.CombineRuleName, i1, i2, children.String())
}

func (m MatchStep) String() string {
	return fmt.Sprintf("(?: %v )", m.MatchString)
}

func (m MatchSaveStep) String() string {
	return fmt.Sprintf("(?{ [$^R->[0], ['%v', $^R->[1], pos(), []]] })", m.SaveRuleName)
}

func indexOfR(indexes []int) string {
	buf := new(strings.Builder)
	buf.WriteString("$^R->[")
	for i, x := range indexes {
		buf.WriteString(fmt.Sprintf("%v", x))
		if i != len(indexes)-1 {
			buf.WriteString("][")
		}
	}
	buf.WriteString("]")
	return buf.String()
}

func (t *tree) indicesList() [][]int {
	if t.left == nil && t.right == nil {
		return [][]int{t.val}
	}
	var ret [][]int
	if t.left != nil {
		ret = append(ret, t.left.indicesList()...)
	}
	if t.right != nil {
		ret = append(ret, t.right.indicesList()...)
	}
	return ret
}

func unfoldAnnotatedRegexTree(n int) *tree {
	var indexes []int
	res := tree{val: indexes}
	ref := &res
	for i := 0; i < n; i++ {
		ref.left = &tree{val: append(ref.val, 0)}
		ref.right = &tree{val: append(ref.val, 1)}
		ref = ref.left
	}
	return &res
}

type tree struct {
	val   []int
	left  *tree
	right *tree
}

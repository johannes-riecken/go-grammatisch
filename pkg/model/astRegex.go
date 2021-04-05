// package model
package main

import (
	"fmt"
	"strings"
	"log"
)

type ASTRegex struct {
	Defines []Define
}

func (a ASTRegex) String() string {
	if len(a.Defines) == 0 {
		log.Fatalf("ASTRegex must contain at least one Define: %v", a)
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

type Define struct {
	DefineName string
	RegexSteps []RegexStep
}

func (d *Define) String() string {
	if len(d.RegexSteps) == 0 {
		log.Fatalf("Define must contain at least one RegexStep: %v", d)
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
}

func (PositionSaveStep) RegexStepMarker() {}

type CallStep struct {
	Callee string
}

func (CallStep) RegexStepMarker() {}

type MatchCombineStep struct {
	CombineRuleName string
	Depth int
}

func (MatchCombineStep) RegexStepMarker() {}

type MatchSaveStep struct {
	SaveRuleName string
}

func (MatchSaveStep) RegexStepMarker() {}

type MatchStep struct {
	MatchString string
}

func (MatchStep) RegexStepMarker() {}

func (PositionSaveStep) String() string {
	return "(?{ [$^R, pos()] })"
}

func (c *CallStep) String() string {
	return fmt.Sprintf("(?&%v)", c.Callee)
}

func (m *MatchCombineStep) String() string {
	beginIdx := 1
	endIdx := 2
	arr := unfoldAnnotatedRegexTree(m.Depth).indicesList()
	i0 := indexOfR(arr[0])
	i1 := indexOfR(append(arr[1], beginIdx))
	i2 := indexOfR(append(arr[len(arr)-1], endIdx))
	children := new(strings.Builder)
	for i, x := range arr[1:] {
		children.WriteString(indexOfR(x))
		if i < len(arr) - 1 {
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

func indexOfR(idxs []int) string {
	buf := new(strings.Builder)
	buf.WriteString("$^R->[")
	for i, x := range idxs {
		buf.WriteString(fmt.Sprintf("%v", x))
		if i != len(idxs) - 1 {
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
	ret := [][]int{}
	if t.left != nil {
		ret = append(ret, t.left.indicesList()...)
	}
	if t.right != nil {
		ret = append(ret, t.right.indicesList()...)
	}
	return ret
}

func unfoldAnnotatedRegexTree(n int) *tree {
	idxs := []int{}
	res := tree{val: idxs}
	ref := &res
	for i := 0; i < n; i++ {
		ref.left = &tree{val: append(ref.val, 0)}
		ref.right = &tree{val: append(ref.val, 1)}
		ref = ref.left
	}
	return &res
}

type tree struct {
	val []int
	left *tree
	right *tree
}

func main() {
	re := ASTRegex { Defines: []Define{{DefineName: "Foo", RegexSteps: []RegexStep{PositionSaveStep{}, MatchStep{"bar"}, MatchSaveStep{"Foo"}}}}}
	fmt.Println(re)
}

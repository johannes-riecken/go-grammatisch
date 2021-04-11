package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/johannes-riecken/go-grammatisch/pkg/go-grammatisch/model"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("pkg/templates/*")
	r.GET("/index", func(c *gin.Context) {
		grammarJSON := c.Query("grammar-content")
		var grammar model.Grammar
		var generatedSource model.ASTRegex
		if grammarJSON == "" {
			c.HTML(http.StatusOK, "index.gohtml", gin.H{})
			return
		}
		err := fakeUnmarshal([]byte(grammarJSON), &grammar)
		if err != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("Unmarshaling grammar JSON failed: %v", err))
			return
		}
		generatedSource = grammar.ToRegex()

		inputDoc := c.Query("input-doc")
		var outputAST []byte
		if inputDoc != "" {
			inputFile, err := os.Create("input.txt")
			if err != nil {
				c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to create input file: %v", err))
				return
			}
			defer inputFile.Close()
			inputFile.WriteString(inputDoc)
			regexFile, err := os.Create("regex.txt")
			if err != nil {
				c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to create regex file: %v", err))
				return
			}
			defer regexFile.Close()
			regexFile.WriteString(generatedSource.String())
			outputAST, err = exec.Command("perl", "generateSyntaxTree.pl").Output()
			if err != nil {
				c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to execute perl script: %v", err))
				return
			}
		}
		c.HTML(http.StatusOK, "index.gohtml", gin.H{"generatedSource": generatedSource, "outputAST": string(outputAST)})
	})
	_ = r.Run()
}

func fakeUnmarshal(b []byte, grammar *model.Grammar) error {
	if string(b) == `Foo : 'bar' ;` {
		*grammar = model.Grammar{
			RuleSpecs: []model.RuleSpec{
				{
					RuleRef: "Foo",
					Alternatives: []model.Alternative{{
						Elements: []model.Element{
							model.Quoted{
								Quoted: "'bar'",
							},
											},
										}},
				},
			},
		}
		return nil
	}
	if string(b) == "foo : Bar Bar;\r\n" +
"Bar : 'baz';" {
		*grammar = model.Grammar{
			RuleSpecs: []model.RuleSpec{{RuleRef: "foo", Alternatives: []model.Alternative{{Elements: []model.Element{model.RuleRef{RuleRefName: "Bar"},
				model.RuleRef{RuleRefName: "Bar"}}}}}, {RuleRef: "Bar", Alternatives: []model.Alternative{{Elements: []model.Element{model.Quoted{Quoted: "'baz'"}}}}}}}
		return nil
	}
	return errors.New("im Moment werden leider nur Grammatiken akzeptiert, die bereits hart im Quelltext kodiert sind")
}

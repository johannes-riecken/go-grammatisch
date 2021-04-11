package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/johannes-riecken/go-grammatisch/pkg/go-grammatisch/model"
	"net/http"
	"os"
	"os/exec"
)

func getIndexHandler(c *gin.Context) {
	grammarJSON := c.Query("grammar-content")
	inputDoc := c.Query("input-doc")
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

	var outputAST []byte
	if inputDoc != "" {
		inputFile, err := os.Create("input.txt")
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to create input file: %v", err))
			return
		}
		defer inputFile.Close()
		_, err = inputFile.WriteString(inputDoc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to write to input file: %v", err))
			return
		}
		regexFile, err := os.Create("regex.txt")
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to create regex file: %v", err))
			return
		}
		defer regexFile.Close()
		_, err = regexFile.WriteString(generatedSource.String())
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to write to regex file: %v", err))
			return
		}
		outputAST, err = exec.Command("perl", "generateSyntaxTree.pl").Output()
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to execute perl script: %v", err))
			return
		}
	}
	c.HTML(http.StatusOK, "index.gohtml", gin.H{"generatedSource": generatedSource, "outputAST": string(outputAST)})
}

package server

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johannes-riecken/go-grammatisch/pkg/go-grammatisch/model"
)

func AddRoutes(r *gin.Engine) {
	globals := make(map[string]string)
	r.LoadHTMLGlob("pkg/templates/*.gohtml")
	r.GET("/step00", func(c *gin.Context) {
		c.HTML(http.StatusOK, "step00.gohtml", globals)
	})
	r.POST("/step01", func(c *gin.Context) {
		addPostFormToGlobals(c, globals)
		grammars := model.Grammars(globals["grammars"])
		grammar, err := grammars.Process()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		grammarJSON, err := json.Marshal(grammar)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		globals["grammar"] = string(grammarJSON)
		c.HTML(http.StatusOK, "step01.gohtml", globals)
	})
	r.POST("/step02", func(c *gin.Context) {
		addPostFormToGlobals(c, globals)
		var grammar model.Grammar
		err := json.Unmarshal([]byte(globals["grammar"]), &grammar)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		astRegex := grammar.Process()
		astRegexJSON, err := json.Marshal(astRegex)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		globals["astRegex"] = string(astRegexJSON)
		c.HTML(http.StatusOK, "step02.gohtml", globals)
	})
	r.POST("/step03", func(c *gin.Context) {
		addPostFormToGlobals(c, globals)
		c.HTML(http.StatusOK, "step03.gohtml", globals)
	})
	r.POST("/step04", func(c *gin.Context) {
		addPostFormToGlobals(c, globals)
		var astRegex model.ASTRegex
		err := json.Unmarshal([]byte(globals["astRegex"]), &astRegex)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		inputDoc := globals["inputDoc"]
		ast, err := astRegex.Process(inputDoc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		globals["ast"] = ast
		c.HTML(http.StatusOK, "step04.gohtml", globals)
	})
	r.POST("/step05", func(c *gin.Context) {
		addPostFormToGlobals(c, globals)
		c.Status(http.StatusOK)
	})
}

func addPostFormToGlobals(c *gin.Context, globals map[string]string) {
	_ = c.PostForm("")
	for k, v := range c.Request.PostForm {
		globals[k] = v[0]
	}
}

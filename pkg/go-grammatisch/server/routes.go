package server

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/johannes-riecken/go-grammatisch/pkg/go-grammatisch/model"
	"log"
	"net/http"
)

func AddRoutes(r *gin.Engine) {
	globals := make(map[string]string)
	r.GET("/index", getIndexHandler)
	//r.Static("/static", "/Users/riecken/repos/go-grammatisch/pkg/templates/step00.gohtml")
	r.LoadHTMLGlob("/Users/rieckenj/repos/go-grammatisch/pkg/templates/*.gohtml")
	r.GET("/step00", func(c *gin.Context) {
		c.HTML(http.StatusOK, "step00.gohtml", globals)
	})
	r.POST("/step01", func(c *gin.Context) {
		log.Println(c.PostForm("grammars"))
		addPostFormToGlobals(c, globals)
		c.HTML(http.StatusOK, "step01.gohtml", globals)
	})
	r.POST("/step02", func(c *gin.Context) {
		log.Println(c.PostForm("grammar"))
		addPostFormToGlobals(c, globals)
		globals["astRegex"] = "" // TODO: Construct Grammar from string and call .ToRegex(). In the beginning, do it all with JSON
		c.HTML(http.StatusOK, "step02.gohtml", globals)
	})
	r.POST("/step03", func(c *gin.Context) {
		log.Println(c.PostForm("astRegex"))
		addPostFormToGlobals(c, globals)
		c.HTML(http.StatusOK, "step03.gohtml", globals)
	})
	r.POST("/step04", func(c *gin.Context) {
		log.Println(c.PostForm("inputDoc"))
		addPostFormToGlobals(c, globals)
		c.HTML(http.StatusOK, "step04.gohtml", globals)
	})
	r.POST("/step05", func(c *gin.Context) {
		log.Println(c.PostForm("ast"))
		addPostFormToGlobals(c, globals)
		c.Status(http.StatusOK)
	})
}

func addPostFormToGlobals(c *gin.Context, globals map[string]string) {
	for k, v := range c.Request.PostForm {
		globals[k] = v[0]
	}
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
	if string(b) == "foo : Bar Bar;\r\n"+
		"Bar : 'baz';" {
		*grammar = model.Grammar{
			RuleSpecs: []model.RuleSpec{{RuleRef: "foo", Alternatives: []model.Alternative{{Elements: []model.Element{model.RuleRef{RuleRefName: "Bar"},
				model.RuleRef{RuleRefName: "Bar"}}}}}, {RuleRef: "Bar", Alternatives: []model.Alternative{{Elements: []model.Element{model.Quoted{Quoted: "'baz'"}}}}}}}
		return nil
	}
	return errors.New("im Moment werden leider nur Grammatiken akzeptiert, die bereits hart im Quelltext kodiert sind")
}

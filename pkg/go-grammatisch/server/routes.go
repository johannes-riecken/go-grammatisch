package server

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/johannes-riecken/go-grammatisch/pkg/go-grammatisch/model"
)

func AddRoutes(r *gin.Engine) {
	r.GET("/index", getIndexHandler)
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

package server

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/johannes-riecken/go-grammatisch/pkg/go-grammatisch/model"
)

func AddRoutes() {
	globals := make(map[string]string)
	tmpl, err := template.ParseGlob("pkg/templates/*.gohtml")
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/step00", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl.ExecuteTemplate(w, "step00", globals)
		}
	})
	http.HandleFunc("/step01", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			addPostFormToGlobals(r, globals)
			grammars := model.Grammars(globals["grammars"])
			grammar, err := grammars.Process()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			grammarJSON, err := json.Marshal(grammar)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			globals["grammar"] = string(grammarJSON)
			tmpl.ExecuteTemplate(w, "step01", globals)
		}
	})
	http.HandleFunc("/step02", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			addPostFormToGlobals(r, globals)
			var grammar model.Grammar
			err := json.Unmarshal([]byte(globals["grammar"]), &grammar)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			astRegex := grammar.Process()
			astRegexJSON, err := json.Marshal(astRegex)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			globals["astRegex"] = string(astRegexJSON)
			tmpl.ExecuteTemplate(w, "step02", globals)
		}
	})
	http.HandleFunc("/step03", func(w http.ResponseWriter, r *http.Request) {
		addPostFormToGlobals(r, globals)
		tmpl.ExecuteTemplate(w, "step03", globals)
	})
	http.HandleFunc("/step04", func(w http.ResponseWriter, r *http.Request) {
		addPostFormToGlobals(r, globals)
		var astRegex model.ASTRegex
		err := json.Unmarshal([]byte(globals["astRegex"]), &astRegex)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		inputDoc := globals["inputDoc"]
		ast, err := astRegex.Process(inputDoc)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		globals["ast"] = ast
		tmpl.ExecuteTemplate(w, "step04", globals)
	})
	http.HandleFunc("/step05", func(w http.ResponseWriter, r *http.Request) {
		addPostFormToGlobals(r, globals)
		w.WriteHeader(http.StatusOK)
	})
}

func addPostFormToGlobals(r *http.Request, globals map[string]string) {
	err := r.ParseForm()
	if err != nil {
		return
	}
	for k, v := range r.PostForm {
		globals[k] = v[0]
	}
}

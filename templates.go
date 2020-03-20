package authentication

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

var templates *template.Template

func initTemplates() {
	log.Println("initialising templates")
	templates = template.New("").Funcs(template.FuncMap{
		"Title": strings.Title,
	})
	templates = template.Must(templates.ParseGlob("./templates/providerSelection.gohtml"))
}

func providerSelectionTemplate(w http.ResponseWriter, r *http.Request, providers []string) {
	handleError(w, r, templates.ExecuteTemplate(w, "providerSelection.gohtml", Providers()))
}

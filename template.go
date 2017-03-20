package gostuff

import (
	"html/template"
	"log"
	"os"
)

// Goes through all templates and parses then on startup
func ParseTemplates() {

	var allNewProviders AllNewsProviders
	allNewProviders.ReadAllNews()
	parseTemplate(allNewProviders, "news.html", []string{"templates/newsTemplate.html",
		"templates/guestHeader.html"}...)

	tempArgs := struct {
		PageTitle string // Title of the web page
	}{
		"Free Online Chess",
	}
	parseTemplate(tempArgs, "index.html", []string{"templates/index.html",
		"templates/guestHeader.html"}...)
}

// @templateArgs Template arguments that will be parsed
// @outputPath the output file of the parsed template
// @templatePath relative location to template that is to be parsed
func parseTemplate(templateArgs interface{}, outputPath string, templatePaths ...string) {

	var t = template.Must(template.ParseFiles(templatePaths...))

	f, err := os.Create(outputPath)
	if err != nil {
		log.Println(err)
		return
	}

	err = t.Execute(f, templateArgs)
	if err != nil {
		log.Println(err)
		return
	}
}

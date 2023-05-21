package gostuff

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

// Goes through all templates and parses then on startup
func OneTimeParseTemplates() {

	tempArgs := struct {
		PageTitle string // Title of the web page
	}{
		"Free Online Chess",
	}
	ParseTemplates(tempArgs, nil, "index.html", []string{"templates/indexTemplate.html",
		"templates/guestHeader.html"}...)

	tempArgs = struct {
		PageTitle string // Title of the web page
	}{
		"Help",
	}
	ParseTemplates(tempArgs, nil, "help.html", []string{"templates/helpTemplate.html",
		"templates/guestHeader.html"}...)

	tempArgs = struct {
		PageTitle string // Title of the web page
	}{
		"Screenshots",
	}
	ParseTemplates(tempArgs, nil, "screenshots.html", []string{"templates/screenshotsTemplate.html",
		"templates/guestHeader.html"}...)
}

func Show404Page(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	var doesNotExist = template.Must(template.ParseFiles("404.html"))

	p := struct {
		Url string
	}{
		r.Host,
	}

	if err := doesNotExist.Execute(w, &p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// @templateArgs Template arguments that will be parsed
// @writer http.ResponseWriter or nil, if nil a file will be created and served
// @outputPath the output file of the parsed template
// @templatePath relative location to template that is to be parsed
func ParseTemplates(templateArgs interface{}, writer http.ResponseWriter, outputPath string,
	templatePaths ...string) {

	log := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	var t = template.Must(template.ParseFiles(templatePaths...))

	if writer == nil {
		f, err := os.Create(outputPath)
		defer f.Close()

		if err != nil {
			log.Println(err)
			return
		}
		err = t.Execute(f, templateArgs)
		if err != nil {
			log.Println(err)
		}
	} else {
		err := t.Execute(writer, templateArgs)
		if err != nil {
			log.Println(err)
		}
	}
}

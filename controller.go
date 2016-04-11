package main

import (
	"html/template"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

var templates *template.Template

func routerInit() *httprouter.Router {
	router := httprouter.New()

	templates = parseTemplates()
	router.GET("/", IndexHandle)
	router.ServeFiles("/css/*filepath", http.Dir("/templates/css/"))

	return router
}

//IndexHandle Стандартная функция обработчик
func IndexHandle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := templates.Lookup("index.html").Execute(w, nil)
	if err != nil {
		http.NotFound(w, r)
	}

}

//func ServeStyles(w http.ResponseWriter, r *http.Request, _ httprouter.Params)

func parseTemplates() *template.Template {
	result := template.New("templates")

	basePath := "templates"
	templateFolder, _ := os.Open(basePath)
	defer templateFolder.Close()

	templatePathsRaw, _ := templateFolder.Readdir(-1)
	// -1 means all of the contents
	templatePaths := new([]string)
	for _, pathInfo := range templatePathsRaw {
		if !pathInfo.IsDir() {
			*templatePaths = append(*templatePaths,
				basePath+"/"+pathInfo.Name())
		}
	}

	result.ParseFiles(*templatePaths...)
	return result
}

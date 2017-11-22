package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"golang.org/x/net/websocket"
	"github.com/kYem/dota-dashboard/ws"
	"github.com/kYem/dota-dashboard/controller"
)

func main() {
	templates := populateTemplates()

	http.HandleFunc("/",
		func(w http.ResponseWriter, req *http.Request) {
			controller.SetDefaultHeaders(w)
			requestedFile := req.URL.Path[1:]
			tmpl := templates.Lookup(requestedFile + ".html")

			var context interface{} = nil
			if tmpl != nil {
				tmpl.Execute(w, context)
			} else {
				w.WriteHeader(404)
			}
		})

	http.HandleFunc("/img/", serveResource)
	http.HandleFunc("/css/", serveResource)
	http.HandleFunc("/live/stats", controller.LiveGamesStats)
	http.HandleFunc("/live", controller.LiveGames)
	http.HandleFunc("/matches", controller.HomePage)
	http.Handle("/socket", websocket.Handler(ws.Echo))
	log.Fatal(http.ListenAndServe(":8008", nil))
}

func serveResource(w http.ResponseWriter, req *http.Request) {
	path := "public" + req.URL.Path
	var contentType string

	if strings.HasSuffix(path, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(path, ".png") {
		contentType = "image/png"
	} else {
		contentType = "text/html"
	}
	f, err := os.Open(path)

	if err == nil {
		defer f.Close()
		w.Header().Add("Content-Type", contentType)

		br := bufio.NewReader(f)
		br.WriteTo(w)
	} else {
		w.WriteHeader(404)
	}
}

func populateTemplates() *template.Template {

	result := template.New("templates")

	basePath := "templates"
	templateFolder, _ := os.Open(basePath)
	defer templateFolder.Close()

	templatePathsRaw, _ := templateFolder.Readdir(-1)

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

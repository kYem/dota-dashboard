package main

import (
	"io"
	"log"
	"net/http"
	"github.com/kYem/dota-dashboard/api"
	"github.com/kYem/dota-dashboard/config"
	"io/ioutil"
	"os"
	"text/template"
	"strings"
	"bufio"
)


func HomePage(w http.ResponseWriter, req *http.Request) {

	client := api.GetClient(config.LoadConfig())

	resp := client.GetTopLiveGames("1")

	if resp.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	io.WriteString(w, string(body))
}

func LiveGames(w http.ResponseWriter, req *http.Request) {

	setDefaultHeaders(w)
	partner := req.URL.Query().Get("partner")
	if partner == "" {
		partner = "0"
	}
	client := api.GetClient(config.LoadConfig())
	resp := client.GetTopLiveGames(partner)

	if resp.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	io.WriteString(w, string(body))
}

func LiveGamesStats(w http.ResponseWriter, req *http.Request) {
	setDefaultHeaders(w)
	serverSteamId := req.URL.Query().Get("server_steam_id")
	client := api.GetClient(config.LoadConfig())
	resp := client.GetRealTimeStats(serverSteamId)

	if resp.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	io.WriteString(w, string(body))
}

func main() {
	templates := populateTemplates()

	http.HandleFunc("/",
		func(w http.ResponseWriter, req *http.Request) {
			setDefaultHeaders(w)
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
	http.HandleFunc("/live/stats", LiveGamesStats)
	http.HandleFunc("/live", LiveGames)
	http.HandleFunc("/matches", HomePage)

	log.Fatal(http.ListenAndServe(":8888", nil))
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

func setDefaultHeaders(w http.ResponseWriter) {
	w.Header().Add("access-control-allow-credentials", "true")
	w.Header().Add("access-control-allow-origin", "http://dotatv.com:3000")
	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
}



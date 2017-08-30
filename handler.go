package gographiql

import (
	"log"
	"text/template"
	"github.com/GeertJohan/go.rice"
	"net/http"
	"strings"
	"fmt"
)

type Config struct {
	BasePath   string
	GraphqlUrl string
	OauthScope string
}
type Handler struct {
	config        *Config
	staticHandler http.Handler
	template      *template.Template
	static        *rice.Box
}

func New(config *Config) http.Handler {
	conf := rice.Config{
		LocateOrder: []rice.LocateMethod{rice.LocateEmbedded, rice.LocateAppended, rice.LocateFS},
	}
	static, err := conf.FindBox("graphiql")
	if err != nil {
		log.Fatalf("error opening static: %s\n", err)
	}
	viewBox, err := rice.FindBox("graphiql-template")
	if err != nil {
		log.Fatal(err)
	}

	indexContent, err := viewBox.String("index.html")
	if err != nil {
		log.Fatal(err)
	}
	// parse and execute the template
	indexTemplate, err := template.New("message").Parse(indexContent)
	if err != nil {
		log.Fatal(err)
	}

	return &Handler{
		config:        config,
		static:        static,
		staticHandler: http.FileServer(static.HTTPBox()),
		template:      indexTemplate,
	}

}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	urlPath := r.URL.Path
	prefix := strings.TrimPrefix(r.URL.Path, s.config.BasePath)
	fmt.Println(urlPath, prefix)
	if file, err := s.static.Open(prefix); err == nil {
		if f, e := file.Stat(); e == nil && !f.IsDir() {
			r.URL.Path = prefix
			s.staticHandler.ServeHTTP(w, r)
			return
		}
	}
	error := s.template.Execute(w, s.config)
	if error != nil {
		log.Fatal(error)
	}
}

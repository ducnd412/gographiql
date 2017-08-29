package gographiql

import (
	"log"
	"text/template"
	"github.com/GeertJohan/go.rice"
	"net/http"
	"strings"
)

type Config struct {
	GraphqlUrl string
	OauthScope string
}
type Handler struct {
	config        *Config
	staticHandler http.Handler
	template      *template.Template
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
		staticHandler: http.FileServer(static.HTTPBox()),
		template:      indexTemplate,
	}

}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	upath := r.URL.Path
	basePath := upath[0:strings.LastIndex(upath, "/")]
	if strings.HasSuffix(upath, ".png") || strings.HasSuffix(upath, ".css") || strings.HasSuffix(upath, ".js") {
		if p := strings.TrimPrefix(r.URL.Path, basePath); len(p) < len(r.URL.Path) {
			r.URL.Path = p
			s.staticHandler.ServeHTTP(w, r)
		}
		return
	}
	data := map[string]interface{}{
		"BasePath": basePath,
		"Config":   s.config,
	}
	error := s.template.Execute(w, data)
	if error != nil {
		log.Fatal(error)
	}
}

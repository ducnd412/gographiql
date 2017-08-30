package gographiql

import (
	"testing"
	"fmt"
	"net/http"
	"log"
)

func Test(t *testing.T) {
	swaggerHandler := New(&Config{

		BasePath:   "/graphql",
		GraphqlUrl: "http://localhost:8083/graphql",
		OauthScope: "scope",
	})

	http.Handle("/graphql/", swaggerHandler)
	go func() {
		fmt.Println("Serving files on :8082, press ctrl-C to exit")
		err := http.ListenAndServe(":8082", nil)
		if err != nil {
			log.Fatalf("error serving files: %v", err)
		}
	}()
	select {}
}

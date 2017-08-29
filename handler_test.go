package swaggerui

import (
	"testing"
	"fmt"
	"net/http"
	"log"
)

func Test(t *testing.T) {
	swaggerHandler := New(&Config{

		GraphqlUrl: "http://localhost:8083/graphql",
		OauthScope: "scope",
	})

	http.Handle("/swagger/", swaggerHandler)
	go func() {
		fmt.Println("Serving files on :8082, press ctrl-C to exit")
		err := http.ListenAndServe(":8082", nil)
		if err != nil {
			log.Fatalf("error serving files: %v", err)
		}
	}()
	select {}
}

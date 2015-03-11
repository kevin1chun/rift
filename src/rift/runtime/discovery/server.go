package discovery

import (
	"io"
	"log"
	"net/http"
	"rift/runtime"
	"rift/support/collections"
)

func Start() {
	ctx := runtime.Context{collections.NewPersistentMap()}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "OPTION":
			io.WriteString(w, "OPTIONS")
		case "GET":
			
		}	
	})

	log.Fatal(http.ListenAndServe(":8830", nil))
}
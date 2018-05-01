package analyzer

import (
	"fmt"
	// "encoding/json"
    "log"
	"net/http"
)

func NewServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("test")
	})
	
	log.Fatal(http.ListenAndServe(":3001", nil))
}

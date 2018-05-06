package analyzer

import (
	"encoding/json"
	"fmt"
	"mime"
	"mime/multipart"
    "log"
	"net/http"
	"os"
	"strings"
)

func Serve() {
	
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		mediaType, params, mimeErr := mime.ParseMediaType(req.Header.Get("Content-Type"))

		if mimeErr != nil {
			

		} else if !strings.HasPrefix(mediaType, "multipart/") {
			

		} else {
			partReader := multipart.NewReader(req.Body, params["boundary"])
			manager := ImagesFromParts(partReader)
			manager.ProcessItems()

			bytes, err := json.Marshal(manager)

			if err == nil {
				fmt.Fprint(w, string(bytes[:]))
			} else {
				fmt.Fprint(w, err)
			}
		}
	})
	
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil))
}

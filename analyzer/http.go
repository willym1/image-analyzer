package analyzer

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"mime"
	"mime/multipart"
    "log"
	"net/http"
	"strings"
)

func NewServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		mediaType, params, mimeErr := mime.ParseMediaType(req.Header.Get("Content-Type"))

		if mimeErr != nil {
			

		} else if !strings.HasPrefix(mediaType, "multipart/") {
			

		} else {
			reader := multipart.NewReader(req.Body, params["boundary"])
			part, _ := reader.NextPart()
	
			b, _ := ioutil.ReadAll(part)
			bytesReader := bytes.NewReader(b)
			_, format, err := image.Decode(bytesReader)

			if err != nil {
				fmt.Println("Error decoding image")

			} else {
				fmt.Println(format)
			}

			fmt.Fprintln(w, "test")
		}
	})
	
	log.Fatal(http.ListenAndServe(":3001", nil))
}

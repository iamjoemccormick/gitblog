package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gomarkdown/markdown"
)

var (
	staticBase         string = "static/"
	supportedFileTypes        = [3]string{".md", ".html", ".htm"}
)

func main() {
	log.SetOutput(os.Stdout)
	http.HandleFunc("/", urlHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func urlHandler(w http.ResponseWriter, r *http.Request) {

	requestedPath := r.URL.Path[1:]
	ft, content, err := readFile(requestedPath)

	if err != nil {
		//TODO: Better handling of 404 and internal server errors include returning proper HTTP response codes.

		if os.IsNotExist(err) {
			fmt.Fprintf(w, "404: Page '%s' does not exist!", requestedPath)
		} else {
			fmt.Fprintf(w, "500: Internal server error attempting to access '%s'.", requestedPath)
		}
	} else {

		fmt.Fprintf(w, "<head></head>")
		if ft == ".md" {
			content = markdown.ToHTML(content, nil, nil)
		}
		w.Write(content)
		//TODO: Sanitize untrusted content: https://github.com/gomarkdown/markdown#sanitize-untrusted-content
	}

}
func readFile(path string) (fileType string, content []byte, err error) {

	for _, t := range supportedFileTypes {
		content, err = ioutil.ReadFile(staticBase + path + t)

		if err != nil && os.IsNotExist(err) {
			log.Printf("Did not find a file named '%s'", staticBase+path+t)
			continue
		} else if err != nil {
			log.Printf("Unknown error attempting to read file '%s': %s", path, err)
			return "", nil, err
		} else {
			return t, content, nil
		}

	}

	log.Printf("No file found at '%s' with the following file types %s", path, supportedFileTypes)
	return "", nil, err
}

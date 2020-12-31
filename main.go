package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
)

var (
	staticBase            string = "static/"
	supportedFileTypes           = [3]string{".md", ".html", ".htm"}                                   // In order of greatest to least priority in case of duplicates.
	supportedDefaultFiles        = [2]string{"index", "home"}                                          // In order of greatest to least priority in case of duplicates.
	baseTemplatePath      string = "base.html"                                                         // If a path is omitted the base template will be "{{ content }}".
	baseTemplate                 = template{path: baseTemplatePath, template: []byte("{{ content }}")} // Global variable used to store the contents of a base template.
)

func main() {
	log.SetOutput(os.Stdout)
	baseTemplate.loadTemplate()
	http.HandleFunc("/", handleURL)
	http.HandleFunc("/gitposthook", handleGitPostHook)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// WIP: gitPostHook will be used to pull down changes when an HTTP POST is received from the configured repository.
// For now we just ensure the base template is reloaded.
func handleGitPostHook(w http.ResponseWriter, r *http.Request) {

	baseTemplate.loadTemplate()
	fmt.Fprintf(w, "Successfully reloaded base template.")
}

func handleURL(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path[1:]
	var (
		ft      string
		content []byte
		err     error
	)

	if path == "" || strings.HasSuffix(path, "/") {
		ft, content, err = handleDirectory(path)
	} else {
		ft, content, err = readFile(path)
	}

	if err != nil {
		//TODO: Better handling of 404 and internal server errors include returning proper HTTP response codes.

		if os.IsNotExist(err) {
			fmt.Fprintf(w, "404: Page '%s' does not exist!", path)
		} else {
			fmt.Fprintf(w, "500: Internal server error attempting to access '%s'.", path)
		}
	} else {
		if ft == ".md" {
			//TODO: Sanitize untrusted content: https://github.com/gomarkdown/markdown#sanitize-untrusted-content
			content = (markdown.ToHTML(content, nil, nil))
		}
		w.Write(baseTemplate.parseTemplate(content))
	}
}

// handleDirectory checks if there is a supportedDefaultFile in a given directory path and returns the result of readFile if it exists, or an error.
func handleDirectory(path string) (fileType string, content []byte, err error) {

	for _, d := range supportedDefaultFiles {
		ft, content, err := readFile(path + d)

		if err != nil {
			continue
		}
		return ft, content, err
	}

	//TODO: Consider automatically generating a default page if one was not specified.

	log.Printf("No file found under '%s' with the following names %s and file types %s", path, supportedDefaultFiles, supportedFileTypes)
	return "", nil, err
}

// readFile takes a path that may or may not end in a filename extension.
// If the path does not end in a supportedFileType it will check if there is a file with any of the supportedFileTypes at that path.
// If the path ends in a supportedFileType but no file exists for that type, an error will be returned.
func readFile(path string) (fileType string, content []byte, err error) {

	for _, t := range supportedFileTypes {

		if strings.HasSuffix(path, t) {
			content, err = ioutil.ReadFile(staticBase + path)
		} else {
			content, err = ioutil.ReadFile(staticBase + path + t)
		}

		if err != nil && os.IsNotExist(err) {
			log.Printf("Did not find a file at '%s' with supported file type '%s'", staticBase+path, t)
			continue
		} else if err != nil {
			log.Printf("Unknown error attempting to read file '%s' with ending '%s': %s", staticBase+path, t, err)
			return "", nil, err
		} else {
			return t, content, nil
		}

	}

	log.Printf("No file found at '%s' with the following file types %s", path, supportedFileTypes)
	return "", nil, err
}

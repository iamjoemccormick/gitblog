package main

import (
	"log"
	"strings"
)

type template struct {
	path     string
	template []byte
}

// loadTemplate takes a template and if a path is provided the contents are loaded from disk.
// If the path is empty then contents are set to "".
// Returns a fatal error if a path was provided but an error occured reading the file.
func (t *template) loadTemplate() (err error) {

	if t.path != "" {
		_, t.template, err = readFile(t.path)

		if err != nil {
			log.Fatalf("FATAL: A path to a base template was provided, but the following error occured trying to read the file: %s", err)
			return err
		}
		log.Printf("Successfully loaded base template from path '%s'", t.path)
		return err
	}
	t.template = []byte("")
	log.Printf("No path to a base template was provided, defaulting to '%s'", t.template)
	return err
}

func (t *template) parseTemplate(content []byte) (parsedTemplate []byte) {

	//w.Write([]byte(strings.Replace(string(baseTemplate.template), "{{ content }}", string(content), 1)))

	return []byte(strings.Replace(string(baseTemplate.template), "{{ content }}", string(content), 1))
}

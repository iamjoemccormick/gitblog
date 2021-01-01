package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type template struct {
	path     string
	template []byte
}

type navItem struct {
	text     string
	path     string
	children []navItem
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

func buildNav() (nav string) {

	var directories []string

	filepath.Walk(staticBase, func(path string, info os.FileInfo, err error) error {
		if strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		if info.IsDir() {
			directories = append(directories, path)
		}
		return nil
	})

	// TODO: Parse directories into a nav structure.
	return nav
}

func (t *template) parseTemplate(content []byte) (parsedTemplate []byte) {

	// TODO: Is there a more efficient way to do this?
	return []byte(strings.Replace(string(baseTemplate.template), "{{ content }}", string(content), 1))
}

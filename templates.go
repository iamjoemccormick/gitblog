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
	url      string
	children []*navItem
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

// buildNav() dynamically generates the site navigation based on the directory structure under staticBase.
func buildNav() (nav string) {

	var directories []string

	filepath.Walk(staticBase, func(path string, info os.FileInfo, err error) error {
		if strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		if info.IsDir() {
			directories = append(directories, strings.Join(strings.Split(path, staticBase)[1:], ""))
		}
		return nil
	})

	baseNav := []*navItem{}

	for _, dir := range directories {
		currentNav := &baseNav

		for _, section := range strings.Split(dir, "/") {
			appendSection := true

			for _, c := range *currentNav {
				if c.text == section {
					appendSection = false
					currentNav = &c.children
				}
			}

			if appendSection == true {
				newNav := &navItem{text: section, url: dir}
				*currentNav = append(*currentNav, newNav)

				// Note this is technically unnecessary since filepath.Walk will return in lexical order, in other words we'll always get tech/ before tech/golang.
				// However this ensures this function still works even if we get a child directory before we have an entry for it's parent in currentNav.
				currentNav = &newNav.children
			}
		}
	}

	//TODO: Figure out what we want to return from this function.
	// OPTION 1: Parse out nav structure (regardless of depth) into HTML, probably has an unordered list that can just be inserted into a template.
	// OPTION 2: Return a structure that can be iterated over from a template allowing highly customizable navigation.
	// This would require implementing templating support for for loops then puts the burden on the user to worry about how deep the navigation goes
	// (or how deep a navigation bar they want to display).

	return nav
}

func (t *template) parseTemplate(content []byte) (parsedTemplate []byte) {

	// TODO: Is there a more efficient way to do this?
	return []byte(strings.Replace(string(baseTemplate.template), "{{ content }}", string(content), 1))
}

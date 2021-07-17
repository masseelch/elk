package internal

import (
	"fmt"
	"golang.org/x/tools/imports"
	"io/ioutil"
	"os"
)

type (
	file struct {
		path    string
		content []byte
	}
	assets struct {
		dirs  []string
		files []file
	}
)

// write files and dirs in the assets.
func (a assets) write() error {
	for _, dir := range a.dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("create dir %q: %w", dir, err)
		}
	}
	for _, file := range a.files {
		if err := ioutil.WriteFile(file.path, file.content, 0644); err != nil {
			return fmt.Errorf("write file %q: %w", file.path, err)
		}
	}
	return nil
}

// format runs "goimports" on all assets.
func (a assets) format() error {
	for _, file := range a.files {
		path := file.path
		src, err := imports.Process(path, file.content, nil)
		if err != nil {
			return fmt.Errorf("format file %s: %v", path, err)
		}
		if err := ioutil.WriteFile(path, src, 0644); err != nil {
			return fmt.Errorf("write file %s: %v", path, err)
		}
	}
	return nil
}

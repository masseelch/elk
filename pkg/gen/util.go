package gen

import (
	"fmt"
	"github.com/facebook/ent/schema/field"
	"golang.org/x/tools/imports"
	"io/ioutil"
	"os"
	"os/exec"
)

var (
	dartTypeNames = [...]string{
		field.TypeInvalid: "dynamic",
		field.TypeBool:    "bool",
		field.TypeTime:    "DateTime",
		field.TypeJSON:    "Map<String, dynamic>",
		field.TypeUUID:    "String",
		field.TypeBytes:   "dynamic",
		field.TypeEnum:    "String",
		field.TypeString:  "String",
		field.TypeInt:     "int",
		field.TypeInt8:    "int",
		field.TypeInt16:   "int",
		field.TypeInt32:   "int",
		field.TypeInt64:   "int",
		field.TypeUint:    "int",
		field.TypeUint8:   "int",
		field.TypeUint16:  "int",
		field.TypeUint32:  "int",
		field.TypeUint64:  "int",
		field.TypeFloat32: "double",
		field.TypeFloat64: "double",
	}
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

// formatGo runs "goimports" on all assets.
func (a assets) formatGo() error {
	for _, file := range a.files {
		path := file.path
		src, err := imports.Process(path, file.content, nil)
		if err != nil {
			return fmt.Errorf("formatGo file %s: %v", path, err)
		}
		if err := ioutil.WriteFile(path, src, 0644); err != nil {
			return fmt.Errorf("write file %s: %v", path, err)
		}
	}
	return nil
}

// formatDart runs "dartfmt" on all assets.
func (a assets) formatDart() error {
	args := []string{"-w"}
	for _, dir := range a.dirs {
		args = append(args, dir)
	}

	if err := exec.Command("dartfmt", args...).Run(); err != nil {
		if err == exec.ErrNotFound {
			// "dartfmt" is not available
			fmt.Println("The command 'dartfmt' was not found on your system. Generated code remains unformatted.")
		} else {
			return err
		}
	}

	return nil
}

// Dart type for a given go type.
func dartType(t *field.TypeInfo) string {
	return dartTypeNames[t.Type]
}

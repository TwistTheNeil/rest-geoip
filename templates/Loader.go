package templates

import (
	"html/template"
	"io/ioutil"
	"strings"
)

// LoadTemplate loads templates embedded by go-assets-builder
func LoadTemplate() (*template.Template, error) {
	t := template.New("")
	for name, file := range Assets.Files {
		defer file.Close()
		if file.IsDir() || !strings.HasSuffix(name, ".tmpl") {
			continue
		}
		h, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		t, err = t.New(name).Parse(string(h))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

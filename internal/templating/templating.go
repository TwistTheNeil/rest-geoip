package templating

import (
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	"github.com/markbates/pkger"
)

func ParseTemplates(root string) (*template.Template, error) {
	t := template.New("")

	// Since Walk receives a dynamic value, pkger won't be able to find the
	// actual directory to package from the next line, which is why we used
	// pkger.Include() in main().
	err := pkger.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		f, err := pkger.Open(path)
		if err != nil {
			return err
		}
		// We read from pkger's fs here so the template can be parsed
		contents, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		_, err = t.Parse(string(contents))
		if err != nil {
			return err
		}
		return nil
	})

	return t, err
}

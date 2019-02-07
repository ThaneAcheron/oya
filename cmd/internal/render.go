package internal

import (
	"io"

	"github.com/bilus/oya/pkg/project"
	"github.com/bilus/oya/pkg/template"
	"github.com/pkg/errors"
)

func Render(oyafilePath, templatePath, outputPath, alias string, stdout, stderr io.Writer) error {
	proj, err := project.Detect(outputPath)
	if err != nil {
		return err
	}

	o, found, err := proj.Oyafile(oyafilePath)
	if err != nil {
		return err
	}

	var values template.Scope
	if found {
		if alias != "" {
			av, found := o.Values[alias]
			if !found {
				return errors.Errorf("Alias %s not found in project", alias)
			}
			aliasScope, found := av.(template.Scope)
			if !found {
				return errors.Errorf("Alias %s not found in project", alias)
			}
			values = aliasScope
		} else {
			values = o.Values
		}
	}

	return template.RenderAll(templatePath, outputPath, values)
}

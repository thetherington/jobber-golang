package mailtemplates

import (
	"bytes"
	"embed"
	"html/template"
)

var (
	//go:embed templates
	templateFolder embed.FS
)

type TemplateMaker struct {
	File     string
	Template *template.Template
}

func NewTemplateMaker(file string) (*TemplateMaker, error) {
	data, err := templateFolder.ReadFile("templates/" + file)
	if err != nil {
		return nil, err
	}

	template := template.Must(template.New(file).Parse(string(data)))

	return &TemplateMaker{
		file,
		template,
	}, nil
}

func (t *TemplateMaker) Render(locals interface{}) (string, error) {
	var body bytes.Buffer

	err := t.Template.ExecuteTemplate(&body, t.File, locals)
	if err != nil {
		return "", err
	}

	return body.String(), nil
}

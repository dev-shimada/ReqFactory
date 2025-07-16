package request

import (
	"bytes"
	"net/http"
	"strings"
	"text/template"
)

type Factory struct {
	method         string
	urlTemplate    *template.Template
	headerTemplate *template.Template
	bodyTemplate   *template.Template
}

func NewFactory(method, urlTemplate, headerTemplate, bodyTemplate string) *Factory {
	return &Factory{
		method:         method,
		urlTemplate:    template.Must(template.New("url").Parse(urlTemplate)),
		headerTemplate: template.Must(template.New("header").Parse(headerTemplate)),
		bodyTemplate:   template.Must(template.New("body").Parse(bodyTemplate)),
	}
}

func (f *Factory) Build(header, row []string) (*http.Request, error) {
	data := make(map[string]string)
	for i, h := range header {
		data[h] = row[i]
	}

	var url bytes.Buffer
	if err := f.urlTemplate.Execute(&url, data); err != nil {
		return nil, err
	}

	var body bytes.Buffer
	if err := f.bodyTemplate.Execute(&body, data); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(f.method, url.String(), &body)
	if err != nil {
		return nil, err
	}

	var headers bytes.Buffer
	if err := f.headerTemplate.Execute(&headers, data); err != nil {
		return nil, err
	}

	for _, line := range strings.Split(headers.String(), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
	}

	return req, nil
}
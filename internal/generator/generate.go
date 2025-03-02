package main

import (
	"bytes"
	"encoding/json"
	"go/format"
	"html/template"
	"os"
)

const (
	registryTemplate = `// This file is auto generated, DO NOT EDIT.
package libdnsregistry

import (
	{{- range .}}
	libdns{{.Name}} "{{.Package}}"
	{{- end}}
)

var registry = RegistryStore{
	{{- range .}}
	"{{.Name}}": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdns{{.Name}}.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL: "{{ .URL }}",
				Description: "{{ .Description }}",
				Configuration: configurationDetails[libdns{{.Name}}.Provider](),
			}
		},
	},
	{{- end}}
}
`
)

type providerSource struct {
	Name        string `json:"name"`
	Package     string `json:"package"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

func main() {
	args := os.Args[1:]
	if len(args) != 3 {
		panic("Usage: go run generate code <json path> <output path>")
	}

	providers := new([]providerSource)

	listJSON, err := os.ReadFile(args[1])
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(listJSON, providers)
	if err != nil {
		panic(err)
	}

	switch args[0] {
	case "code":
		generateCode(args[2], *providers)

	default:
		panic("type not supported")
	}
}

func generateCode(path string, providers []providerSource) {
	registryTpl, err := template.New("libdns").Parse(registryTemplate)
	if err != nil {
		panic(err)
	}

	content := &bytes.Buffer{}
	if err := registryTpl.Execute(content, providers); err != nil {
		panic(err)
	}

	formattedContent, err := format.Source(content.Bytes())
	if err != nil {
		panic(err)
	}

	//nolint:gosec,mnd
	if err := os.WriteFile(path, formattedContent, 0o644); err != nil {
		panic(err)
	}
}

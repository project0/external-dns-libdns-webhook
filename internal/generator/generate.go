package main

import (
	"bytes"
	"encoding/json"
	"go/format"
	"html/template"
	"os"
)

const (
	registryTemplate = `package libdns

// This file is auto generated, do not modify directly
import (
	{{- range .}}
	libdns{{.Name}} "{{.Package}}"
	{{- end}}
)

var Registry = map[string]func(conf [][]byte) (LibdnsProvider, error) {
	{{- range .}}
	"{{.Name}}": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdns{{.Name}}.Provider](conf)
	},
	{{- end}}
}
`
)

type provderSource struct {
	Name        string `json:"name"`
	Package     string `json:"package"`
	Description string `json:"description,omitempty"`
}

func main() {
	args := os.Args[1:]
	if len(args) != 3 {
		panic("Usage: go run generate code <json path> <output path>")
	}

	providers := new([]provderSource)
	listJson, err := os.ReadFile(args[1])
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(listJson, providers)
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

func generateCode(path string, providers []provderSource) {

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
	if err := os.WriteFile(path, formattedContent, 0644); err != nil {
		panic(err)
	}
}

//go:build ignore

// DISCLAIMER: There is **absolutely** no need for creating a file/flow. Why I've done it?
// - My laziness: I want to automate everything as much as possible (but the benefit here is very dubious)
// - Learning experience: I'm always curious about auto-generated code, and this approach is the easiest. This use
//                        case does not justify using go/ast or alike.
//
// I follow Rob's Pike advice: Please use go generate creatively. It’s there to encourage experimentation.

package main

import (
	"ddia/src/server"
	"encoding/json"
	"log"
	"os"
	"strings"
	"text/template"
	"time"
)

func main() {
	err := generate()
	if err != nil {
		log.Fatal(err)
	}
}

func generate() error {
	if err := generateJSON(); err != nil {
		return err
	}

	if err := generateConst(); err != nil {
		return err
	}

	return nil
}

func generateConst() error {
	commands := make([]string, 0, len(server.Commands))
	for _, c := range server.Commands {
		if c.Status == "won't-do" {
			continue
		}
		commands = append(commands, c.Name)
	}

	f, err := os.Create("commands_const.go")
	if err != nil {
		return err
	}
	defer f.Close()

	return packageTemplate.
		Execute(f, struct {
			Timestamp time.Time
			Commands  []string
		}{
			Timestamp: time.Now(),
			Commands:  commands,
		})
}

func generateJSON() error {
	fd, err := os.Create("commands.json")
	if err != nil {
		return err
	}
	defer fd.Close()

	raw, err := json.MarshalIndent(server.Commands, "", "    ")
	if err != nil {
		return err
	}

	_, err = fd.Write(raw)
	return err
}

var templ = `// Code generated by go generate; DO NOT EDIT.
// To recreate run: make generate
// {{ .Timestamp }}
package server


const (
{{- range .Commands }}
	{{ printf "// %s command" . }}
	{{ printf "%s = " . }}{{ printf "\"%s\"" . | ToUpper }}

{{- end }}
)
`

var packageTemplate = template.Must(
	template.
		New("").
		Funcs(template.FuncMap{
			"ToUpper": strings.ToUpper,
		}).
		Parse(templ),
)

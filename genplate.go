package main

import (
	"log"
	"os"
	"strings"
	"text/template"
)

func main() {
	anyIsHelp := false
	for _, arg := range os.Args {
		if arg == "-h" || arg == "--help" {
			anyIsHelp = true
			break
		}
	}

	if anyIsHelp || len(os.Args) < 3 {
		usage()
	}

	templateFile := os.Args[1]
	tmpl := template.Must(template.ParseFiles(templateFile))

	outFile := os.Args[2]
	out, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	templateArgs := os.Args[3:]
	argMap := map[string]string{}
	for _, arg := range templateArgs {
		split := strings.SplitN(arg, "=", 2)
		argMap[split[0]] = split[1]
	}

	err = tmpl.Execute(out, argMap)
	if err != nil {
		log.Fatal(err)
	}
}

func usage() {
	println(
		`genplate -- code generation with the template package
USAGE:

	genplate template_file out_file [argname=argval...]`,
	)

	os.Exit(1)
}

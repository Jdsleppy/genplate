package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"
	"unicode"
)

func main() {
	anyIsHelp := false
	for _, arg := range os.Args {
		if arg == "-h" || arg == "--help" {
			anyIsHelp = true
			break
		}
	}

	if anyIsHelp || len(os.Args) != 4 {
		usage()
	}

	// template file
	templateFileArg := os.Args[1]
	tmpl := template.New(templateFileArg)
	tmpl = tmpl.Funcs(usefulFuncs())
	tmpl = template.Must(tmpl.ParseFiles(templateFileArg))

	// output file
	outFileArg := os.Args[2]
	outFile, err := os.Create(outFileArg)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	// data file
	dataFileArg := os.Args[3]
	dataFileContents, err := ioutil.ReadFile(dataFileArg)
	if err != nil {
		log.Fatal(err)
	}
	var templateData interface{}
	if err := json.Unmarshal(dataFileContents, &templateData); err != nil {
		log.Fatal(err)
	}

	// render
	if err := tmpl.Execute(outFile, templateData); err != nil {
		log.Fatal(err)
	}
}

func usage() {
	println(
		`genplate -- code generation with the template package
USAGE:

	genplate template_file out_file data_file
	
template_file    relative path to a go-style template file
out_file         relative path to the output, which will be truncated
data_file        relative path to a JSON file to be passed in as template data`,
	)

	os.Exit(1)
}

func usefulFuncs() template.FuncMap {
	return template.FuncMap{
		"Pluralize":  pluralize,
		"CamelCase":  camelCase,
		"PascalCase": pascalCase,
		"SnakeCase":  snakeCase,
	}
}

func pluralize(in string) (string, error) {
	runes := []rune(in)
	var out []rune

	// copy all but last rune
	for i := 0; i < len(runes)-1; i++ {
		out = append(out, runes[i])
	}

	lastRune := runes[len(runes)-1]

	switch lastRune {
	case 's':
		out = append(out, 's', 'e', 's')
	case 'y':
		out = append(out, 'i', 'e', 's')
	default:
		out = append(out, lastRune, 's')
	}

	return string(out), nil
}

func camelCase(in string) (string, error) {
	runes := []rune(in)
	var out []rune

	switch {
	case isCamelCase(in):
		return in, nil
	case isPascalCase(in):
		for i, r := range runes {
			if i == 0 {
				out = append(out, unicode.ToLower(r))
			} else {
				out = append(out, r)
			}
		}
	case isSnakeCase(in):
		for i, r := range runes {
			if r == '_' {
				continue
			}
			if i > 0 && runes[i-i] == '_' {
				out = append(out, unicode.ToUpper(r))
			} else {
				out = append(out, r)
			}
		}
	default:
		return "", fmt.Errorf("cannot convert %s to camelCase", in)
	}

	return string(out), nil
}

func pascalCase(in string) (string, error) {
	runes := []rune(in)
	var out []rune

	switch {
	case isCamelCase(in):
		for i, r := range runes {
			if i == 0 {
				out = append(out, unicode.ToUpper(r))
			} else {
				out = append(out, r)
			}
		}
	case isPascalCase(in):
		return in, nil
	case isSnakeCase(in):
		for i, r := range runes {
			if r == '_' {
				continue
			}
			if i == 0 || runes[i-i] == '_' {
				out = append(out, unicode.ToUpper(r))
			} else {
				out = append(out, r)
			}
		}
	default:
		return "", fmt.Errorf("cannot convert %s to pascalCase", in)
	}

	return string(out), nil
}

func snakeCase(in string) (string, error) {
	runes := []rune(in)
	var out []rune

	switch {
	case isCamelCase(in):
		for _, r := range runes {
			if unicode.IsUpper(r) {
				out = append(out, '_', unicode.ToLower(r))
			} else {
				out = append(out, r)
			}
		}
	case isPascalCase(in):
		for i, r := range runes {
			if i == 0 {
				out = append(out, unicode.ToLower(r))
			} else if unicode.IsUpper(r) {
				out = append(out, '_', unicode.ToLower(r))
			} else {
				out = append(out, r)
			}
		}
	case isSnakeCase(in):
		return in, nil
	default:
		return "", fmt.Errorf("cannot convert %s to snakeCase", in)
	}

	return string(out), nil
}

func isCamelCase(in string) bool {
	runes := []rune(in)

	for i, r := range runes {
		if i == 0 {
			if unicode.IsUpper(r) {
				return false
			}
		}
		if !unicode.IsLetter(r) {
			return false
		}
	}

	return true
}

func isPascalCase(in string) bool {
	runes := []rune(in)

	for i, r := range runes {
		if i == 0 {
			if unicode.IsLower(r) {
				return false
			}
		}
		if !unicode.IsLetter(r) {
			return false
		}
	}

	return true
}

func isSnakeCase(in string) bool {
	runes := []rune(in)

	for _, r := range runes {
		if !unicode.IsLower(r) && r != '_' {
			return false
		}
	}

	return true
}

# genplate
~Code~ Anything generation with golang templates.

```
genplate -- file generation with golang templates
USAGE:

        genplate template_file out_file data_file

template_file    relative path to a golang template file
out_file         relative path to the output, which will be truncated
data_file        relative path to a JSON file to be passed in as template data
```

## Examples

Full examples [here](https://github.com/Jdsleppy/genplate/tree/master/examples)

### `go generate` example

`genplate/examples/generate.go`
```
package examples

//go:generate genplate example.template generatedcowboy.go cowboylingo.json
//go:generate genplate example.template generatedrobot.go robotlingo.json
```

`genplate/examples/cowboylingo.json`
```json
{
    "FuncName": "CowboyLingo",
    "PrintVals": [
        "YeeHaw",
        "howdyPartner",
        "goodbye_buddy"
    ]
}
```

```
cd genplate
go generate ./...
# generatedcowboy.go and generatedrobot.go now exist
```

### Command line example

```
cd genplate
genplate example.template generatedcowboy.go cowboylingo.json
# generatedcowboy.go now exists
```

## Built-in template funcs

Add your custom funcs to a fork, or open a PR if they are sufficiently generic.

- `Pluralize`
- `CamelCase`
- `PascalCase`
- `SnakeCase`

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/urfave/cli"
)

type FieldMergeData struct {
	FieldName string
	ZeroValue string
}

type StructMergeData struct {
	StructName string
	Fields     []FieldMergeData
}

type StructMergeDataSlice []StructMergeData

func (x StructMergeDataSlice) Len() int      { return len(x) }
func (x StructMergeDataSlice) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x StructMergeDataSlice) Less(i, j int) bool {
	return x[i].StructName < x[j].StructName
}

type PackageStructMergeData struct {
	Package string
	Imports map[string]bool
	Structs StructMergeDataSlice
}

// TODO(pratik): Refactor ZeroValue to just use a conditional.
const structMerge = `package {{.Package}}

// This file is generated by the gen-merge package.
// Please avoid editing this manually.
import ({{range $k, $v  := .Imports}}
 {{$k}} {{end}})

{{range .Structs}}
// For any field not defined on s1, set it to the corresponding field in s2
func (s1 {{.StructName}}) Merge(s2 {{.StructName}}) {{.StructName}} {
{{range .Fields}}
	if s1.{{.FieldName}} == {{.ZeroValue}} {
		s1.{{.FieldName}} = s2.{{.FieldName}}
	}
{{end}}
	return s1
}

// For any field defined on s2, set it to the corresponding field in s1
func (s1 {{.StructName}}) MergeOverride(s2 {{.StructName}}) {{.StructName}} {
{{range .Fields}}
	if s2.{{.FieldName}} != {{.ZeroValue}} {
		s1.{{.FieldName}} = s2.{{.FieldName}}
	}
{{end}}
	return s1
}
{{end}}`

func PrintMergePackage(writer io.Writer, pkg PackageStructMergeData) error {
	pkgStructs := pkg.Structs
	sort.Sort(pkgStructs)

	structTemplate := template.New("Struct Template")
	// Filter out structs that don't have any fields so it's more minimal.
	structs := StructMergeDataSlice{}
	for _, p := range pkgStructs {
		if len(p.Fields) != 0 {
			structs = append(structs, p)
		}
	}
	pkg.Structs = structs
	structTemplate, err := structTemplate.Parse(structMerge)
	if err != nil {
		return err
	}

	return structTemplate.Execute(writer, pkg)
}

func ParsePackage(filename string) (PackageStructMergeData, error) {
	fset := token.NewFileSet() // positions are relative to fset
	pkg := PackageStructMergeData{}
	parsedPkgs, err := parser.ParseDir(fset, filename, nil, 0)
	if err != nil {
		return PackageStructMergeData{}, err
	}
	for packageName, parsedPkg := range parsedPkgs {
		pkg.Package = packageName
		for filePath, f := range parsedPkg.Files {

			if strings.Contains(filePath, "_test") {
				continue
			}

			for _, i := range f.Imports {
				alias := i.Name
				if alias != nil {
					if pkg.Imports == nil {
						pkg.Imports = map[string]bool{}
					}
					complete := alias.Name + " " + i.Path.Value
					pkg.Imports[complete] = true
				}
			}

			for name, object := range f.Scope.Objects {
				if !ast.IsExported(name) {
					continue
				}
				s := StructMergeData{
					StructName: name,
				}
				fields := []FieldMergeData{}
				switch object.Decl.(type) {
				case *ast.TypeSpec:
					typeSpec := object.Decl.(*ast.TypeSpec)
					switch typeSpec.Type.(type) {
					case *ast.StructType:

						// TODO: REMOVE
						//_ = typeSpec.Type
						//fmt.Printf("typeSpec.Type: %+v\n", typeSpec.Type.(*ast.InterfaceType))
						// TODO: REMOVE
						fieldsList := typeSpec.Type.(*ast.StructType).Fields.List
						for _, fieldDecl := range fieldsList {
							var zeroValue, currentFieldType string
							currentFieldName := fieldDecl.Names[0].Name
							switch fieldDecl.Type.(type) {
							case *ast.Ident:
								currentFieldType = fieldDecl.Type.(*ast.Ident).Name
								zeroValue = fmt.Sprintf("*(new(%s))", currentFieldType)
							case *ast.StarExpr:
								zeroValue = "nil"
							case *ast.ArrayType:
								zeroValue = "nil" // Nil is different from empty array.
							case *ast.SelectorExpr:
								selectorExpr := fieldDecl.Type.(*ast.SelectorExpr)
								zeroValue = fmt.Sprintf("*(new(%s.%s))", selectorExpr.X, selectorExpr.Sel.Name)
							}
							fields = append(fields, FieldMergeData{
								FieldName: currentFieldName,
								ZeroValue: zeroValue,
							})
						}
					}
					s.Fields = fields
					pkg.Structs = append(pkg.Structs, s)
				}
			}
		}
	}
	return pkg, nil
}

/* Example usage of PrintMergePackage
pkg := PackageStructMergeData{
	Package: "structs",
	Structs:[]StructMergeData{
	{
		StructName: "Person",
		Fields: []FieldMergeData{
			{
				FieldName: "Age",
				ZeroValue: "0",
			},
		},
	},
},
}
PrintMergePackage(os.Stdout, pkg)

will generate the following:

---------------------------------------------
package structs

func (s1 Person) Merge(s2 Person) Person {
  if s1.Age == 0 {
  	s1.Age = s2.Age
  }
  return s1
}
---------------------------------------------
*/

func main() {
	app := cli.NewApp()
	app.Name = "go-merge"
	app.Usage = "Generate some merge functions"
	app.Action = func(c *cli.Context) error {
		filename := c.Args().Get(0)
		pkg, err := ParsePackage(filename)
		if err != nil {
			return err
		}
		return PrintMergePackage(os.Stdout, pkg)
	}
	app.Run(os.Args)
}

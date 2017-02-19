package main

import (
	"fmt"
	"os"

	"go/ast"
	"go/parser"
	"go/token"
	"io"
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

type PackageStructMergeData struct {
	Package string
	Structs []StructMergeData
}

// TODO(pratik): Refactor ZeroValue to just use a conditional.
const structMerge = `package {{.Package}}

// This file is generated by the gen-merge package.
// Please avoid editing this manually.

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
{{end}}
`

func PrintMergePackage(writer io.Writer, pkg PackageStructMergeData) {
	structTemplate := template.New("Struct Template")
	structTemplate, err := structTemplate.Parse(structMerge)
	if err != nil {
		panic(err)
	}

	err = structTemplate.Execute(writer, pkg)
	if err != nil {
		panic(err)
	}
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
		fset := token.NewFileSet() // positions are relative to fset
		pkg := PackageStructMergeData{}
		parsedPkgs, err := parser.ParseDir(fset, filename, nil, 0)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		for packageName, parsedPkg := range parsedPkgs {
			pkg.Package = packageName
			for _, f := range parsedPkg.Files {
				for name, object := range f.Scope.Objects {
					s := StructMergeData{
						StructName: name,
					}
					fields := []FieldMergeData{}
					switch object.Decl.(type) {
					case *ast.TypeSpec:
						typeSpec := object.Decl.(*ast.TypeSpec)
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
							}
							fields = append(fields, FieldMergeData{
								FieldName: currentFieldName,
								ZeroValue: zeroValue,
							})
						}
						s.Fields = fields
						pkg.Structs = append(pkg.Structs, s)
					}
				}
			}
		}
		PrintMergePackage(os.Stdout, pkg)
		return nil
	}
	app.Run(os.Args)
}

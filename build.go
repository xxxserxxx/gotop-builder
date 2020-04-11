package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	help := flag.Bool("h", false, "Print help")
	rev := flag.String("r", "v3.5.1", "Tagged gotop version")
	flag.Parse()
	if *help {
		fmt.Printf("USAGE: %s -r <tag> <extension-URL>...\n", os.Args[0])
		flag.Usage()
		fmt.Print(`
This program creates the files necessary to compile gotop with selected
extensions. You need:

1. The tagged version of gotop to compile (e.g., "v3.5.1")
2. One (or more) extensions to enable (e.g. "github.com/xxxserxxx/gotop-nvidia")

Example:

$ go run ./build.go -r v3.5.1 github.com/xxxserxxx/gotop-dummy
$ go build -o gotop ./gotop.go
$ sudo cp gotop /usr/local/bin
`)
		os.Exit(0)
	}

	//////////////////////////////////
	// fetch gotop's main.go
	resp, err := http.Get("https://raw.githubusercontent.com/xxxserxxx/gotop/" + *rev + "/cmd/gotop/main.go")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//////////////////////////////////
	// add imports for the extensions
	fset := token.NewFileSet() // positions are relative to fset
	// Parse src but stop after processing the imports.
	f, err := parser.ParseFile(fset, "gotop.go", string(bs), parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, k := range flag.Args() {
		AddNamedImport(fset, f, "_", k)
	}

	//////////////////////////////////
	// write out the program and go mod file
	fout, err := os.Create("gotop.go")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fout.Close()
	err = printer.Fprint(fout, fset, f)

	mods := `module gotop

require github.com/xxxserxxx/gotop/%s %s

go 1.14
`
	gm, err := os.Create("go.mod")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer gm.Close()
	major := strings.Split(*rev, ".")[0]
	fmt.Fprintf(gm, mods, major, *rev)
}

// Most of this copied from
// https://github.com/golang/tools/blob/master/go/ast/astutil/imports.go
// That had features not needed for this job; this is greatly trimmed down.
// A little copying is better than a little dependency.
func AddNamedImport(fset *token.FileSet, f *ast.File, name, path string) {
	newImport := &ast.ImportSpec{
		Name: &ast.Ident{Name: name},
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: strconv.Quote(path),
		},
	}

	var (
		lastImport = -1         // index in f.Decls of the file's final import decl
		impDecl    *ast.GenDecl // import decl containing the best match
	)
	for i, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if ok && gen.Tok == token.IMPORT {
			lastImport = i
			impDecl = gen
			break
		}
	}
	impDecl = &ast.GenDecl{
		Tok: token.IMPORT,
	}
	impDecl.TokPos = f.Decls[lastImport].End()
	f.Decls = append(f.Decls, nil)
	copy(f.Decls[lastImport+2:], f.Decls[lastImport+1:])
	f.Decls[lastImport+1] = impDecl

	// Insert new import at insertAt.
	insertAt := 0
	impDecl.Specs = append(impDecl.Specs, nil)
	copy(impDecl.Specs[insertAt+1:], impDecl.Specs[insertAt:])
	impDecl.Specs[insertAt] = newImport
	pos := impDecl.Pos()
	newImport.Name.NamePos = pos
	newImport.Path.ValuePos = pos
	newImport.EndPos = pos

	// Clean up parens. impDecl contains at least one spec.
	if len(impDecl.Specs) == 1 {
		// Remove unneeded parens.
		impDecl.Lparen = token.NoPos
	} else if !impDecl.Lparen.IsValid() {
		// impDecl needs parens added.
		impDecl.Lparen = impDecl.Specs[0].Pos()
	}

	f.Imports = append(f.Imports, newImport)
}

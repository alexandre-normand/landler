package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	app = kingpin.New("landler", "Tool to list all http handlers in the given directory")
)

type errWriter struct {
	b   *bytes.Buffer
	err error
}

func (ew *errWriter) writeString(value string) {
	if ew.err != nil {
		return
	}
	_, ew.err = ew.b.WriteString(value)
}

func main() {
	kingpin.Version("1.0.0")
	kingpin.Parse()

	inputPaths := []string{"."}

	if handlers, err := run(inputPaths); err != nil {
		fmt.Fprintf(os.Stderr, "Error finding handlers: %v", err)
	} else {
		for _, handler := range handlers {
			fmt.Println(handler)
		}
	}
}

func run(inputPaths []string) (handlers []string, err error) {
	handlers = make([]string, 0)
	for _, path := range inputPaths {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			if isFile(file.Name()) && strings.HasSuffix(file.Name(), ".go") {
				fullPath := filepath.Join(path, file.Name())
				h, err := findFunctions(fullPath)
				if err != nil {
					return nil, err
				}

				handlers = append(handlers, h...)
			}
		}
	}

	return handlers, err
}

// isFile reports whether the named file is a file (not a directory).
func isFile(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return !info.IsDir()
}

func findFunctions(path string) (handlerNames []string, err error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}

	handlerNames = make([]string, 0)
	for _, decl := range f.Decls {
		if fun, ok := decl.(*ast.FuncDecl); ok {
			if isHttpHandler(fun) {
				handlerNames = append(handlerNames, fun.Name.Name)
			}
		}
	}

	return handlerNames, nil
}

func isHttpHandler(funcDeclaration *ast.FuncDecl) bool {
	if funcDeclaration.Recv != nil {
		return false
	} else if funcDeclaration.Type.Results != nil {
		return false
	} else if funcDeclaration.Type.Params == nil {
		return false
	} else if len(funcDeclaration.Type.Params.List) == 2 {
		return true
	}

	return false
}

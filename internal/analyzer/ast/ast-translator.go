package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

var logger slog.Logger

type PathToFile string
type Repository struct {
	files map[PathToFile]ast.File
}
type Repositories map[string]Repository

func Translate(services ...Service) Repositories {
	repos := make(Repositories)
	for _, service := range services {
		files, err := recTraversal(service.root)
		if err != nil {
			logger.Error("Can't get files from directory")
			return nil
		}
		repository := createAstPresentation(files)
		repos[service.Name+"/"+service.Version] = repository
	}
	return repos
}
func recTraversal(rootPath string) (*[]os.File, error) {
	var result []os.File

	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".go") {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			result = append(result, *file)
		}
		return nil
	})

	if err != nil {
		logger.Error("Can't traverse directory:" + err.Error())
		return nil, err
	}
	return &result, nil
}

func createAstPresentation(files *[]os.File) Repository {
	var repository Repository
	for _, file := range *files {
		src, err := io.ReadAll(&file)
		if err != nil {
			logger.Error("Can't read file:" + err.Error())
			continue
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, file.Name(), src, 0)
		if err != nil {
			logger.Error("Can't parse file:" + err.Error())
			continue
		}
		repository.files[PathToFile(file.Name())] = *f
	}
	return repository
}

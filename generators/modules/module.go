package modules

import (
	"encoding/json"
	"fmt"

	"github.com/cihangir/gene/generators/clients"
	"github.com/cihangir/gene/generators/errors"
	"github.com/cihangir/gene/generators/folders"
	"github.com/cihangir/gene/generators/handlers"
	"github.com/cihangir/gene/generators/models"
	"github.com/cihangir/gene/generators/tests"
	"github.com/cihangir/gene/helpers"
	"github.com/cihangir/gene/schema"
	"github.com/cihangir/gene/stringext"
)

type Module struct {
	schema           *schema.Schema
	TargetFolderName string
}

func NewModule(s *schema.Schema) *Module {
	return &Module{
		schema:           s,
		TargetFolderName: "./",
	}
}

func NewFromFile(path string) (*Module, error) {
	fileContent, err := helpers.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var s schema.Schema
	if err := json.Unmarshal(fileContent, &s); err != nil {
		return nil, err
	}

	return NewModule(&s), nil
}

func (m *Module) Create() error {
	rootPath := m.TargetFolderName

	// first ensure that we have the correct folder structure for our system
	if err := folders.EnsureFolders(
		rootPath, // root folder
		folders.FolderStucture,
	); err != nil {
		return err
	}

	// create the module folder structure
	if err := folders.EnsureFolders(
		rootPath, // root folder
		createModuleStructure(stringext.ToLowerFirst(
			m.schema.Title,
		)),
	); err != nil {
		return err
	}

	if err := models.Generate(rootPath, m.schema); err != nil {
		return err
	}

	if err := models.GenerateStatements(rootPath, m.schema); err != nil {
		return err
	}

	if err := handlers.Generate(rootPath, m.schema); err != nil {
		return err
	}

	if err := errors.Generate(rootPath, m.schema); err != nil {
		return err
	}

	if err := m.GenerateMainFile(rootPath); err != nil {
		return err
	}

	if err := clients.Generate(rootPath, m.schema); err != nil {
		return err
	}

	if err := tests.Generate(rootPath, m.schema); err != nil {
		return err
	}

	return nil
}

var moduleFolderStucture = []string{
	"workers/%[1]s",
	"workers/%[1]s/%[1]sapi",
	"workers/%[1]s/cmd",
	"workers/%[1]s/cmd/%[1]s",
	"workers/%[1]s/tests",
	"workers/%[1]s/errors",
	"workers/%[1]s/clients",
	// "workers/%[1]s/handlers",
}

func createModuleStructure(name string) []string {
	modified := make([]string, len(moduleFolderStucture))
	for i, str := range moduleFolderStucture {
		modified[i] = fmt.Sprintf(str, name)
	}

	return modified
}

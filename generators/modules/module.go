// Package modules handle module creation for the given json-schema
package modules

import (
	"fmt"

	"github.com/cihangir/gene/config"
	"github.com/cihangir/gene/generators/clients"
	"github.com/cihangir/gene/generators/common"
	gerr "github.com/cihangir/gene/generators/errors"
	"github.com/cihangir/gene/generators/functions"
	"github.com/cihangir/gene/generators/mainfile"
	"github.com/cihangir/gene/generators/models"
	"github.com/cihangir/gene/generators/scanners/rows"
	"github.com/cihangir/gene/generators/sql/definitions"
	"github.com/cihangir/gene/generators/sql/statements"

	"github.com/cihangir/schema"
)

type Generator interface {
	Name() string
	Generate(*common.Context, *schema.Schema) ([]common.Output, error)
}

var generators []Generator

func init() {
	generators = []Generator{
		statements.New(),
		models.New(),
		rows.New(),
		gerr.New(),
		mainfile.New(),
		clients.New(),
		// tests.New(), //TODO(cihangir) tests are not stable
		functions.New(),
		definitions.New(),
		// js.New(),
	}
}

// Module holds the required parameters for a module
type Module struct {
	schema *schema.Schema

	context *common.Context

	// TargetFolderName holds the folder name for the module
	TargetFolderName string
}

// NewModule creates a new module with the given Schema
func New(conf *config.Config) (*Module, error) {
	s, err := common.Read(conf.Schema)
	if err != nil {
		return nil, err
	}

	context := common.NewContext()
	context.Config = conf

	return &Module{
		schema:           s.Resolve(s),
		context:          context,
		TargetFolderName: "./",
	}, nil
}

// Create creates the module. While creating the module it handles models,
// handlers, errors, servers, clients and tests generation
func (m *Module) Create() error {

	for _, gen := range generators {
		mgen, err := gen.Generate(m.context, m.schema)
		if err != nil {
			return err
		}

		if err := common.WriteOutput(mgen); err != nil {
			return err
		}
	}

	return nil
}

var moduleFolderStucture = []string{
	"cmd/%[1]s/",
	"workers/%[1]s",
	"workers/%[1]s/api",
	"workers/%[1]s/tests",
	"workers/%[1]s/js",
	"workers/%[1]s/errors",
	"workers/%[1]s/clients",
}

func createModuleStructure(name string) []string {
	modified := make([]string, len(moduleFolderStucture))
	for i, str := range moduleFolderStucture {
		modified[i] = fmt.Sprintf(str, name)
	}

	return modified
}

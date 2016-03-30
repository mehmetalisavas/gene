package js

import (
	"fmt"
	"strings"

	"github.com/cihangir/gene/generators/common"
)

// Generator for js client
type Generator struct{}

func pathfunc(data *common.TemplateData) string {
	return fmt.Sprintf(
		"%s/%s/%s.js",
		data.Settings.Get("fullPathPrefix").(string),
		data.ModuleName,
		strings.ToLower(data.Schema.Title),
	)

}

// Generate generates JS client for given schema
func (g *Generator) Generate(req *common.Req, res *common.Res) error {
	o := &common.Op{
		Name:        "js",
		Template:    FunctionsTemplate,
		PathFunc:    pathfunc,
		DoNotFormat: true,
	}

	return common.Proces(o, req, res)
}

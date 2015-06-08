package rows

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"go/format"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/schema"
)

type Generator struct{}

// Generate generates and writes the errors of the schema
func (g *Generator) Generate(context *common.Context, s *schema.Schema) ([]common.Output, error) {
	temp := template.New("rowscanner.tmpl").Funcs(common.TemplateFuncs)
	if _, err := temp.Parse(RowScannerTemplate); err != nil {
		return nil, err
	}

	outputs := make([]common.Output, 0)
	for _, def := range common.SortedObjectSchemas(s.Definitions) {

		var buf bytes.Buffer

		data := struct {
			Schema *schema.Schema
		}{
			Schema: def,
		}

		if err := temp.ExecuteTemplate(&buf, "rowscanner.tmpl", data); err != nil {
			return nil, err
		}

		f, err := format.Source(buf.Bytes())
		if err != nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		outputs = append(outputs, common.Output{
			Content: f,
			Path: fmt.Sprintf(
				"%s%s_rowscanner.go",
				context.Config.Target,
				strings.ToLower(def.Title),
			),
		})
	}

	return outputs, nil
}

package definitions

import (
	"encoding/json"

	"testing"

	"github.com/cihangir/gene/config"
	"github.com/cihangir/gene/testdata"
	"github.com/cihangir/schema"
	"github.com/cihangir/stringext"
)

func TestTypes(t *testing.T) {
	s := &schema.Schema{}
	if err := json.Unmarshal([]byte(testdata.TestDataFull), s); err != nil {
		t.Fatal(err.Error())
	}

	s = s.Resolve(s)
	g := New()

	context := config.NewContext()
	moduleName := context.ModuleNameFunc(s.Title)
	settings := g.generateSettings(moduleName, s)

	index := 0
	for _, def := range s.Definitions {

		// schema should have our generator
		if !def.Generators.Has(generatorName) {
			continue
		}

		settingsDef := g.setDefaultSettings(settings, def)
		settingsDef.Set("tableName", stringext.ToFieldName(def.Title))

		sts := DefineTypes(settingsDef, def)
		equals(t, expectedTypes[index], sts)
		index++
	}
}

var expectedTypes = []string{
	`
-- ----------------------------
--  Types structure for account.profile.enum_bare
-- ----------------------------
DROP TYPE IF EXISTS "account"."profile_enum_bare_enum" CASCADE;
CREATE TYPE "account"."profile_enum_bare_enum" AS ENUM (
  'enum1',
  'enum2',
  'enum3'
);
ALTER TYPE "account"."profile_enum_bare_enum" OWNER TO "social";`,
}

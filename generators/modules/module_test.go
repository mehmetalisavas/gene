package modules

import (
	"encoding/json"

	"github.com/cihangir/gene/schema"
	"github.com/cihangir/gene/testdata"

	"testing"
)

func TestCreateModuleStructure(t *testing.T) {
	expected := []string{
		"workers/name",
		"workers/name/api",
		"workers/name/cmd",
		"workers/name/cmd/name",
		"workers/name/tests",
		"workers/name/errors",
		"workers/name/handlers",
	}

	structure := createModuleStructure("name")

	for _, stc := range structure {
		exists := false
		for _, expt := range expected {
			if expt == stc {
				exists = true
				break
			}
		}
		if !exists {
			t.Fatalf("%s is not expected in the result set", stc)
		}
	}
}

func TestCreateModule(t *testing.T) {
	var s schema.Schema
	if err := json.Unmarshal([]byte(testdata.JSON1), &s); err != nil {
		t.Fatal(err.Error())
	}

	m := NewModule(&s)
	err := m.Create()
	if err != nil {
		t.Fatal(err.Error())
	}

}

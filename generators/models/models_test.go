package models

import (
	"encoding/json"
	"strings"

	"testing"

	"github.com/cihangir/gene/generators/folders"
	"github.com/cihangir/gene/generators/validators"
	"github.com/cihangir/gene/testdata"
	"github.com/cihangir/schema"

	"github.com/cihangir/gene/writers"
)

func TestGenerateModel(t *testing.T) {
	var s schema.Schema
	if err := json.Unmarshal([]byte(testdata.JSON1), &s); err != nil {
		t.Fatal(err.Error())
	}

	model, err := GenerateModel(&s)
	if err != nil {
		t.Fatal(err.Error())
	}

	folders.EnsureFolders("/tmp/", folders.FolderStucture)
	fileName := "/tmp/models/" + s.Title + ".go"

	err = writers.WriteFormattedFile(fileName, model)
	if err != nil {
		t.Fatal(err.Error())
	}

}

func TestGenerateSchema(t *testing.T) {
	s := &schema.Schema{}
	if err := json.Unmarshal([]byte(testdata.JSON1), s); err != nil {
		t.Fatal(err.Error())
	}

	// replace "~" with "`"
	result := strings.Replace(`
// Account represents a registered User
type Account struct {
	CreatedAt              time.Time ~json:"createdAt"~              // Profile's creation time
	EmailAddress           string    ~json:"emailAddress"~           // Email Address of the Account
	EmailStatusConstant    string    ~json:"emailStatusConstant"~    // Status of the Account's Email
	ID                     int64     ~json:"id"~                     // The unique identifier for a Account's Profile
	Password               string    ~json:"password"~               // Salted Password of the Account
	PasswordStatusConstant string    ~json:"passwordStatusConstant"~ // Status of the Account's Password
	ProfileID              int64     ~json:"profileId"~              // The unique identifier for a Account's Profile
	Salt                   string    ~json:"salt"~                   // Salt used to hash Password of the Account
	StatusConstant         string    ~json:"statusConstant"~         // Status of the Account
	URL                    string    ~json:"url"~                    // Salted Password of the Account
	URLName                string    ~json:"urlName"~                // Salted Password of the Account
}`, "~", "`", -1)

	code, err := GenerateSchema(s.Definitions["Account"])
	if err != nil {
		t.Fatal(err.Error())
	}

	if result != string(code) {
		// fmt.Printf("foo %# v", pretty.Formatter(difflib.Diff([]string{result}, []string{string(code)})))
		t.Fatalf("Schema is not same. Wanted: %s, Get: %s", result, string(code))
	}
}

func TestGenerateValidators(t *testing.T) {
	s := &schema.Schema{}
	if err := json.Unmarshal([]byte(testdata.JSON1), s); err != nil {
		t.Fatal(err.Error())
	}
	result := `
// Validate validates the struct
func (a *Account) Validate() error {
	return govalidator.NewMulti(govalidator.Date(a.CreatedAt),
		govalidator.MaxLength(a.Salt, 255),
		govalidator.Min(float64(a.ID), 1.000000),
		govalidator.Min(float64(a.ProfileID), 1.000000),
		govalidator.MinLength(a.Password, 6),
		govalidator.MinLength(a.URL, 6),
		govalidator.MinLength(a.URLName, 6),
		govalidator.OneOf(a.EmailStatusConstant, []string{
			EmailStatusConstant.Verified,
			EmailStatusConstant.NotVerified,
		}),
		govalidator.OneOf(a.PasswordStatusConstant, []string{
			PasswordStatusConstant.Valid,
			PasswordStatusConstant.NeedsReset,
			PasswordStatusConstant.Generated,
		}),
		govalidator.OneOf(a.StatusConstant, []string{
			StatusConstant.Registered,
			StatusConstant.Unregistered,
			StatusConstant.NeedsManualVerification,
		})).Validate()
}`

	code, err := validators.Generate(s.Definitions["Account"])
	if err != nil {
		t.Fatal(err.Error())
	}

	if result != string(code) {
		t.Fatalf("Schema is not same. Wanted: %s, Get: %s", result, string(code))
	}
}

func TestGenerateFunctions(t *testing.T) {
	var s schema.Schema
	if err := json.Unmarshal([]byte(testdata.JSON1), &s); err != nil {
		t.Fatal(err.Error())
	}

	_, err := GenerateFunctions(&s)
	if err != nil {
		t.Fatal(err.Error())
	}
}

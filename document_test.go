package aeo

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func loadFixture(t *testing.T) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "aeo-person.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	return data
}

func TestParseDocument_CanonicalPersonExample(t *testing.T) {
	doc, err := ParseDocument(loadFixture(t))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if doc.AEOVersion != "0.1" {
		t.Errorf("AEOVersion = %q, want %q", doc.AEOVersion, "0.1")
	}
	if doc.Entity.Type != EntityPerson {
		t.Errorf("Entity.Type = %q, want %q", doc.Entity.Type, EntityPerson)
	}
	if doc.Entity.Name != "Miz Causevic" {
		t.Errorf("Entity.Name = %q, want %q", doc.Entity.Name, "Miz Causevic")
	}
	if got := len(doc.Claims); got != 6 {
		t.Errorf("len(Claims) = %d, want 6", got)
	}
}

func TestClaimIDs_AllPresent(t *testing.T) {
	doc, err := ParseDocument(loadFixture(t))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	got := doc.ClaimIDs()
	sort.Strings(got)
	want := []string{
		"authored-spec",
		"current-role",
		"live-products",
		"location",
		"primary-stack",
		"years-experience",
	}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i, id := range want {
		if got[i] != id {
			t.Errorf("ID[%d] = %q, want %q", i, got[i], id)
		}
	}
}

func TestFindClaim_Found(t *testing.T) {
	doc, _ := ParseDocument(loadFixture(t))
	claim := doc.FindClaim("years-experience")
	if claim == nil {
		t.Fatal("FindClaim returned nil for present ID")
	}
	if claim.Predicate != "aeo:yearsOfExperience" {
		t.Errorf("Predicate = %q, want %q", claim.Predicate, "aeo:yearsOfExperience")
	}
	if v, ok := claim.Value.(float64); !ok || v != 30 {
		t.Errorf("Value = %v (%T), want 30 (float64)", claim.Value, claim.Value)
	}
}

func TestFindClaim_Missing(t *testing.T) {
	doc, _ := ParseDocument(loadFixture(t))
	if doc.FindClaim("does-not-exist") != nil {
		t.Error("FindClaim should return nil for unknown ID")
	}
}

func TestRoundTrip_PreservesStructure(t *testing.T) {
	doc, err := ParseDocument(loadFixture(t))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	reSerialized, err := doc.Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	reParsed, err := ParseDocument(reSerialized)
	if err != nil {
		t.Fatalf("re-parse: %v", err)
	}
	if reParsed.Entity.Name != doc.Entity.Name {
		t.Error("entity name not preserved")
	}
	if len(reParsed.Claims) != len(doc.Claims) {
		t.Errorf("claim count: got %d, want %d", len(reParsed.Claims), len(doc.Claims))
	}
}

func TestParseDocument_RejectsUnknownField(t *testing.T) {
	bad := strings.Replace(
		string(loadFixture(t)),
		`"aeo_version": "0.1"`,
		`"aeo_version": "0.1", "unexpected_field": "nope"`,
		1,
	)
	if _, err := ParseDocumentString(bad); err == nil {
		t.Error("expected error on unknown top-level field, got nil")
	}
}

func TestLoadDocument_FromFile(t *testing.T) {
	doc, err := LoadDocument(filepath.Join("testdata", "aeo-person.json"))
	if err != nil {
		t.Fatalf("LoadDocument: %v", err)
	}
	if doc.Entity.Name != "Miz Causevic" {
		t.Errorf("Entity.Name = %q, want %q", doc.Entity.Name, "Miz Causevic")
	}
}

// Package aeo provides parse, build, and serialization support for
// AEO Protocol v0.1 declaration documents.
//
// Specification: https://github.com/mizcausevic-dev/aeo-protocol-spec
package aeo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

// ProtocolVersion is the AEO Protocol version this SDK targets.
const ProtocolVersion = "0.1"

// SDKVersion is the version of this SDK module.
const SDKVersion = "0.1.0"

// EntityType enumerates the entity kinds the AEO Protocol defines.
type EntityType string

// Permitted EntityType values.
const (
	EntityPerson       EntityType = "Person"
	EntityOrganization EntityType = "Organization"
	EntityProduct      EntityType = "Product"
	EntityPlace        EntityType = "Place"
	EntityConcept      EntityType = "Concept"
)

// VerificationType enumerates the verification mechanisms an authority block can declare.
type VerificationType string

// Permitted VerificationType values.
const (
	VerificationDomain       VerificationType = "domain"
	VerificationDNS          VerificationType = "dns"
	VerificationGitHub       VerificationType = "github"
	VerificationLinkedIn     VerificationType = "linkedin"
	VerificationGPG          VerificationType = "gpg"
	VerificationWellKnownURI VerificationType = "well-known-uri"
)

// Confidence enumerates a claim's confidence level.
type Confidence string

// Permitted Confidence values.
const (
	ConfidenceHigh   Confidence = "high"
	ConfidenceMedium Confidence = "medium"
	ConfidenceLow    Confidence = "low"
)

// AuditMode enumerates the document audit modes.
type AuditMode string

// Permitted AuditMode values.
const (
	AuditNone      AuditMode = "none"
	AuditSignature AuditMode = "signature"
	AuditEndpoint  AuditMode = "endpoint"
)

// Entity is the subject of an AEO declaration.
type Entity struct {
	ID           string     `json:"id"`
	Type         EntityType `json:"type"`
	Name         string     `json:"name"`
	Aliases      []string   `json:"aliases,omitempty"`
	CanonicalURL string     `json:"canonical_url"`
}

// Verification is a proof of ownership or control over an identifier.
type Verification struct {
	Type     VerificationType `json:"type"`
	Value    string           `json:"value"`
	ProofURI string           `json:"proof_uri,omitempty"`
}

// Authority captures what the entity considers authoritative about itself.
type Authority struct {
	PrimarySources []string       `json:"primary_sources"`
	EvidenceLinks  []string       `json:"evidence_links,omitempty"`
	Verifications  []Verification `json:"verifications,omitempty"`
}

// Claim is a single asserted fact about the entity.
type Claim struct {
	ID         string      `json:"id"`
	Predicate  string      `json:"predicate"`
	Value      interface{} `json:"value"`
	Evidence   []string    `json:"evidence,omitempty"`
	ValidFrom  string      `json:"valid_from,omitempty"`
	ValidUntil *string     `json:"valid_until,omitempty"`
	Confidence Confidence  `json:"confidence,omitempty"`
}

// CitationPreferences expresses how the entity prefers to be cited.
type CitationPreferences struct {
	PreferredAttribution string   `json:"preferred_attribution,omitempty"`
	CanonicalLinks       []string `json:"canonical_links,omitempty"`
	DoNotCite            []string `json:"do_not_cite,omitempty"`
}

// AnswerConstraints expresses soft constraints for answer engines.
type AnswerConstraints struct {
	MustInclude         []string `json:"must_include,omitempty"`
	MustNotInclude      []string `json:"must_not_include,omitempty"`
	FreshnessWindowDays int      `json:"freshness_window_days,omitempty"`
}

// Audit configures the audit surface for the declaration.
type Audit struct {
	Mode           AuditMode `json:"mode"`
	SigningKeyURI  string    `json:"signing_key_uri,omitempty"`
	Signature      string    `json:"signature,omitempty"`
	EndpointURI    string    `json:"endpoint_uri,omitempty"`
	EndpointSchema string    `json:"endpoint_schema,omitempty"`
}

// Document is a complete AEO Protocol v0.1 declaration document.
type Document struct {
	AEOVersion          string               `json:"aeo_version"`
	Entity              Entity               `json:"entity"`
	Authority           Authority            `json:"authority"`
	Claims              []Claim              `json:"claims"`
	CitationPreferences *CitationPreferences `json:"citation_preferences,omitempty"`
	AnswerConstraints   *AnswerConstraints   `json:"answer_constraints,omitempty"`
	Audit               *Audit               `json:"audit,omitempty"`
}

// ParseDocument parses a byte slice into a Document. It rejects unknown
// top-level fields and unknown fields anywhere in the structure.
func ParseDocument(raw []byte) (*Document, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	var d Document
	if err := dec.Decode(&d); err != nil {
		return nil, fmt.Errorf("aeo: parse document: %w", err)
	}
	return &d, nil
}

// ParseDocumentString is a convenience wrapper around ParseDocument.
func ParseDocumentString(raw string) (*Document, error) {
	return ParseDocument([]byte(raw))
}

// LoadDocument reads and parses a document from a file path.
func LoadDocument(path string) (*Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("aeo: read %s: %w", path, err)
	}
	return ParseDocument(data)
}

// Marshal serializes the document to pretty-printed JSON.
func (d *Document) Marshal() ([]byte, error) {
	return json.MarshalIndent(d, "", "  ")
}

// ClaimIDs returns the IDs of all claims in the document.
func (d *Document) ClaimIDs() []string {
	ids := make([]string, len(d.Claims))
	for i, c := range d.Claims {
		ids[i] = c.ID
	}
	return ids
}

// FindClaim returns a pointer to the claim with the given ID, or nil.
func (d *Document) FindClaim(id string) *Claim {
	for i := range d.Claims {
		if d.Claims[i].ID == id {
			return &d.Claims[i]
		}
	}
	return nil
}

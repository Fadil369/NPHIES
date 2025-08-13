package models

import (
	"time"
)

// CodeSystem represents a terminology code system
type CodeSystem struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	URL         string    `json:"url" db:"url"`
	Version     string    `json:"version" db:"version"`
	Status      string    `json:"status" db:"status"` // active, inactive, draft
	Publisher   string    `json:"publisher" db:"publisher"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Code represents an individual code within a code system
type Code struct {
	ID           string    `json:"id" db:"id"`
	CodeSystemID string    `json:"code_system_id" db:"code_system_id"`
	Code         string    `json:"code" db:"code"`
	Display      string    `json:"display" db:"display"`
	Definition   string    `json:"definition" db:"definition"`
	Status       string    `json:"status" db:"status"` // active, inactive, deprecated
	Parent       *string   `json:"parent,omitempty" db:"parent"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// ValueSet represents a set of codes from one or more code systems
type ValueSet struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	URL         string    `json:"url" db:"url"`
	Version     string    `json:"version" db:"version"`
	Status      string    `json:"status" db:"status"`
	Publisher   string    `json:"publisher" db:"publisher"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ValueSetCode represents a code included in a value set
type ValueSetCode struct {
	ID           string `json:"id" db:"id"`
	ValueSetID   string `json:"value_set_id" db:"value_set_id"`
	CodeSystemID string `json:"code_system_id" db:"code_system_id"`
	Code         string `json:"code" db:"code"`
	Display      string `json:"display" db:"display"`
}

// CodeMapping represents mapping between codes from different systems
type CodeMapping struct {
	ID               string    `json:"id" db:"id"`
	SourceSystem     string    `json:"source_system" db:"source_system"`
	SourceCode       string    `json:"source_code" db:"source_code"`
	TargetSystem     string    `json:"target_system" db:"target_system"`
	TargetCode       string    `json:"target_code" db:"target_code"`
	Equivalence      string    `json:"equivalence" db:"equivalence"` // equivalent, equal, wider, subsumes, narrower, specializes, inexact, unmatched, disjoint
	Comment          string    `json:"comment" db:"comment"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// DTOs for API requests and responses

type CreateCodeSystemRequest struct {
	Name        string `json:"name" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	URL         string `json:"url" binding:"required"`
	Version     string `json:"version" binding:"required"`
	Publisher   string `json:"publisher"`
}

type CreateCodeRequest struct {
	CodeSystemID string `json:"code_system_id" binding:"required"`
	Code         string `json:"code" binding:"required"`
	Display      string `json:"display" binding:"required"`
	Definition   string `json:"definition"`
	Parent       string `json:"parent,omitempty"`
}

type CreateValueSetRequest struct {
	Name        string `json:"name" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	URL         string `json:"url" binding:"required"`
	Version     string `json:"version" binding:"required"`
	Publisher   string `json:"publisher"`
	Codes       []ValueSetCodeRequest `json:"codes"`
}

type ValueSetCodeRequest struct {
	CodeSystemID string `json:"code_system_id" binding:"required"`
	Code         string `json:"code" binding:"required"`
	Display      string `json:"display"`
}

type CodeLookupRequest struct {
	System string `json:"system" binding:"required"`
	Code   string `json:"code" binding:"required"`
}

type CodeLookupResponse struct {
	Found       bool   `json:"found"`
	Code        string `json:"code,omitempty"`
	Display     string `json:"display,omitempty"`
	Definition  string `json:"definition,omitempty"`
	System      string `json:"system,omitempty"`
	SystemName  string `json:"system_name,omitempty"`
}

type CodeValidationRequest struct {
	System    string `json:"system" binding:"required"`
	Code      string `json:"code" binding:"required"`
	ValueSet  string `json:"value_set,omitempty"`
}

type CodeValidationResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message,omitempty"`
	Display string `json:"display,omitempty"`
}

type CodeMappingRequest struct {
	SourceSystem string `json:"source_system" binding:"required"`
	SourceCode   string `json:"source_code" binding:"required"`
	TargetSystem string `json:"target_system" binding:"required"`
}

type CodeMappingResponse struct {
	Mappings []CodeMapping `json:"mappings"`
}

type ImportTerminologyRequest struct {
	Type   string `json:"type" binding:"required"` // "fhir", "csv", "xml"
	Source string `json:"source" binding:"required"` // URL or base64 encoded data
	Format string `json:"format"` // Additional format info
}

type StatisticsResponse struct {
	CodeSystems int64            `json:"code_systems"`
	Codes       int64            `json:"codes"`
	ValueSets   int64            `json:"value_sets"`
	Mappings    int64            `json:"mappings"`
	CacheStats  map[string]int64 `json:"cache_stats"`
}
// Package model defines shared types used across the differ subsystem.
package model

// AttributeDiff holds the before and after values for a single resource attribute.
type AttributeDiff struct {
	Before interface{}
	After  interface{}
}

// ChangeType classifies the kind of drift detected for a resource.
type ChangeType string

const (
	ChangeAdded    ChangeType = "added"
	ChangeRemoved  ChangeType = "removed"
	ChangeModified ChangeType = "modified"
)

// DriftItem represents a single resource that has drifted between two state files.
type DriftItem struct {
	// Address is the fully-qualified Terraform resource address, e.g. "aws_s3_bucket.my_bucket".
	Address string
	// Type is the Terraform resource type, e.g. "aws_s3_bucket".
	Type string
	// Name is the resource label within its type.
	Name string
	// ChangeType classifies the drift.
	ChangeType ChangeType
	// Attributes contains per-attribute before/after diffs; nil for added/removed resources.
	Attributes map[string]AttributeDiff
}

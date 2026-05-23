package differ

import (
	"strings"

	"github.com/user/stackdiff/internal/differ/model"
)

// IgnoreList holds rules for suppressing resources or attributes from drift output.
type IgnoreList struct {
	resourcePrefixes []string
	attributeKeys    []string
}

// NewIgnoreList constructs an IgnoreList from prefix and attribute key slices.
func NewIgnoreList(resourcePrefixes, attributeKeys []string) *IgnoreList {
	return &IgnoreList{
		resourcePrefixes: resourcePrefixes,
		attributeKeys:    attributeKeys,
	}
}

// SuppressResource returns true when the resource address matches any configured prefix.
func (il *IgnoreList) SuppressResource(address string) bool {
	for _, prefix := range il.resourcePrefixes {
		if strings.HasPrefix(address, prefix) {
			return true
		}
	}
	return false
}

// FilterAttributes removes attribute entries whose keys match any configured key rule.
func (il *IgnoreList) FilterAttributes(attrs map[string]model.AttributeDiff) map[string]model.AttributeDiff {
	if len(il.attributeKeys) == 0 {
		return attrs
	}
	out := make(map[string]model.AttributeDiff, len(attrs))
	for k, v := range attrs {
		if !il.matchesAttributeKey(k) {
			out[k] = v
		}
	}
	return out
}

// Apply filters a slice of DriftItems, removing suppressed resources and stripping ignored attributes.
func (il *IgnoreList) Apply(items []model.DriftItem) []model.DriftItem {
	result := make([]model.DriftItem, 0, len(items))
	for _, item := range items {
		if il.SuppressResource(item.Address) {
			continue
		}
		item.Attributes = il.FilterAttributes(item.Attributes)
		result = append(result, item)
	}
	return result
}

func (il *IgnoreList) matchesAttributeKey(key string) bool {
	for _, rule := range il.attributeKeys {
		if key == rule || strings.HasPrefix(key, rule+".") {
			return true
		}
	}
	return false
}

package differ

import "strings"

// IgnoreRule describes a single rule for suppressing drift items.
type IgnoreRule struct {
	// AddressPrefix, if non-empty, suppresses any resource whose address
	// starts with this string (e.g. "module.legacy").
	AddressPrefix string

	// AttributeKey, if non-empty, suppresses changes to this specific
	// attribute key across all resources (e.g. "tags.LastModified").
	AttributeKey string
}

// IgnoreList holds a collection of IgnoreRules.
type IgnoreList struct {
	rules []IgnoreRule
}

// NewIgnoreList constructs an IgnoreList from a slice of rules.
func NewIgnoreList(rules []IgnoreRule) *IgnoreList {
	return &IgnoreList{rules: rules}
}

// SuppressResource returns true when the given resource address matches
// any address-prefix rule.
func (il *IgnoreList) SuppressResource(address string) bool {
	for _, r := range il.rules {
		if r.AddressPrefix != "" && strings.HasPrefix(address, r.AddressPrefix) {
			return true
		}
	}
	return false
}

// FilterAttributes removes attribute keys that are matched by any
// attribute-key rule, returning a new map without those keys.
func (il *IgnoreList) FilterAttributes(attrs map[string]string) map[string]string {
	if len(attrs) == 0 {
		return attrs
	}
	out := make(map[string]string, len(attrs))
	for k, v := range attrs {
		if !il.suppressAttr(k) {
			out[k] = v
		}
	}
	return out
}

func (il *IgnoreList) suppressAttr(key string) bool {
	for _, r := range il.rules {
		if r.AttributeKey != "" && r.AttributeKey == key {
			return true
		}
	}
	return false
}

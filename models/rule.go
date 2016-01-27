// Copyright 2016 Eleme Inc. All rights reserved.

package models

import (
	"fmt"
	"github.com/eleme/banshee/config"
	"github.com/eleme/banshee/util"
	"strings"
)

// Rule is a type to describe alerter rule.
type Rule struct {
	// Rule may be cached.
	cache `sql:"-" json:"-"`
	// ID in db.
	ID int `gorm:"primary_key" json:"id"`
	// Project belongs to
	ProjectID int `sql:"index;not null" json:"projectID"`
	// Pattern is a wildcard string
	Pattern string `sql:"size:400;not null;unique" json:"pattern"`
	// Trend
	TrendUp   bool `json:trendUp`
	TrendDown bool `json:trendDown`
	// Optional thresholds data, only used if the rule condition is about
	// threshold. Although we really don't need any thresholds for trending
	// analyzation and alertings, but we still provide a way to alert by
	// thresholds.
	ThresholdMax float64 `json:"thresholdMax"`
	ThresholdMin float64 `json:"thresholdMin"`
	// String representation.
	Repr string `sql:"-" json:"repr"`
}

// Copy the rule.
func (rule *Rule) Copy() *Rule {
	dst := &Rule{}
	rule.CopyTo(dst)
	return dst
}

// CopyTo copy the rule to another.
func (rule *Rule) CopyTo(r *Rule) {
	rule.RLock()
	defer rule.RUnlock()
	r.Lock()
	defer r.Unlock()
	r.ID = rule.ID
	r.ProjectID = rule.ProjectID
	r.Pattern = rule.Pattern
	r.TrendUp = rule.TrendUp
	r.TrendDown = rule.TrendDown
	r.ThresholdMax = rule.ThresholdMax
	r.ThresholdMin = rule.ThresholdMin
}

// Equal tests rule equality
func (rule *Rule) Equal(r *Rule) bool {
	rule.RLock()
	defer rule.RUnlock()
	r.RLock()
	defer rule.RUnlock()
	return (r.ID == rule.ID &&
		r.ProjectID == rule.ProjectID &&
		r.Pattern == rule.Pattern &&
		r.TrendUp == rule.TrendUp &&
		r.TrendDown == rule.TrendDown &&
		r.ThresholdMax == rule.ThresholdMax &&
		r.ThresholdMin == rule.ThresholdMin)
}

// Test if a metric hits this rule.
//
//	1. For trend related conditions, index.Score will be used.
//	2. For value related conditions, metric.Value will be used.
//
func (rule *Rule) Test(m *Metric, idx *Index, cfg *config.Config) bool {
	// RLock if shared.
	rule.RLock()
	defer rule.RUnlock()
	// Conditions
	ok := false
	if !ok && rule.TrendUp && rule.ThresholdMax == 0 {
		// TrendUp
		ok = idx.Score > 1
	}
	if !ok && rule.TrendUp && rule.ThresholdMax != 0 {
		// TrendUp And Value >= X
		ok = idx.Score > 1 && m.Value >= rule.ThresholdMax
	}
	if !ok && !rule.TrendUp && rule.ThresholdMax != 0 {
		// Value >= X
		ok = m.Value >= rule.ThresholdMax
	}
	if !ok && rule.TrendDown && rule.ThresholdMin == 0 {
		// TrendDown
		ok = idx.Score < -1
	}
	if !ok && rule.TrendDown && rule.ThresholdMin != 0 {
		// TrendDown And Value <= X
		ok = idx.Score < -1 && m.Value <= rule.ThresholdMin
	}
	if !ok && !rule.TrendDown && rule.ThresholdMin != 0 {
		// Value <= X
		ok = m.Value <= rule.ThresholdMin
	}
	return ok
}

// BuildRepr initializes the rule's string repr.
func (rule *Rule) BuildRepr() {
	var parts []string
	if rule.TrendUp && rule.ThresholdMax == 0 {
		parts = append(parts, "trend ↑")
	}
	if rule.TrendUp && rule.ThresholdMax != 0 {
		s := fmt.Sprintf("(trend ↑ && value >= %s)", util.ToFixed(rule.ThresholdMax, 3))
		parts = append(parts, s)
	}
	if !rule.TrendUp && rule.ThresholdMax != 0 {
		s := fmt.Sprintf("value >= %s", util.ToFixed(rule.ThresholdMax, 3))
		parts = append(parts, s)
	}
	if rule.TrendDown && rule.ThresholdMin == 0 {
		parts = append(parts, "trend ↓")
	}
	if rule.TrendDown && rule.ThresholdMin != 0 {
		s := fmt.Sprintf("(trend ↓ && value <= %s)", util.ToFixed(rule.ThresholdMin, 3))
		parts = append(parts, s)
	}
	if !rule.TrendDown && rule.ThresholdMin != 0 {
		s := fmt.Sprintf("value <= %s", util.ToFixed(rule.ThresholdMin, 3))
		parts = append(parts, s)
	}
	rule.Repr = strings.Join(parts, " || ")
}

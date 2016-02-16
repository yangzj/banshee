// Copyright 2016 Eleme Inc. All rights reserved.

package models

import (
	"github.com/eleme/banshee/config"
	"github.com/eleme/banshee/util/assert"
	"testing"
)

func TestRuleTest(t *testing.T) {
	var rule *Rule
	// TrendUp
	rule = &Rule{TrendUp: true}
	assert.Ok(t, rule.Test(&Metric{}, &Index{Score: 1.2}, nil))
	assert.Ok(t, !rule.Test(&Metric{}, &Index{Score: 0.8}, nil))
	// TrendDown
	rule = &Rule{TrendDown: true}
	assert.Ok(t, rule.Test(&Metric{}, &Index{Score: -1.2}, nil))
	assert.Ok(t, !rule.Test(&Metric{}, &Index{Score: 1.2}, nil))
	// TrendUp And Value >= X
	rule = &Rule{TrendUp: true, ThresholdMax: 39}
	assert.Ok(t, rule.Test(&Metric{Value: 50}, &Index{Score: 1.3}, nil))
	assert.Ok(t, !rule.Test(&Metric{Value: 38}, &Index{Score: 1.5}, nil))
	assert.Ok(t, !rule.Test(&Metric{Value: 60}, &Index{Score: 0.9}, nil))
	// TrendDown And Value <= X
	rule = &Rule{TrendDown: true, ThresholdMin: 40}
	assert.Ok(t, rule.Test(&Metric{Value: 10}, &Index{Score: -1.2}, nil))
	assert.Ok(t, !rule.Test(&Metric{Value: 41}, &Index{Score: -1.2}, nil))
	assert.Ok(t, !rule.Test(&Metric{Value: 12}, &Index{Score: -0.2}, nil))
	// (TrendUp And Value >= X) Or TrendDown
	rule = &Rule{TrendUp: true, TrendDown: true, ThresholdMax: 90}
	assert.Ok(t, rule.Test(&Metric{Value: 100}, &Index{Score: 1.1}, nil))
	assert.Ok(t, rule.Test(&Metric{}, &Index{Score: -1.1}, nil))
	assert.Ok(t, !rule.Test(&Metric{}, &Index{Score: -0.1}, nil))
	assert.Ok(t, !rule.Test(&Metric{Value: 89}, &Index{Score: 1.3}, nil))
	assert.Ok(t, !rule.Test(&Metric{Value: 189}, &Index{Score: 0.3}, nil))
	// (TrendUp And Value >= X) Or (TrendDown And Value <= X)
	rule = &Rule{TrendUp: true, TrendDown: true, ThresholdMax: 90, ThresholdMin: 10}
	assert.Ok(t, rule.Test(&Metric{Value: 100}, &Index{Score: 1.2}, nil))
	assert.Ok(t, rule.Test(&Metric{Value: 9}, &Index{Score: -1.2}, nil))
	assert.Ok(t, !rule.Test(&Metric{Value: 12}, &Index{Score: 1.2}, nil))
	assert.Ok(t, !rule.Test(&Metric{Value: 102}, &Index{Score: 0.2}, nil))
	assert.Ok(t, !rule.Test(&Metric{Value: 2}, &Index{Score: 0.9}, nil))
	// Default thresholdMaxs
	cfg := config.New()
	cfg.Detector.DefaultThresholdMaxs["fo*"] = 300
	rule = &Rule{TrendUp: true}
	assert.Ok(t, rule.Test(&Metric{Value: 310, Name: "foo"}, &Index{Score: 1.3}, cfg))
	assert.Ok(t, !rule.Test(&Metric{Value: 120, Name: "foo"}, &Index{Score: 1.3}, cfg))
}

package types

import (
	"encoding/json"
	"testing"
)

func TestAppMetricsRoundTrip(t *testing.T) {
	cpuLimit := 0.5
	memLimit := float64(536870912) // 512 MiB

	roundTrip(t, "AppMetricsResponse/populated", AppMetricsResponse{
		AppID:             42,
		Range:             "1h",
		Start:             1780000000,
		End:               1780003600,
		StepSeconds:       30,
		RateWindowSeconds: 240,
		Metrics: AppMetrics{
			CPU: MetricSeries{
				Unit:   UnitCores,
				Points: []MetricPoint{{T: 1780000000, V: 0.12}, {T: 1780000030, V: 0.18}},
				Limit:  &cpuLimit,
			},
			Memory: MetricSeries{
				Unit:   UnitBytes,
				Points: []MetricPoint{{T: 1780000000, V: 1048576}},
				Limit:  &memLimit,
			},
		},
	})

	// Empty series (never deployed / scaled to zero) + no limit set.
	roundTrip(t, "AppMetricsResponse/empty", AppMetricsResponse{
		AppID:             7,
		Range:             "15m",
		Start:             1780000000,
		End:               1780000900,
		StepSeconds:       15,
		RateWindowSeconds: 240,
		Metrics: AppMetrics{
			CPU:    MetricSeries{Unit: UnitCores, Points: []MetricPoint{}},
			Memory: MetricSeries{Unit: UnitBytes, Points: []MetricPoint{}},
		},
	})
}

// Limit must be omitted from the wire when nil, present (even if 0) when set.
func TestMetricSeriesLimitOmitempty(t *testing.T) {
	noLimit, _ := json.Marshal(MetricSeries{Unit: UnitCores, Points: []MetricPoint{}})
	if got := string(noLimit); got != `{"unit":"cores","points":[]}` {
		t.Errorf("nil limit should be omitted, got %s", got)
	}
	zero := 0.0
	withLimit, _ := json.Marshal(MetricSeries{Unit: UnitCores, Points: []MetricPoint{}, Limit: &zero})
	if got := string(withLimit); got != `{"unit":"cores","points":[],"limit":0}` {
		t.Errorf("set limit (0) should be present, got %s", got)
	}
}

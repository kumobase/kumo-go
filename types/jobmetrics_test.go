package types

import "testing"

func TestJobExecutionMetricsRoundTrip(t *testing.T) {
	cpuLimit := 0.25
	memLimit := float64(134217728) // 128 MiB

	roundTrip(t, "JobExecutionMetricsResponse/populated", JobExecutionMetricsResponse{
		JobID:             42,
		ExecutionID:       7,
		Start:             1780000000,
		End:               1780000090,
		StepSeconds:       15,
		RateWindowSeconds: 240,
		Metrics: JobExecutionMetrics{
			CPU: MetricSeries{
				Unit:   UnitCores,
				Points: []MetricPoint{{T: 1780000000, V: 0.10}, {T: 1780000015, V: 0.22}},
				Limit:  &cpuLimit,
			},
			Memory: MetricSeries{
				Unit:   UnitBytes,
				Points: []MetricPoint{{T: 1780000000, V: 1048576}},
				Limit:  &memLimit,
			},
		},
	})

	// Empty series (pending / sub-scrape run / aged out) + no limit.
	roundTrip(t, "JobExecutionMetricsResponse/empty", JobExecutionMetricsResponse{
		JobID:       9,
		ExecutionID: 1,
		Metrics: JobExecutionMetrics{
			CPU:    MetricSeries{Unit: UnitCores, Points: []MetricPoint{}},
			Memory: MetricSeries{Unit: UnitBytes, Points: []MetricPoint{}},
		},
	})
}

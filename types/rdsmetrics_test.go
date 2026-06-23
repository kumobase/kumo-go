package types

import "testing"

func TestRDSMetricsRoundTrip(t *testing.T) {
	cpuLimit := 2.0
	memLimit := float64(2147483648)  // 2 GiB
	diskLimit := float64(10737418240) // 10 GiB volume capacity
	connLimit := 100.0

	roundTrip(t, "RDSMetricsResponse/populated", RDSMetricsResponse{
		InstanceID:        42,
		Range:             "1h",
		Start:             1780000000,
		End:               1780003600,
		StepSeconds:       30,
		RateWindowSeconds: 240,
		Metrics: RDSMetrics{
			CPU: MetricSeries{
				Unit:   UnitCores,
				Points: []MetricPoint{{T: 1780000000, V: 0.42}, {T: 1780000030, V: 0.51}},
				Limit:  &cpuLimit,
			},
			Memory: MetricSeries{
				Unit:   UnitBytes,
				Points: []MetricPoint{{T: 1780000000, V: 536870912}},
				Limit:  &memLimit,
			},
			Disk: MetricSeries{
				Unit:   UnitBytes,
				Points: []MetricPoint{{T: 1780000000, V: 3221225472}},
				Limit:  &diskLimit,
			},
			Database: RDSDatabaseMetrics{
				Up:             MetricSeries{Unit: UnitCount, Points: []MetricPoint{{T: 1780000000, V: 1}}},
				Connections:    MetricSeries{Unit: UnitCount, Points: []MetricPoint{{T: 1780000000, V: 12}}, Limit: &connLimit},
				ReplicationLag: MetricSeries{Unit: UnitSeconds, Points: []MetricPoint{{T: 1780000000, V: 0.5}}},
				Commits:        MetricSeries{Unit: UnitPerSecond, Points: []MetricPoint{{T: 1780000000, V: 42.5}}},
				Rollbacks:      MetricSeries{Unit: UnitPerSecond, Points: []MetricPoint{{T: 1780000000, V: 0.1}}},
				CacheHitRatio:  MetricSeries{Unit: UnitPercent, Points: []MetricPoint{{T: 1780000000, V: 99.3}}},
			},
		},
	})

	// Empty series (just provisioned / suspended) — disk limit still present
	// (authoritative from storage_gb) even with no usage samples yet.
	roundTrip(t, "RDSMetricsResponse/empty", RDSMetricsResponse{
		InstanceID:        7,
		Range:             "15m",
		Start:             1780000000,
		End:               1780000900,
		StepSeconds:       15,
		RateWindowSeconds: 240,
		Metrics: RDSMetrics{
			CPU:    MetricSeries{Unit: UnitCores, Points: []MetricPoint{}},
			Memory: MetricSeries{Unit: UnitBytes, Points: []MetricPoint{}},
			Disk:   MetricSeries{Unit: UnitBytes, Points: []MetricPoint{}, Limit: &diskLimit},
			Database: RDSDatabaseMetrics{
				Up:             MetricSeries{Unit: UnitCount, Points: []MetricPoint{}},
				Connections:    MetricSeries{Unit: UnitCount, Points: []MetricPoint{}},
				ReplicationLag: MetricSeries{Unit: UnitSeconds, Points: []MetricPoint{}},
				Commits:        MetricSeries{Unit: UnitPerSecond, Points: []MetricPoint{}},
				Rollbacks:      MetricSeries{Unit: UnitPerSecond, Points: []MetricPoint{}},
				CacheHitRatio:  MetricSeries{Unit: UnitPercent, Points: []MetricPoint{}},
			},
		},
	})
}

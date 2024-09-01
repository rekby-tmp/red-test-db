package utils

import (
	"slices"
	"sync"
	"time"
)

type LatencyMetric struct {
	m         sync.Mutex
	durations []time.Duration
}

func (m *LatencyMetric) Add(duration time.Duration) {
	m.m.Lock()
	m.durations = append(m.durations, duration)
	m.m.Unlock()
}

func (m *LatencyMetric) AddSince(t time.Time) {
	m.Add(time.Since(t))
}

func (m *LatencyMetric) Stat(p ...float64) LatencyStats {
	m.m.Lock()
	durations := m.durations
	m.durations = make([]time.Duration, 0, len(durations)*2)
	m.m.Unlock()

	slices.Sort(durations)

	res := make([]time.Duration, len(p))
	if len(durations) == 0 {
		return LatencyStats{
			Durations: res,
		}
	}

	for resIndex, percentile := range p {
		index := calcIndex(len(durations), percentile)
		d := durations[index]
		res[resIndex] = d
	}

	return LatencyStats{
		TotalCount:  len(durations),
		Percentiles: p,
		Durations:   res,
	}
}

type LatencyStats struct {
	TotalCount  int
	Percentiles []float64
	Durations   []time.Duration
}

func calcIndex(l int, p float64) int {
	if p >= 1.0 {
		return l - 1
	}

	return int(float64(l) * p)
}

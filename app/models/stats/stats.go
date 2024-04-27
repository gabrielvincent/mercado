package stats

import "time"

type Number interface {
	~int | ~int32 | ~int64 | ~float32 | ~float64 | ~uint | ~uint32 | ~uint64
}

type MetricType string

type Metric interface{}

const (
	Percent MetricType = "percent"
	Amount  MetricType = "amount"
)

type ComparisonMetric[T Number] struct {
	Name          string
	Type          MetricType
	TargetValue   T
	CompareValue  T
	IncreaseValue T
}

type PeriodComparison struct {
	TargetDateStart  time.Time
	TargetDateEnd    time.Time
	CompareDateStart time.Time
	CompareDateEnd   time.Time
	Metrics          map[string]Metric
}

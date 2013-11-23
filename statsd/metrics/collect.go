package metrics

import (
    "sync"
)

type StatsCollection struct {
    Timers map[string][]float64
    TimerCounters map[string]float64

    Counters map[string]float64

    Gauges map[string]float64

    Sets map[string]map[string]interface{}
}

func NewStatsCollection (past *StatsCollection) *StatsCollection {

    s := &StatsCollection{
        Timers: make(map[string][]float64),
        TimerCounters: make(map[string]float64),
        Counters: make(map[string]float64),
        Sets: make(map[string]map[string]interface{}),
    }

    if past != nil {
        s.Gauges = past.Gauges
    } else {
        s.Gauges = make(map[string]float64)
    }

    return s
}

type MetricCollector interface {
    Flush()
}

func NewMetricCollector(inbound <-chan *Metric, workers int) (MetricCollector, <-chan *StatsCollection) {

    outbound := make(chan *StatsCollection)

    p := metricCollectorStruct{
        inbound_metrics: inbound,
        outbound_stats: outbound,
        stats: NewStatsCollection(nil),
    }

    for i := 0; i < workers; i++ {
        go p.collect()
    }

    return p, outbound
}

type metricCollectorStruct struct {
    sync.Mutex

    inbound_metrics <-chan *Metric
    outbound_stats chan<- *StatsCollection

    stats *StatsCollection
}

func (p metricCollectorStruct) Flush() {
        p.Lock()

        stats := p.stats
        p.stats = NewStatsCollection(stats)

        defer p.Unlock()

        p.outbound_stats <- stats
}

func (p metricCollectorStruct) collect () {

    for m := range p.inbound_metrics {
        p.Lock()
        switch m.Type {
        case "ms":
            if _, ok := p.stats.Timers[m.Key]; !ok {
                p.stats.Timers[m.Key] = make([]float64, 0)
            }
            p.stats.Timers[m.Key] = append(p.stats.Timers[m.Key], m.FloatValue)
            p.stats.TimerCounters[m.Key] += (1 / m.SampleRate)
        case "g":
            if m.Value[0] == '+' || m.Value[0] == '-' {
                p.stats.Gauges[m.Key] += m.FloatValue
            } else {
                p.stats.Gauges[m.Key] = m.FloatValue
            }
        case "s":
            if _, ok := p.stats.Sets[m.Key]; !ok {
                p.stats.Sets[m.Key] = make(map[string]interface{})
            }
            p.stats.Sets[m.Key][m.Value] = nil
        default:
            p.stats.Counters[m.Key] += m.FloatValue * (1 / m.SampleRate)
        }

        p.Unlock()
    }
}


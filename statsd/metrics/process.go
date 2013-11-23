package metrics

import (
    "time"
)

type ProcessedStats struct {
    StatsCollection

    CounterRates map[string]float64
}

func NewProcessedStats() *ProcessedStats {
    return &ProcessedStats{
        CounterRates: make(map[string]float64),
    }
}

func ProcessFromChannel(inbound <-chan *StatsCollection, flushInterval time.Duration, workers int) <-chan *ProcessedStats {
    out := make(chan *ProcessedStats)

    for i := 0; i < workers; i++ {
        go func (){
            for rawStats := range inbound {
                out <- ProcessStats(rawStats, flushInterval)
            }
        }()
    }
    return out
}

func ProcessStats(rawStats *StatsCollection, flushInterval time.Duration) *ProcessedStats {
    stats := NewProcessedStats()

    flushSeconds := flushInterval.Seconds()

    for key, value := range rawStats.Counters {
        stats.CounterRates[key] = value / flushSeconds
    }

    return stats
}

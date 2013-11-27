package metrics

import (
    "time"
    "sort"
    "math"
)

type ProcessedStats struct {
    StatsCollection

    CounterRates map[string]float64

    ExtraTimerData map[string]TimerData
}

type TimerData struct {
    Std float64
    Upper float64
    Lower float64
    Count float64
    Count_ps float64
    Sum float64
    Mean float64
    Median float64

    Percents []TimerPercentile
}

type TimerPercentile struct {
    Percent float64

    Mean float64
    Sum float64
    Boundry float64
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

    for key, values := range rawStats.Timers {

        data := TimerData{
            Percents: make([]TimerPercentile)
        }

        sort.Float64s(values)

        data.Sum = cumulativeValues[-1]
        data.Min

        var count = len(values)
        var min = values[0]
        var max = values[-1]

        var cumulativeValues = make([]float64, 0, count)

        cumulativeValues[0] = min
        for i := 1; i < count; i++ {
            cumulativeValues[i] = values[i] + cumulativeValues[i-1]
        }

        for _, percent := range percentThreshold {
            pData := {
                Percent:  percent,

                Sum: min,
                Mean: min,
                Boundry: max,
            }

            if count > 1 {
                var numInThreshold = math.Floor(math.Abs(percent) / 100 * count)
                if numInThreshold == 0 {
                    continue
                }

                if (percent > 0) {
                    pData.Boundry = values[numInThreshold - 1]
                    pData.Sum = cumulativeValues[numInThreshold - 1]
                } else {
                    pData.Boundry = values[count - numInThreshold]
                    pData.Sum = cumulativeValues[count - 1] - cumulativeValues[count - numInThreshold - 1]
                }
                pData.Mean = pData.Sum / numInThreshold
            }
            data.Percents = append(data.Percents, pData)
        }

        // process_metrics.js +69

    }

    return stats
}


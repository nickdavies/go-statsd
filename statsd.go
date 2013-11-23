package main

import (
    "time"
    "fmt"
)

import (
    "./statsd/servers"
    "./statsd/metrics"
)

func main() {

    server, err := servers.GetServer("udp")
    if err != nil {
        panic(err)
    }

    inbound, _, err := server.Run(servers.Config{":1234"})
    if err != nil {
        panic(err)
    }

    rawMetrics := make(chan *metrics.Metric)

    collector, rawStats := metrics.NewMetricCollector(rawMetrics, 1)

    flushInterval := 1 * time.Second

    go func () {
        for _ = range time.Tick(flushInterval) {
            collector.Flush()
        }
    }()

    stats := metrics.ProcessFromChannel(rawStats, flushInterval, 2)

    go func () {
        for _ = range stats {
            // write to backend??
        }
    }()

    for p := range inbound {
        fmt.Println("Recv: ", p.Message)

        ms, err := metrics.ParseMessage(p.Message)
        if err != nil {
            //TODO: Deal with error
            continue
        }
        for _, m := range ms {
            rawMetrics <-m
        }

        server.DestroyPacket(p)
    }

}

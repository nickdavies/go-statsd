package metrics

import (
    "fmt"
    "regexp"
    "strconv"
    "strings"
)

// Cleaning Regex
var whitespace = regexp.MustCompile(`\s+`)
var slash = regexp.MustCompile(`/`)
var notallowed = regexp.MustCompile(`[^a-zA-Z_\-0-9\.]`)

// Validation Regex
var number = regexp.MustCompile(`^([\d\.]+$)`)
var signedNumber = regexp.MustCompile(`^([\-\+\d\.]+$)`)

type MetricType string

type Metric struct {
    Key string
    Value string
    Type MetricType
    SampleRate float64
}

type ValidationError struct {
    message string
}

func (e ValidationError) Error() string {
    return e.message
}

func ParseMessage(message string) ([]*Metric, error) {
    var metrics = make([]*Metric, 0, 5)

    commands := strings.Split(message, "\n")
    for _, command := range commands {
        if len(command) == 0 {
            continue
        }
        err := parseCommand(command, &metrics)
        if err != nil {
            //TODO: handle error
            fmt.Println("parse command ERROR!", err)
            continue
        }
    }

    return metrics, nil
}

func parseCommand(command string, metrics *[]*Metric) error {
    rawMetrics := strings.Split(command, ":")

    key := rawMetrics[0]
    rawMetrics = rawMetrics[1:]

    key = whitespace.ReplaceAllLiteralString(key, "_")
    key = slash.ReplaceAllLiteralString(key, "-")
    key = notallowed.ReplaceAllLiteralString(key, "")

    if len(rawMetrics) == 0 {
        rawMetrics = []string{"1"}
    }

    for _, rawMetric := range rawMetrics {
        metric, err := parseMetric(key, rawMetric)
        if err != nil {
            //TODO: handle error
            fmt.Println("parse metric ERROR!", err)
            continue
        }

        *metrics = append(*metrics, metric)
    }

    return nil
}

func parseMetric(key, rawMetric string) (*Metric, error) {
    var err error

    fields := strings.Split(rawMetric, "|")

    if len(fields) < 2 {
        return nil, ValidationError{"Validation failed: no type field"}
    }

    var metric = &Metric{
        Key: key,
        Value: fields[0],
        Type: MetricType(strings.TrimSpace(fields[1])),
    }

    if len(fields) >= 3 && len(fields[2]) != 0 {
        if fields[2][0] != '@' {
            return nil, ValidationError{"Validation failed: sample rate has no @ at the start"}
        }

        metric.SampleRate, err = strconv.ParseFloat(fields[2][1:], 64)
        if err != nil {
            return nil, ValidationError{"Validation failed: sample rate is not a valid float"}
        }

        if metric.SampleRate < 0 {
            return nil, ValidationError{"Validation failed: sample rate is < 0"}
        }
    } else {
        metric.SampleRate = 1
    }

    switch metric.Type {
    case "s":
    case "g":
        if !signedNumber.MatchString(fields[0]) {
            return nil, ValidationError{"Validation failed: guage value is not a number"}
        }
    default:
        if !number.MatchString(fields[0]) {
            return nil, ValidationError{"Validation failed: value is not a number"}
        }
    }

    return metric, nil
}


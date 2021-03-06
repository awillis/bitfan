// Code generated by "bitfanDoc "; DO NOT EDIT
package statsd

import "github.com/awillis/bitfan/processors/doc"

func (p *processor) Doc() *doc.Processor {
	return &doc.Processor{
		Name:       "statsd",
		ImportPath: "github.com/awillis/bitfan/processors/output-statsd",
		Doc:        "",
		DocShort:   "",
		Options: &doc.ProcessorOptions{
			Doc: "",
			Options: []*doc.ProcessorOption{
				&doc.ProcessorOption{
					Name:           "Host",
					Alias:          "host",
					Doc:            "The address of the statsd server.",
					Required:       false,
					Type:           "string",
					DefaultValue:   nil,
					PossibleValues: []string{},
					ExampleLS:      "",
				},
				&doc.ProcessorOption{
					Name:           "Port",
					Alias:          "port",
					Doc:            "The port to connect to on your statsd server.",
					Required:       false,
					Type:           "int",
					DefaultValue:   nil,
					PossibleValues: []string{},
					ExampleLS:      "",
				},
				&doc.ProcessorOption{
					Name:           "Protocol",
					Alias:          "protocol",
					Doc:            "",
					Required:       false,
					Type:           "string",
					DefaultValue:   nil,
					PossibleValues: []string{},
					ExampleLS:      "",
				},
				&doc.ProcessorOption{
					Name:           "Sender",
					Alias:          "sender",
					Doc:            "The name of the sender. Dots will be replaced with underscores.",
					Required:       false,
					Type:           "string",
					DefaultValue:   nil,
					PossibleValues: []string{},
					ExampleLS:      "",
				},
				&doc.ProcessorOption{
					Name:           "Count",
					Alias:          "count",
					Doc:            "A count metric. metric_name => count as hash",
					Required:       false,
					Type:           "hash",
					DefaultValue:   nil,
					PossibleValues: []string{},
					ExampleLS:      "",
				},
				&doc.ProcessorOption{
					Name:           "Decrement",
					Alias:          "decrement",
					Doc:            "A decrement metric. Metric names as array.",
					Required:       false,
					Type:           "array",
					DefaultValue:   nil,
					PossibleValues: []string{},
					ExampleLS:      "",
				},
				&doc.ProcessorOption{
					Name:           "Gauge",
					Alias:          "gauge",
					Doc:            "A gauge metric. metric_name => gauge as hash.",
					Required:       false,
					Type:           "hash",
					DefaultValue:   nil,
					PossibleValues: []string{},
					ExampleLS:      "",
				},
				&doc.ProcessorOption{
					Name:           "Increment",
					Alias:          "increment",
					Doc:            "An increment metric. Metric names as array.",
					Required:       false,
					Type:           "array",
					DefaultValue:   nil,
					PossibleValues: []string{},
					ExampleLS:      "",
				},
				&doc.ProcessorOption{
					Name:           "NameSpace",
					Alias:          "namespace",
					Doc:            "The statsd namespace to use for this metric.",
					Required:       false,
					Type:           "string",
					DefaultValue:   nil,
					PossibleValues: []string{},
					ExampleLS:      "",
				},
				&doc.ProcessorOption{
					Name:           "SampleRate",
					Alias:          "sample_rate",
					Doc:            "The sample rate for the metric.",
					Required:       false,
					Type:           "float32",
					DefaultValue:   nil,
					PossibleValues: []string{},
					ExampleLS:      "",
				},
				&doc.ProcessorOption{
					Name:           "Set",
					Alias:          "set",
					Doc:            "A set metric. metric_name => \"string\" to append as hash",
					Required:       false,
					Type:           "hash",
					DefaultValue:   nil,
					PossibleValues: []string{},
					ExampleLS:      "",
				},
				&doc.ProcessorOption{
					Name:           "Timing",
					Alias:          "timing",
					Doc:            "A timing metric. metric_name => duration as hash",
					Required:       false,
					Type:           "hash",
					DefaultValue:   nil,
					PossibleValues: []string{},
					ExampleLS:      "",
				},
			},
		},
		Ports: []*doc.ProcessorPort{},
	}
}

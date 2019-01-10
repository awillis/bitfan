// Code generated by "bitfanDoc "; DO NOT EDIT
package uuid

import "github.com/awillis/bitfan/processors/doc"

func (p *processor) Doc() *doc.Processor {
	return &doc.Processor{
		Name:       "uuid",
		ImportPath: "github.com/awillis/bitfan/processors/filter-uuid",
		Doc:        "The uuid filter allows you to generate a UUID and add it as a field to each processed event.\n\nThis is useful if you need to generate a string that’s unique for every event, even if the same input is processed multiple times. If you want to generate strings that are identical each time a event with a given content is processed (i.e. a hash) you should use the fingerprint filter instead.\n\nThe generated UUIDs follow the version 4 definition in RFC 4122).",
		DocShort:   "Adds a UUID to events",
		Options: &doc.ProcessorOptions{
			Doc: "",
			Options: []*doc.ProcessorOption{
				&doc.ProcessorOption{
					Name:           "processors.CommonOptions",
					Alias:          ",squash",
					Doc:            "",
					Required:       false,
					Type:           "processors.CommonOptions",
					DefaultValue:   nil,
					PossibleValues: []string{},
					ExampleLS:      "",
				},
				&doc.ProcessorOption{
					Name:           "Overwrite",
					Alias:          "",
					Doc:            "If the value in the field currently (if any) should be overridden by the generated UUID.\nDefaults to false (i.e. if the field is present, with ANY value, it won’t be overridden)",
					Required:       false,
					Type:           "bool",
					DefaultValue:   nil,
					PossibleValues: []string{},
					ExampleLS:      "",
				},
				&doc.ProcessorOption{
					Name:           "Target",
					Alias:          "target",
					Doc:            "Add a UUID to a field",
					Required:       true,
					Type:           "string",
					DefaultValue:   nil,
					PossibleValues: []string{},
					ExampleLS:      "",
				},
			},
		},
		Ports: []*doc.ProcessorPort{},
	}
}
